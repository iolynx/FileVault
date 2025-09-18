package files

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: sqlc.New(pool),
	}
}

func (r *Repository) GetBlobByID(ctx context.Context, id uuid.UUID) (sqlc.Blob, error) {
	return r.queries.GetBlobByID(ctx, id)
}

func (r *Repository) GetBlobBySha(ctx context.Context, sha string) (sqlc.Blob, error) {
	return r.queries.GetBlobBySha(ctx, sha)
}

func (r *Repository) GetFileByUUID(ctx context.Context, id uuid.UUID) (sqlc.File, error) {
	return r.queries.GetFileByUUID(ctx, id)
}

func (r *Repository) CreateBlob(ctx context.Context, sha, storagePath string, mimeType string, size int64) (sqlc.Blob, error) {
	return r.queries.CreateBlob(ctx, sqlc.CreateBlobParams{
		Sha256:      sha,
		StoragePath: storagePath,
		MimeType:    util.NewText(mimeType),
		Size:        size,
		Refcount:    1,
	})
}

func (r *Repository) IncrementBlobRefcount(ctx context.Context, blobID uuid.UUID) (int32, error) {
	return r.queries.IncrementBlobRefcount(ctx, blobID)
}

func (r *Repository) DecrementBlobRefCount(ctx context.Context, blobID uuid.UUID) (int32, error) {
	return r.queries.DecrementBlobRefcount(ctx, blobID)
}

func (r *Repository) CreateFile(ctx context.Context, ownerID int64, blobID uuid.UUID, filename string, declaredMime string, size int64) (sqlc.File, error) {
	return r.queries.CreateFile(ctx, sqlc.CreateFileParams{
		OwnerID:      ownerID,
		BlobID:       blobID,
		Filename:     filename,
		DeclaredMime: util.NewText(declaredMime),
		Size:         size,
	})
}

func (r *Repository) ListFilesByOwner(ctx context.Context, ownerID int64, search string, limit, offset int32) ([]sqlc.ListFilesByOwnerRow, error) {
	return r.queries.ListFilesByOwner(ctx, sqlc.ListFilesByOwnerParams{
		OwnerID: ownerID,
		Column2: search,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *Repository) CreateFileShare(ctx context.Context, fileID uuid.UUID, targetUserID int64) (sqlc.FileShare, error) {
	return r.queries.CreateFileShare(ctx, sqlc.CreateFileShareParams{
		FileID:     fileID,
		SharedWith: targetUserID,
	})
}
