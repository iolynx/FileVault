package folders

import (
	"context"
	"log"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service handles folder-related business logic, including creation, updating, deletion,
// moving folders, and listing selectable folders.
type Service struct {
	repo    *Repository
	storage storage.Storage
}

// NewService creates a new instance of the folder Service.
// - repo: repository providing database operations for folders and files.
// - storage: storage interface used for managing file blobs associated with folders.
func NewService(repo *Repository, storage storage.Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

// CreateFolder creates a new folder for the authenticated user.
// - Validates folder name is not empty.
// - Validates the parent folder belongs to the user (if provided).
// Returns the created folder or an error.
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

// UpdateFolder renames a folder for the authenticated user.
// - Validates user ownership of the folder.
// - Returns an error if the folder is not found, the user does not own it, or the name is empty.
// - Note: This only handles renaming; moving folders is not handled here.
// Returns the updated folder as FolderResponse.
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

// GetSelectableFolders returns a list of folders that the authenticated user
// can select (for moving files and folders). If folderID is provided, it
// validates ownership of that folder, and uses the folderID to determine what folders are selectable.
// Returns an error if the user is unauthorized or if the folder does not exist or is forbidden.
func (s *Service) GetSelectableFolders(ctx context.Context, folderID *uuid.UUID) ([]Folder, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return nil, apierror.NewUnauthorizedError()
	}

	if folderID != nil {
		folder, err := s.repo.GetFolderByID(ctx, *folderID)
		if err != nil {
			return nil, apierror.NewNotFoundError("Folder")
		}
		if folder.OwnerID != userID {
			return nil, apierror.NewForbiddenError()
		}
	}

	params := sqlc.ListSelectableFoldersParams{
		OwnerID: userID,
	}
	if folderID != nil {
		params.CurrentFolderID = pgtype.UUID{Bytes: *folderID, Valid: true}
	}

	rows, err := s.repo.ListSelectableFolders(ctx, params)
	if err != nil {
		return nil, apierror.NewInternalServerError()
	}

	folders := make([]Folder, len(rows))
	for i, row := range rows {
		folders[i] = Folder{
			ID:             row.ID,
			Name:           row.Name,
			CreatedAt:      row.CreatedAt.Time,
			ParentFolderID: util.ToUUIDPtr(row.ParentFolderID),
		}
	}
	return folders, nil
}

// DeleteFolder deletes a folder and all its contents from the database and storage.
// - Validates user ownership of the folder.
// - Deletes subfolders and file records automatically via ON DELETE CASCADE.
// - Checks all blobs in the folder hierarchy for cleanup and deletes unreferenced blobs from storage.
// Returns an error if the user is unauthorized, the folder does not exist, or deletion fails.
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
	blobIDs, err := s.repo.GetBlobIDsInFolderHierarchy(ctx, folderID)
	if err != nil && err != pgx.ErrNoRows {
		return apierror.NewInternalServerError("Could not retrieve files for deletion")
	}

	// Delete the folder record from the database
	// ON DELETE CASCADE deletes all subfolders and file records automatically
	if err = s.repo.DeleteFolder(ctx, folderID); err != nil {
		return apierror.NewInternalServerError("Failed to delete folder from database")
	}
	log.Printf("Deleted folder %s and all its contents from database records.", folderID)

	// Checking blobs for cleanup
	log.Printf("Checking %d blobs for cleanup...", len(blobIDs))
	for _, blobID := range blobIDs {
		storagePath, err := s.repo.queries.DeleteBlobIfUnused(ctx, blobID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// the blob is still referenced by another file, we do nothing.
				continue
			}
			log.Printf("Error during blob cleanup for %s: %v", blobID, err)
			continue
		}

		// The blob record has been deleted, and we can safely delete from storage.
		log.Printf("Blob %s is now unreferenced, deleting object %s from storage.", blobID, storagePath)
		if err := s.storage.DeleteBlob(ctx, storagePath); err != nil {
			log.Printf("CRITICAL: Failed to delete object %s from storage: %v", storagePath, err)
		}
	}
	return nil
}

// UpdateFolderParent updates the parent folder of the specified folder.
// - Validates that the authenticated user owns both the folder and the target parent (if provided).
// - Moves the folder under the new parent or to root if TargetFolderID is nil.
// Returns an error if the user is unauthorized, the folder or target parent is forbidden, or if the operation fails.
func (s *Service) UpdateFolderParent(ctx context.Context, folderID uuid.UUID, req UpdateFolderParentRequest) error {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	// ownership checks
	folder, err := s.repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return apierror.NewInternalServerError()
	}
	if folder.OwnerID != userID {
		return apierror.NewForbiddenError()
	}

	// if the destination (parent) folder is not null, check its ownership
	if req.TargetFolderID != nil {
		newParentFolder, err := s.repo.GetFolderByID(ctx, *req.TargetFolderID)
		if err != nil {
			return apierror.NewInternalServerError("Could not find parent folder")
		}
		if newParentFolder.OwnerID != userID {
			return apierror.NewForbiddenError()
		}
		if newParentFolder.ID == folderID {
			return apierror.NewBadRequestError("Source and Destination cannot be the same")
		}
	}

	params := sqlc.UpdateFolderParentFolderParams{
		ID: folderID,
	}

	if req.TargetFolderID != nil {
		params.ParentFolderID = pgtype.UUID{Bytes: *req.TargetFolderID, Valid: true}
	}

	return s.repo.UpdateFolderParentFolder(ctx, params)
}
