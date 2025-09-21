package folders

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/google/uuid"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db *sqlc.Queries) *Repository {
	return &Repository{
		queries: db,
	}
}

func (r *Repository) CreateFolder(ctx context.Context, arg sqlc.CreateFolderParams) (sqlc.Folder, error) {
	return r.queries.CreateFolder(ctx, arg)
}

func (r *Repository) GetFolderByID(ctx context.Context, folderID uuid.UUID) (sqlc.Folder, error) {
	return r.queries.GetFolderByID(ctx, folderID)
}

func (r *Repository) UpdateFolder(ctx context.Context, arg sqlc.UpdateFolderParams) (sqlc.UpdateFolderRow, error) {
	return r.queries.UpdateFolder(ctx, arg)
}

func (r *Repository) DeleteFolder(ctx context.Context, folderID uuid.UUID) error {
	return r.queries.DeleteFolder(ctx, folderID)
}

func (r *Repository) GetObjectKeysInFolderHierarchy(ctx context.Context, folderID uuid.UUID) ([]string, error) {
	return r.queries.GetObjectKeysInFolderHierarchy(ctx, folderID)
}
