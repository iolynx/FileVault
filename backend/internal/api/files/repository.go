package files

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

// BeginTx starts a new database transaction.
func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

// WithTx returns a new repository instance with its queries scoped to the provided transaction.
func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	// We return a pointer to a new repository struct
	return &Repository{
		pool:    r.pool,
		queries: r.queries.WithTx(tx),
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

func (r *Repository) CreateBlob(ctx context.Context, arg sqlc.CreateBlobParams) (sqlc.Blob, error) {
	return r.queries.CreateBlob(ctx, arg)
}

func (r *Repository) IncrementBlobRefcount(ctx context.Context, blobID uuid.UUID) (int32, error) {
	return r.queries.IncrementBlobRefcount(ctx, blobID)
}

func (r *Repository) DecrementBlobRefCount(ctx context.Context, blobID uuid.UUID) (int32, error) {
	return r.queries.DecrementBlobRefcount(ctx, blobID)
}

func (r *Repository) CreateFile(ctx context.Context, arg sqlc.CreateFileParams) (sqlc.File, error) {
	return r.queries.CreateFile(ctx, arg)
}

func (r *Repository) ListFilesByOwner(ctx context.Context, ownerID int64, search string, limit, offset int32) ([]sqlc.ListFilesByOwnerRow, error) {
	return r.queries.ListFilesByOwner(ctx, sqlc.ListFilesByOwnerParams{
		OwnerID: ownerID,
		Column2: search,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *Repository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	return r.queries.DeleteFile(ctx, fileID)
}

func (r *Repository) ListFilesSharedWithUser(ctx context.Context, userID int64, search string, limit, offset int32) ([]sqlc.File, error) {
	return r.queries.ListFilesSharedWithUser(ctx, sqlc.ListFilesSharedWithUserParams{
		SharedWith: userID,
		Column2:    search,
		Limit:      limit,
		Offset:     offset,
	})
}

func (r *Repository) UpdateFilename(ctx context.Context, arg sqlc.UpdateFilenameParams) (sqlc.File, error) {
	return r.queries.UpdateFilename(ctx, arg)
}

func (r *Repository) DoesUserExist(ctx context.Context, userID int64) (bool, error) {
	return r.queries.UserExists(ctx, userID)
}

func (r *Repository) ListUsersWithAccessToFile(ctx context.Context, fileID uuid.UUID) ([]sqlc.ListUsersWithAccessToFileRow, error) {
	return r.queries.ListUsersWithAccessToFile(ctx, fileID)
}

// this function takes the userID in its SharedWith Param, but also checks if the user owns the file
func (r *Repository) UserHasAccess(ctx context.Context, sharedWith int64, fileID uuid.UUID) (bool, error) {
	return r.queries.UserHasAccess(ctx, sqlc.UserHasAccessParams{
		SharedWith: sharedWith,
		ID:         fileID,
	})
}

func (r *Repository) ListFilesForUser(ctx context.Context, userID int64, filename string, ownershipStatus int32, mimeType string, uploadedAfter, uploadedBefore pgtype.Timestamptz, limit, offset int32) ([]sqlc.ListFilesForUserRow, error) {

	return r.queries.ListFilesForUser(ctx, sqlc.ListFilesForUserParams{
		UserID:          userID,
		Filename:        filename,
		MimeType:        mimeType,
		UploadedAfter:   uploadedAfter,
		UploadedBefore:  uploadedBefore,
		OwnershipStatus: ownershipStatus,
		Limit:           limit,
		Offset:          offset,
	})
}

func (r *Repository) ListFolderContents(ctx context.Context, arg sqlc.ListFolderContentsParams) ([]sqlc.ListFolderContentsRow, error) {
	return r.queries.ListFolderContents(ctx, arg)
}

func (r *Repository) ListRootContents(ctx context.Context, arg sqlc.ListRootContentsParams) ([]sqlc.ListRootContentsRow, error) {
	return r.queries.ListRootContents(ctx, arg)
}

func (r *Repository) IncrementDownloadCount(ctx context.Context, fileID uuid.UUID) error {
	return r.queries.IncrementFileDownloadCount(ctx, fileID)
}

func (r *Repository) GetFolderByID(ctx context.Context, folderID uuid.UUID) (sqlc.Folder, error) {
	return r.queries.GetFolderByID(ctx, folderID)
}

func (r *Repository) ListAllFiles(ctx context.Context, arg sqlc.ListAllFilesParams) ([]sqlc.ListAllFilesRow, error) {
	return r.queries.ListAllFiles(ctx, arg)
}

func (r *Repository) DeleteBlobIfUnused(ctx context.Context, blobID uuid.UUID) (string, error) {
	return r.queries.DeleteBlobIfUnused(ctx, blobID)
}

func (r *Repository) UpdateFileFolder(ctx context.Context, arg sqlc.UpdateFileFolderParams) error {
	return r.queries.UpdateFileFolder(ctx, arg)
}

func (r *Repository) DeleteAllSharesForFile(ctx context.Context, fileID uuid.UUID) error {
	return r.queries.DeleteAllSharesForFile(ctx, fileID)
}

func (r *Repository) AddSharesToFile(ctx context.Context, arg []sqlc.AddSharesToFileParams) (int64, error) {
	return r.queries.AddSharesToFile(ctx, arg)
}
