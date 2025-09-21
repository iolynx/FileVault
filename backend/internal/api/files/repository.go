package files

import (
	"context"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db *sqlc.Queries) *Repository {
	return &Repository{
		queries: db,
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

func (r *Repository) CreateFileShare(ctx context.Context, fileID uuid.UUID, targetUserID int64) (sqlc.FileShare, error) {
	return r.queries.CreateFileShare(ctx, sqlc.CreateFileShareParams{
		FileID:     fileID,
		SharedWith: targetUserID,
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

func (r *Repository) DeleteFileShare(ctx context.Context, fileID uuid.UUID, sharedWith int64) error {
	return r.queries.DeleteFileShare(ctx, sqlc.DeleteFileShareParams{
		FileID:     fileID,
		SharedWith: sharedWith,
	})
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

func (r *Repository) IncrementUserStorage(ctx context.Context, userID int64, original_storage_increment, dedup_storage_increment int) error {
	return r.queries.IncrementUserStorage(ctx, sqlc.IncrementUserStorageParams{
		ID:                   userID,
		OriginalStorageBytes: int64(original_storage_increment),
		DedupStorageBytes:    int64(dedup_storage_increment),
	})
}

func (r *Repository) DecrementUserStorage(ctx context.Context, userID int64, original_storage_decrement, dedup_storage_decrement int) error {
	return r.queries.DecrementUserStorage(ctx, sqlc.DecrementUserStorageParams{
		ID:                   userID,
		OriginalStorageBytes: int64(original_storage_decrement),
		DedupStorageBytes:    int64(dedup_storage_decrement),
	})
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

func (r *Repository) DeleteBlobIfUnused(ctx context.Context, blobID uuid.UUID) error {
	return r.queries.DeleteBlobIfUnused(ctx, blobID)
}
