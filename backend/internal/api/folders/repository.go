package folders

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/google/uuid"
)

// Repository handles database operations related to folders
type Repository struct {
	queries *sqlc.Queries
}

// NewRepository creates a new Repository instance with the provided database queries.
// This Repository can be used to perform folder related database operations.
func NewRepository(db *sqlc.Queries) *Repository {
	return &Repository{
		queries: db,
	}
}

// CreateFolder creates a new folder in the db with the given parameters.
// Returns the created Folder and an error if the creation fails.
func (r *Repository) CreateFolder(ctx context.Context, arg sqlc.CreateFolderParams) (sqlc.Folder, error) {
	return r.queries.CreateFolder(ctx, arg)
}

// GetFolderByID fetches a folder by its UUID.
// Returns an error if the folder does not exist.
func (r *Repository) GetFolderByID(ctx context.Context, folderID uuid.UUID) (sqlc.Folder, error) {
	return r.queries.GetFolderByID(ctx, folderID)
}

// UpdateFolder renames a folder with the given filename
// Returns a FolderRow, and an error if renaming fails.
func (r *Repository) UpdateFolder(ctx context.Context, arg sqlc.UpdateFolderParams) (sqlc.UpdateFolderRow, error) {
	return r.queries.UpdateFolder(ctx, arg)
}

// DeleteFolder deletes the folder with the given ID.
// Returns an error if the deletion fails.
func (r *Repository) DeleteFolder(ctx context.Context, folderID uuid.UUID) error {
	return r.queries.DeleteFolder(ctx, folderID)
}

// GetBlobIDsInFolderHierarchy returns all blob IDs contained within
// the specified folder and all of its subfolders.
// Returns a slice of UUIDs and an error if the query fails.
func (r *Repository) GetBlobIDsInFolderHierarchy(ctx context.Context, folderID uuid.UUID) ([]uuid.UUID, error) {
	return r.queries.GetBlobIDsInFolderHierarchy(ctx, folderID)
}

// UpdateFolderParentFolder updates the parent folder of a folder
// according to the provided parameters.
// Returns an error if the update fails.
func (r *Repository) UpdateFolderParentFolder(ctx context.Context, arg sqlc.UpdateFolderParentFolderParams) error {
	return r.queries.UpdateFolderParentFolder(ctx, arg)
}

// ListSelectableFolders returns a list of folders that can be selected
// for moving, based on the given parameters.
// Returns a slice of ListSelectableFoldersRow and an error if the query fails.
func (r *Repository) ListSelectableFolders(ctx context.Context, args sqlc.ListSelectableFoldersParams) ([]sqlc.ListSelectableFoldersRow, error) {
	return r.queries.ListSelectableFolders(ctx, args)
}
