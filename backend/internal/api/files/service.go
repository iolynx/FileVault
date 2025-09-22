package files

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/folders"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service provides file-related operations, including uploading and managing files,
// managing file metadata, and interacting with storage and related repositories.
type Service struct {
	userRepo   *users.Repository
	folderRepo *folders.Repository
	repo       *Repository
	storage    storage.Storage
}

// NewService constructs a new Service instance with the provided repositories and storage.
func NewService(filesRepo *Repository, userRepo *users.Repository, folderRepo *folders.Repository, storage storage.Storage) *Service {
	return &Service{
		repo:       filesRepo,
		userRepo:   userRepo,
		folderRepo: folderRepo,
		storage:    storage,
	}
}

// UploadFile handles uploading a file to the storage backend and creating
// the corresponding database records. It performs ownership checks, computes
// a SHA-256 hash for deduplication, and updates blob reference counts (using a database trigger).
// Returns the created File record or an error.
func (s *Service) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folderID *uuid.UUID) (sqlc.File, error) {
	// Ownership checks
	ownerID, ok := userctx.GetUserID(ctx)
	if !ok {
		return sqlc.File{}, apierror.NewUnauthorizedError()
	}
	if folderID != nil {
		folder, err := s.repo.GetFolderByID(ctx, *folderID)
		if err != nil {
			return sqlc.File{}, apierror.NewNotFoundError("Folder")
		}
		if folder.OwnerID != ownerID {
			return sqlc.File{}, apierror.NewForbiddenError()
		}
	}

	// Compute hash (sha256)
	hasher := sha256.New()
	buf, err := io.ReadAll(io.TeeReader(file, hasher))
	if err != nil {
		return sqlc.File{}, err
	}
	sha := hex.EncodeToString(hasher.Sum(nil))

	var blob sqlc.Blob

	// Check if blob exists
	existingBlob, err := s.repo.GetBlobBySha(ctx, sha)
	if err != nil && err != pgx.ErrNoRows {
		return sqlc.File{}, apierror.NewInternalServerError("Failed to check for existing blob")
	}

	if err == nil {
		// Existing blob: update refcount
		log.Print("blob already exists, updating refcount")
		blob = existingBlob
	} else {
		user, err := s.userRepo.GetUserByID(ctx, ownerID)
		if err != nil {
			return sqlc.File{}, apierror.NewInternalServerError("Could not retrieve user data")
		}
		newBlobSize := int64(len(buf))
		log.Println("storage quota and blobsize:")
		log.Print(user.StorageQuota, newBlobSize)
		if user.StorageUsed+newBlobSize > user.StorageQuota {
			return sqlc.File{}, apierror.New(http.StatusRequestEntityTooLarge, "Storage quota exceeded")
		}

		// Upload to MinIO
		storagePath := fmt.Sprintf("%s_%s", sha, header.Filename)
		reader := bytes.NewReader(buf)
		_, err = s.storage.UploadBlob(ctx, reader, storagePath, newBlobSize, header.Header.Get("Content-Type"))
		if err != nil {
			return sqlc.File{}, err
		}
		log.Print("Uploaded Blob to storage")

		// Create blob record in DB with refcount = 0 (default). The trigger will increment it.
		blobParams := sqlc.CreateBlobParams{
			Sha256:      sha,
			StoragePath: storagePath,
			Size:        newBlobSize,
			MimeType:    util.NewText(header.Header.Get("Content-Type")),
		}
		newBlob, err := s.repo.CreateBlob(ctx, blobParams)
		if err != nil {
			return sqlc.File{}, err
		}
		log.Print("Created Blob record in db")
		blob = newBlob
	}

	fileParams := sqlc.CreateFileParams{
		OwnerID:      ownerID,
		BlobID:       blob.ID,
		Filename:     header.Filename,
		DeclaredMime: util.NewText(header.Header.Get("Content-Type")),
		Size:         blob.Size,
	}
	if folderID != nil {
		fileParams.FolderID = pgtype.UUID{Bytes: *folderID, Valid: true}
	}

	// Create the file record, which triggers blob refcount update
	log.Println("Creating file record with params:", fileParams)
	return s.repo.CreateFile(ctx, fileParams)
}

// GetFileURL returns a signed URL for accessing the file identified by fileID.
// It ensures the requesting user owns the file and fetches the corresponding blob from storage.
func (s *Service) GetFileURL(ctx context.Context, fileID uuid.UUID) (string, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return "", apierror.NewInternalServerError("Failed to get UserID")
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return "", err
	}

	// Ownership check
	if file.OwnerID != userID {
		return "", apierror.NewForbiddenError()
	}

	blob, err := s.repo.GetBlobByID(ctx, file.BlobID)
	if err != nil {
		return "", apierror.NewInternalServerError("Unable to fetch blob")
	}

	return s.storage.GetBlobURL(ctx, blob.StoragePath)
}

// GetFileByUUID retrieves a file record from the database by its UUID.
func (s *Service) GetFileByUUID(ctx context.Context, fileID uuid.UUID) (sqlc.File, error) {
	return s.repo.GetFileByUUID(ctx, fileID)
}

// DownloadFile returns a ReadCloser for the file content along with its filename.
// It checks if the user owns or has access to the file and fetches the corresponding blob.
func (s *Service) DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, string, error) {

	ownerID, ok := userctx.GetUserID(ctx)
	if !ok {
		return nil, "", apierror.NewUnauthorizedError()
	}

	log.Printf("received request from user %d to download file %s", ownerID, fileID)

	// Check if the user owns the file / is shared the file
	userHasAccess, err := s.repo.UserHasAccess(ctx, ownerID, fileID)
	if !userHasAccess || err != nil {
		log.Printf("no access")
		return nil, "", apierror.NewForbiddenError()
	}

	file, err := s.GetFileByUUID(ctx, fileID)
	if err != nil {
		return nil, "", apierror.NewInternalServerError("File not found")
	}

	blobReader, err := s.GetBlobReader(ctx, file)
	if err != nil {
		return nil, "", err
	}

	return blobReader, file.Filename, nil
}

// DeleteFile deletes a file record and its associated blob from storage if no other references exist.
// The blob record's refcount is automatically decremented and deleted through a database trigger.
// Only the owner of the file can perform this action.
func (s *Service) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return apierror.NewNotFoundError("File")
		}
		return err
	}

	// Ownership check
	if file.OwnerID != userID {
		return apierror.NewForbiddenError()
	}

	// delete the file record
	if err := s.repo.DeleteFile(ctx, fileID); err != nil {
		log.Printf("error while trying to delete file: %v", err)
		return apierror.NewInternalServerError("Failed to delete file record")
	}

	storagePath, err := s.repo.DeleteBlobIfUnused(ctx, file.BlobID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Blob %s is still referenced, not deleting from storage.", file.BlobID)
			return nil
		}
		return apierror.NewInternalServerError("Failed to clean up blob record")
	}

	log.Printf("Blob %s is unreferenced, deleting object %s from storage.", file.BlobID, storagePath)
	if err := s.storage.DeleteBlob(ctx, storagePath); err != nil {
		// Critical error: The DB record is gone, but the physical file remains.
		log.Printf("CRITICAL: Failed to delete object %s from storage: %v", storagePath, err)
	}

	log.Printf("Successfully deleted file %s", fileID)
	return nil
}

// GetBlobReader returns a ReadCloser for the blob content corresponding to the given file.
// It fetches the blob from storage using the blob's storage path.
func (s *Service) GetBlobReader(ctx context.Context, file sqlc.File) (io.ReadCloser, error) {
	blob, err := s.repo.GetBlobByID(ctx, file.BlobID)
	blobFileName := blob.StoragePath
	log.Printf("Looking up object: key=%s", blobFileName)
	obj, err := s.storage.GetBlob(ctx, blobFileName)
	if err != nil {
		return nil, err
	}

	// we dont need to check if the object exists here as the storage layer already does that for us
	return obj, nil
}

// UpdateFilename renames a file owned by the current user.
// Returns the updated FileResponse or an error if the user
// is unauthorized, forbidden, or the update fails.
func (s *Service) UpdateFilename(ctx context.Context, newFilename string, fileID uuid.UUID) (FileResponse, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return FileResponse{}, apierror.NewUnauthorizedError()
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return FileResponse{}, err
	}

	// Ownership check
	if file.OwnerID != userID {
		return FileResponse{}, apierror.NewForbiddenError()
	}

	file, err = s.repo.UpdateFilename(ctx, sqlc.UpdateFilenameParams{
		Filename: newFilename,
		ID:       fileID,
	})
	if err != nil {
		return FileResponse{}, err
	}
	return FileResponse{
		ID:            file.ID,
		Filename:      file.Filename,
		Size:          file.Size,
		ContentType:   file.DeclaredMime.String,
		UploadedAt:    file.UploadedAt.Time,
		UserOwnsFile:  file.OwnerID == userID,
		DownloadCount: &file.DownloadCount.Int64,
		ItemType:      "file",
	}, nil

}

// ListUsersWithAccesToFile returns all users who currently
// have access to a given file, this includes the owner,
// and the users the file is shared with.
// Returns a slice of User or an error if the caller lacks access.
func (s *Service) ListUsersWithAccessToFile(ctx context.Context, fileID uuid.UUID) ([]User, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return nil, apierror.NewUnauthorizedError()
	}

	file, _ := s.repo.GetFileByUUID(ctx, fileID)

	if file.OwnerID != userID {
		return nil, apierror.NewForbiddenError()
	}

	userRows, err := s.repo.ListUsersWithAccessToFile(ctx, fileID)
	if err != nil {
		return nil, apierror.NewInternalServerError()
	}

	usersWithAccess := make([]User, 0, len(userRows))
	for _, r := range userRows {
		usersWithAccess = append(usersWithAccess, User{
			ID:         r.ID,
			Name:       r.Name,
			Email:      r.Email,
			Permission: r.Permission,
		})
	}

	return usersWithAccess, nil
}

// ListContents retrieves files and folders for the authenticated user,
// either within a specified folder or at the root. It applies filters,
// pagination, and sorting as specified in the request.
func (s *Service) ListContents(ctx context.Context, req ListContentsRequest) (ListContentsResponse, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return ListContentsResponse{}, apierror.NewUnauthorizedError()
	}
	var items []ContentItem
	var totalCount int64

	log.Printf("Retrieving contents of folder: %v, with sort: %s %s", req.FolderID, req.SortBy, req.SortOrder)

	if req.FolderID != nil {
		// --- Handle Listing a Specific Folder ---
		folder, err := s.repo.GetFolderByID(ctx, *req.FolderID)
		if err != nil {
			return ListContentsResponse{}, apierror.NewNotFoundError("Folder")
		}
		if folder.OwnerID != userID {
			return ListContentsResponse{}, apierror.NewForbiddenError()
		}

		params := sqlc.ListFolderContentsParams{
			UserID:         userID,
			MimeType:       req.MimeType,
			ParentFolderID: *req.FolderID,
			Search:         req.Search,
			MinSize:        req.MinSize,
			MaxSize:        req.MaxSize,
			SortBy:         req.SortBy,
			SortOrder:      req.SortOrder,
			Limit:          req.Limit,
			Offset:         req.Offset,
			UploadedAfter:  util.ToPgTimestamptz(req.UploadedAfter),
			UploadedBefore: util.ToPgTimestamptz(req.UploadedBefore),
		}
		rows, repoErr := s.repo.ListFolderContents(ctx, params)
		if repoErr != nil {
			return ListContentsResponse{}, repoErr
		}

		if len(rows) > 0 {
			totalCount = rows[0].TotalCount
		}
		items = mapListFolderContentsRows(rows)

	} else {
		// --- Handle Listing the Root Folder ---
		params := sqlc.ListRootContentsParams{
			UserID:          userID,
			MimeType:        req.MimeType,
			Search:          req.Search,
			OwnershipStatus: req.OwnershipStatus,
			MinSize:         req.MinSize,
			MaxSize:         req.MaxSize,
			SortBy:          req.SortBy,
			SortOrder:       req.SortOrder,
			Limit:           req.Limit,
			Offset:          req.Offset,
			UploadedAfter:   util.ToPgTimestamptz(req.UploadedAfter),
			UploadedBefore:  util.ToPgTimestamptz(req.UploadedBefore),
		}
		rows, repoErr := s.repo.ListRootContents(ctx, params)
		if repoErr != nil {
			return ListContentsResponse{}, repoErr
		}

		if len(rows) > 0 {
			totalCount = rows[0].TotalCount
		}
		items = mapListRootContentsRows(rows)
	}

	response := ListContentsResponse{
		Data:       items,
		TotalCount: totalCount,
	}

	return response, nil
}

// mapListFolderContentsRows converts sqlc folder-content rows
// into standardized ContentItem structs for API responses.
func mapListFolderContentsRows(rows []sqlc.ListFolderContentsRow) []ContentItem {
	items := make([]ContentItem, len(rows))
	for i, r := range rows {
		item := ContentItem{
			ID:           r.ID,
			ItemType:     r.ItemType,
			Filename:     r.Filename,
			UploadedAt:   r.UploadedAt.Time,
			UserOwnsFile: r.UserOwnsFile,
		}
		// Safely assign nullable fields
		if r.Size.Valid {
			item.Size = &r.Size.Int64
		}
		if r.ContentType.Valid {
			item.ContentType = &r.ContentType.String
		}
		if r.DownloadCount.Valid {
			item.DownloadCount = &r.DownloadCount.Int64
		}
		items[i] = item
	}
	return items
}

// mapListRootContentsRows converts sqlc root-content rows
// into standardized ContentItem structs for API responses.
func mapListRootContentsRows(rows []sqlc.ListRootContentsRow) []ContentItem {
	items := make([]ContentItem, len(rows))
	for i, r := range rows {
		item := ContentItem{
			ID:           r.ID,
			ItemType:     r.ItemType,
			Filename:     r.Filename,
			UploadedAt:   r.UploadedAt.Time,
			UserOwnsFile: r.UserOwnsFile,
		}

		// Safely assign nullable fields
		if r.Size.Valid {
			item.Size = &r.Size.Int64
		}
		if r.ContentType.Valid {
			item.ContentType = &r.ContentType.String
		}
		if r.DownloadCount.Valid {
			item.DownloadCount = &r.DownloadCount.Int64
		}
		items[i] = item
	}
	return items
}

// IncrementDownloadCount increments the download counter
// for the given file in the database.
func (s *Service) IncrementDownloadCount(ctx context.Context, fileID uuid.UUID) error {
	return s.repo.IncrementDownloadCount(ctx, fileID)
}

// ListAllFiles retrieves all files across the system with
// pagination and sorting, primarily for admin use.
func (s *Service) ListAllFiles(ctx context.Context, limit, offset int32, sortBy, sortOrder string) ([]sqlc.ListAllFilesRow, error) {
	return s.repo.ListAllFiles(ctx, sqlc.ListAllFilesParams{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
}

// GetShareInfo returns sharing details for a file owned by
// the current user, including its share URL and the list
// of users who have been granted access.
func (s *Service) GetShareInfo(ctx context.Context, fileID uuid.UUID) (ShareInfoResponse, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return ShareInfoResponse{}, apierror.NewUnauthorizedError()
	}

	// Ownership check
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return ShareInfoResponse{}, apierror.NewNotFoundError("File")
	}
	if file.OwnerID != userID {
		return ShareInfoResponse{}, apierror.NewForbiddenError()
	}

	// Get the blob, the blob's storagePath is the share URL
	blob, err := s.repo.GetBlobByID(ctx, file.BlobID)
	if err != nil {
		return ShareInfoResponse{}, apierror.NewInternalServerError("Unable to fetch blob")
	}

	shareURL, err := s.storage.GetBlobURL(ctx, blob.StoragePath)
	if err != nil {
		return ShareInfoResponse{}, err
	}

	// Get the list of users the file is currently shared with
	sharedWithRows, err := s.repo.ListUsersWithAccessToFile(ctx, fileID)
	if err != nil {
		return ShareInfoResponse{}, err
	}

	// Get the list of all other users to share with
	allUsersRows, err := s.userRepo.ListOtherUsers(ctx, userID)
	if err != nil {
		return ShareInfoResponse{}, err
	}

	sharedWith := make([]User, 0, len(sharedWithRows))
	for _, r := range sharedWithRows {
		sharedWith = append(sharedWith, User{
			ID:         r.ID,
			Name:       r.Name,
			Email:      r.Email,
			Permission: r.Permission,
		})
	}

	allUsers := make([]User, 0, len(allUsersRows))
	for _, r := range allUsersRows {
		allUsers = append(allUsers, User{
			ID:    r.ID,
			Name:  r.Name,
			Email: r.Email,
		})
	}

	// Bundle and return response
	return ShareInfoResponse{
		ShareURL:   shareURL,
		SharedWith: sharedWith,
		AllUsers:   allUsers,
	}, nil
}

// MoveFile moves a file into a different folder. Ensures
// the caller owns both the file and the target folder (if provided).
// If the target folder is not provided, it is moved to the root Folder.
func (s *Service) MoveFile(ctx context.Context, fileID uuid.UUID, req MoveFileRequest) error {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	// Get file, ownership check
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return err
	}
	if file.OwnerID != userID {
		return apierror.NewForbiddenError()
	}

	if req.TargetFolderID != nil {
		// verify folder exists and is owned by user
		folder, err := s.folderRepo.GetFolderByID(ctx, *req.TargetFolderID)
		if err != nil {
			return err
		}
		if folder.OwnerID != userID {
			return apierror.NewForbiddenError()
		}
	}

	params := sqlc.UpdateFileFolderParams{
		ID: fileID,
	}
	if req.TargetFolderID != nil {
		params.FolderID = pgtype.UUID{Bytes: *req.TargetFolderID, Valid: true}
	}

	// update DB
	return s.repo.UpdateFileFolder(ctx, params)
}

// UpdateFileShares updates the list of users a file is shared with.
// It removes all existing shares for the file, then inserts the new list
// of user IDs in a single transaction to ensure atomicity.
func (s *Service) UpdateFileShares(ctx context.Context, req UpdateFileSharesRequest) error {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	// Ownership Check
	file, err := s.repo.GetFileByUUID(ctx, req.FileID)
	if err != nil {
		return apierror.NewNotFoundError("File")
	}
	if file.OwnerID != userID {
		return apierror.NewForbiddenError()
	}

	// Starting a database transaction
	// to perform atomic changes on the file_shares table
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return apierror.NewInternalServerError("could not start transaction")
	}
	defer tx.Rollback(ctx) // rollback on error

	// use this for all succeeding operations
	qtx := s.repo.WithTx(tx)

	// delete all existing shares for this file.
	if err := qtx.DeleteAllSharesForFile(ctx, req.FileID); err != nil {
		return apierror.NewInternalServerError("could not update shares")
	}

	// if there are new users to share with, perform a bulk insert
	// everything succeeded, commit the transaction.
	if len(req.UserIDs) > 0 {
		params := make([]sqlc.AddSharesToFileParams, len(req.UserIDs))

		for i, targetUserID := range req.UserIDs {
			params[i] = sqlc.AddSharesToFileParams{
				FileID:     req.FileID,
				SharedWith: targetUserID,
			}
		}

		_, err = qtx.AddSharesToFile(ctx, params)
		if err != nil {
			return apierror.NewInternalServerError("could not add new shares")
		}
	}
	return tx.Commit(ctx)

}
