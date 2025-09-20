package folders

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo    *Repository
	storage storage.Storage
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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

func (s *Service) UpdateFolder(ctx context.Context, folderID uuid.UUID, req UpdateFolderRequest) (sqlc.Folder, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return sqlc.Folder{}, apierror.NewUnauthorizedError()
	}

	folderToUpdate, err := s.repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return sqlc.Folder{}, apierror.NewInternalServerError("Folder not found")
	}
	if folderToUpdate.OwnerID != userID {
		return sqlc.Folder{}, apierror.NewForbiddenError()
	}

	if req.Name == "" {
		return sqlc.Folder{}, apierror.NewBadRequestError("Folder name cannot be empty!")
	}

	// Note: This implementation only handles renaming, not moving.
	params := sqlc.UpdateFolderParams{
		ID:             folderID,
		Name:           req.Name,
		ParentFolderID: folderToUpdate.ParentFolderID,
	}
	return s.repo.UpdateFolder(ctx, params)
}

func (s *Service) DeleteFolder(ctx context.Context, folderID uuid.UUID) error {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	folderToDelete, err := s.repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return apierror.NewInternalServerError("Folder not found")
	}
	if folderToDelete.OwnerID != userID {
		return apierror.NewForbiddenError()
	}

	return s.repo.DeleteFolder(ctx, folderID)
}
