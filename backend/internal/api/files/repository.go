package files

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations related to files
type Repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewRepository creates a new Repository instance with the provided database queries.
// This Repository can be used to perform file related database operations.
// It initializes with  *pgxpool.Pool instead of sqlc.Queries in order to perform database transactions.
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

// GetBlobByID retrieves a blob record by its UUID.
// Returns an error if no blob is found.
func (r *Repository) GetBlobByID(ctx context.Context, id uuid.UUID) (sqlc.Blob, error) {
	return r.queries.GetBlobByID(ctx, id)
}

// GetBlobBySha retrieves a blob record by its SHA checksum.
// Returns an error if no blob is found.
func (r *Repository) GetBlobBySha(ctx context.Context, sha string) (sqlc.Blob, error) {
	return r.queries.GetBlobBySha(ctx, sha)
}

// GetFileByUUID retrieves a file record by its UUID.
// Returns an error if no file is found.
func (r *Repository) GetFileByUUID(ctx context.Context, id uuid.UUID) (sqlc.File, error) {
	return r.queries.GetFileByUUID(ctx, id)
}

// CreateBlob inserts a new blob record into the database.
// Returns the created blob or an error if the operation fails
func (r *Repository) CreateBlob(ctx context.Context, arg sqlc.CreateBlobParams) (sqlc.Blob, error) {
	return r.queries.CreateBlob(ctx, arg)
}

// CreateFile inserts a new file record into the database.
// Returns the created file or an error if the operation fails.
func (r *Repository) CreateFile(ctx context.Context, arg sqlc.CreateFileParams) (sqlc.File, error) {
	return r.queries.CreateFile(ctx, arg)
}

// DeleteFile removes a file record by its UUID.
// Returns an error if the deletion fails.
func (r *Repository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	return r.queries.DeleteFile(ctx, fileID)
}

// ListFilesSharedWithUser returns a paginated list of files shared with a specific user.
// Supports optional search by filename. Returns an error if the query fails.
func (r *Repository) ListFilesSharedWithUser(ctx context.Context, userID int64, search string, limit, offset int32) ([]sqlc.File, error) {
	return r.queries.ListFilesSharedWithUser(ctx, sqlc.ListFilesSharedWithUserParams{
		SharedWith: userID,
		Column2:    search,
		Limit:      limit,
		Offset:     offset,
	})
}

// UpdateFilename updates the filename of a file record.
// Returns the updated file or an error if the operation fails.
func (r *Repository) UpdateFilename(ctx context.Context, arg sqlc.UpdateFilenameParams) (sqlc.File, error) {
	return r.queries.UpdateFilename(ctx, arg)
}

// ListUsersWithAccessToFile returns all users who have access to a specific file.
// Returns an error if the query fails.
func (r *Repository) ListUsersWithAccessToFile(ctx context.Context, fileID uuid.UUID) ([]sqlc.ListUsersWithAccessToFileRow, error) {
	return r.queries.ListUsersWithAccessToFile(ctx, fileID)
}

// UserHasAccess checks if user owns the file / is shared the file
// It takes the userID in its SharedWith Param, and returns a boolean value and an error if the query fails
func (r *Repository) UserHasAccess(ctx context.Context, sharedWith int64, fileID uuid.UUID) (bool, error) {
	return r.queries.UserHasAccess(ctx, sqlc.UserHasAccessParams{
		SharedWith: sharedWith,
		ID:         fileID,
	})
}

// ListFolderContents returns the contents of a folder, sorted and paginated, with filters (including search) applied
// It returns a slice of ListFolderContentsRow, and an error if the query fails
func (r *Repository) ListFolderContents(ctx context.Context, arg sqlc.ListFolderContentsParams) ([]sqlc.ListFolderContentsRow, error) {
	return r.queries.ListFolderContents(ctx, arg)
}

// ListRootContents returns the contents of the root folder, sorted and paginated, with filters (including search) applied
// It returns a slice of ListRootContentsRow, and an error if the query fails
// Note that this function is simply ListFolderContents but for when the folder = nil (root folder)
func (r *Repository) ListRootContents(ctx context.Context, arg sqlc.ListRootContentsParams) ([]sqlc.ListRootContentsRow, error) {
	return r.queries.ListRootContents(ctx, arg)
}

// IncrementDownloadCount increments (by 1) the download count of the file in the database
// it returns an error if the query fails
func (r *Repository) IncrementDownloadCount(ctx context.Context, fileID uuid.UUID) error {
	return r.queries.IncrementFileDownloadCount(ctx, fileID)
}

// GetFolderByID retrieves a folder record by its UUID.
// Returns an error if no folder is found.
func (r *Repository) GetFolderByID(ctx context.Context, folderID uuid.UUID) (sqlc.Folder, error) {
	return r.queries.GetFolderByID(ctx, folderID)
}

// ListAllFiles returns all file records matching the given parameters.
// Supports filtering, pagination, or other criteria via ListAllFilesParams.
// Used for Admin Routes
func (r *Repository) ListAllFiles(ctx context.Context, arg sqlc.ListAllFilesParams) ([]sqlc.ListAllFilesRow, error) {
	return r.queries.ListAllFiles(ctx, arg)
}

// DeleteBlobIfUnused deletes a blob if its reference count is zero.
// Returns the SHA of the deleted blob or an error if deletion fails.
func (r *Repository) DeleteBlobIfUnused(ctx context.Context, blobID uuid.UUID) (string, error) {
	return r.queries.DeleteBlobIfUnused(ctx, blobID)
}

// UpdateFileFolder updates the parent folder of a file.
// Returns an error if the operation fails.
func (r *Repository) UpdateFileFolder(ctx context.Context, arg sqlc.UpdateFileFolderParams) error {
	return r.queries.UpdateFileFolder(ctx, arg)
}

// DeleteAllSharesForFile removes all sharing records for a given file.
// The operation is performed atomically within a transaction.
// Returns an error if the operation fails.
func (r *Repository) DeleteAllSharesForFile(ctx context.Context, fileID uuid.UUID) error {
	return r.queries.DeleteAllSharesForFile(ctx, fileID)
}

// AddSharesToFile adds new share records for a file to multiple users.
// The operation is performed atomically within a transaction.
// Returns the number of shares successfully added or an error.
func (r *Repository) AddSharesToFile(ctx context.Context, arg []sqlc.AddSharesToFileParams) (int64, error) {
	return r.queries.AddSharesToFile(ctx, arg)
}
