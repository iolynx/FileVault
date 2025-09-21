package folders

import (
	"context"
	"log"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo    *Repository
	storage storage.Storage
}

func NewService(repo *Repository, storage storage.Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

type CreateFolderRequest struct {
	Name           string     `json:"name"`
	ParentFolderID *uuid.UUID `json:"parent_folder_id"`
}

func (s *Service) CreateFolder(ctx context.Context, req CreateFolderRequest) (sqlc.Folder, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return sqlc.Folder{}, apierror.NewUnauthorizedError()
	}

	if req.Name == "" {
		return sqlc.Folder{}, apierror.NewBadRequestError("Folder name cannot be empty")
	}

	if req.ParentFolderID != nil {
		parentFolder, err := s.repo.GetFolderByID(ctx, *req.ParentFolderID)
		if err != nil {
			return sqlc.Folder{}, apierror.NewInternalServerError("Could not find parent folder")
		}
		if parentFolder.OwnerID != userID {
			return sqlc.Folder{}, apierror.NewForbiddenError()
		}
	}

	params := sqlc.CreateFolderParams{
		Name:    req.Name,
		OwnerID: userID,
	}

	if req.ParentFolderID != nil {
		params.ParentFolderID = pgtype.UUID{Bytes: *req.ParentFolderID, Valid: true}
	}

	return s.repo.CreateFolder(ctx, params)
}

type UpdateFolderRequest struct {
	Name string `json:"name"`
}

type FolderResponse struct {
	ID           uuid.UUID `json:"id"`
	Filename     string    `json:"filename"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UserOwnsFile bool      `json:"user_owns_file"`
	ItemType     string    `json:"item_type"`
}

func (s *Service) UpdateFolder(ctx context.Context, folderID uuid.UUID, req UpdateFolderRequest) (FolderResponse, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return FolderResponse{}, apierror.NewUnauthorizedError()
	}

	folderToUpdate, err := s.repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return FolderResponse{}, apierror.NewInternalServerError("Folder not found")
	}
	if folderToUpdate.OwnerID != userID {
		return FolderResponse{}, apierror.NewForbiddenError()
	}

	if req.Name == "" {
		return FolderResponse{}, apierror.NewBadRequestError("Folder name cannot be empty!")
	}

	// Note: This implementation only handles renaming, not moving.
	params := sqlc.UpdateFolderParams{
		ID:             folderID,
		Name:           req.Name,
		ParentFolderID: folderToUpdate.ParentFolderID,
	}

	res, err := s.repo.UpdateFolder(ctx, params)
	if err != nil {
		return FolderResponse{}, err
	}

	return FolderResponse{
		ID:           res.ID,
		Filename:     res.Filename,
		UploadedAt:   res.CreatedAt.Time,
		UserOwnsFile: true,
		ItemType:     "Folder",
	}, nil
}

func (s *Service) DeleteFolder(ctx context.Context, folderID uuid.UUID) error {
	ownerID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	folder, err := s.repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return apierror.NewNotFoundError("Folder")
	}

	// Ownership check
	if folder.OwnerID != ownerID {
		return apierror.NewForbiddenError()
	}

	// get all object keys for files within this folder and its subfolders
	storagePaths, err := s.repo.GetObjectKeysInFolderHierarchy(ctx, folderID)
	if err != nil && err != pgx.ErrNoRows {
		return apierror.NewInternalServerError("Could not retrieve files for deletion")
	}

	// Delete the actual objects from MinIO storage
	if len(storagePaths) > 0 {
		err = s.storage.DeleteBlobs(ctx, storagePaths)
		if err != nil {
			return apierror.NewInternalServerError("Failed to delete files from storage")
		}
	}

	// Delete the folder record from the database
	// ON DELETE CASCADE deletes all subfolders and file records automatically
	err = s.repo.DeleteFolder(ctx, folderID)
	if err != nil {
		return apierror.NewInternalServerError("Failed to delete folder from database")
	}

	if len(storagePaths) > 0 {
		err = s.repo.queries.DeleteBlobsByStoragePaths(ctx, storagePaths)
		if err != nil {
			log.Printf("Could not delete some blobs records: %v", err)
		}
	}

	return nil
}
