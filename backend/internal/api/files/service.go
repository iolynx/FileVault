package files

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	usersRepo *users.Repository
	repo      *Repository
	storage   storage.Storage
}

type File struct {
	ID            uuid.UUID `json:"id"`
	Filename      string    `json:"filename"`
	Size          int64     `json:"size"`
	ContentType   string    `json:"content_type"`
	UploadedAt    time.Time `json:"uploaded_at"`
	UserOwnsFile  bool      `json:"user_owns_file"`
	DownloadCount *int64    `json:"download_count,omitempty"`
}

type FileResponse struct {
	ID            uuid.UUID `json:"id"`
	Filename      string    `json:"filename"`
	Size          int64     `json:"size"`
	ContentType   string    `json:"content_type"`
	UploadedAt    time.Time `json:"uploaded_at"`
	UserOwnsFile  bool      `json:"user_owns_file"`
	DownloadCount *int64    `json:"download_count,omitempty"`
	ItemType      string    `json:"item_type"`
}

type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Permission string `json:"permission"`
}

func NewService(filesRepo *Repository, usersRepo *users.Repository, storage storage.Storage) *Service {
	return &Service{
		repo:      filesRepo,
		usersRepo: usersRepo,
		storage:   storage,
	}
}

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
		user, err := s.usersRepo.GetUserByID(ctx, ownerID)
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

		// Create blob record in DB with refcount=0. The trigger will increment it.
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

	log.Println("Creating file record with params:", fileParams)
	return s.repo.CreateFile(ctx, fileParams)
}

// ListFiles returns a list of files owned by a particular User
func (s *Service) ListFilesByOwner(ctx context.Context, ownerID int64, search string, limit, offset int32) ([]File, error) {
	fileRows, err := s.repo.ListFilesByOwner(ctx, ownerID, search, limit, offset)
	if err != nil {
		return []File{}, err
	}

	files := make([]File, 0, len(fileRows))
	for _, r := range fileRows {
		files = append(files, File{
			ID:            r.ID,
			Filename:      r.Filename,
			Size:          r.Size,
			ContentType:   r.ContentType.String,
			UploadedAt:    r.UploadedAt.Time,
			UserOwnsFile:  true,
			DownloadCount: &r.DownloadCount.Int64,
		})
	}

	return files, nil
}

func (s *Service) GetFileURL(ctx context.Context, fileID uuid.UUID) (string, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return "", apierror.NewInternalServerError("Failed to get UserID")
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return "", err
	}

	if file.OwnerID != userID {
		return "", apierror.NewForbiddenError()
	}

	blob, err := s.repo.GetBlobByID(ctx, file.BlobID)
	if err != nil {
		return "", apierror.NewInternalServerError("Unable to fetch blob")
	}

	return s.storage.GetBlobURL(ctx, blob.StoragePath)
}

func (s *Service) GetFileByUUID(ctx context.Context, fileID uuid.UUID) (sqlc.File, error) {
	return s.repo.GetFileByUUID(ctx, fileID)
}

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

func (s *Service) UpdateFilename(ctx context.Context, newFilename string, fileID uuid.UUID) (FileResponse, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return FileResponse{}, apierror.NewUnauthorizedError()
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return FileResponse{}, err
	}

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

func (s *Service) ShareFile(ctx context.Context, fileID uuid.UUID, targetUserID int64) error {
	ownerID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return err
	}

	if targetUserID <= 0 {
		return apierror.NewBadRequestError("Invalid UserID")
	}
	if int64(targetUserID) == ownerID {
		return apierror.NewBadRequestError("Target User ID cannot be your ID")
	}

	if file.OwnerID != ownerID {
		return apierror.NewForbiddenError()
	}

	userExists, _ := s.repo.DoesUserExist(ctx, targetUserID)
	if !userExists {
		return apierror.NewInternalServerError("User does not exist")
	}

	_, err = s.repo.CreateFileShare(ctx, fileID, targetUserID)
	return err
}

func (s *Service) ListFilesSharedWithUser(ctx context.Context, userID int64, search string, limit, offset int32) ([]File, error) {
	fileRows, err := s.repo.ListFilesSharedWithUser(ctx, userID, search, limit, offset)
	if err != nil {
		return nil, err
	}

	files := make([]File, 0, len(fileRows))
	for _, r := range fileRows {
		files = append(files, File{
			ID:           r.ID,
			Filename:     r.Filename,
			Size:         r.Size,
			ContentType:  r.DeclaredMime.String,
			UploadedAt:   r.UploadedAt.Time,
			UserOwnsFile: false,
		})
	}

	return files, nil
}

func (s *Service) RemoveFileShare(ctx context.Context, fileID uuid.UUID, sharedWith int64) error {
	ownerID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return apierror.NewNotFoundError("File")
	}

	if file.OwnerID != ownerID {
		return apierror.NewForbiddenError()
	}

	return s.repo.DeleteFileShare(ctx, fileID, sharedWith)
}

func (s *Service) ListUsersWithAccesToFile(ctx context.Context, fileID uuid.UUID) ([]User, error) {
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

type ContentItem struct {
	ID            uuid.UUID `json:"id"`
	ItemType      string    `json:"item_type"` // "file" or "folder"
	Filename      string    `json:"filename"`
	Size          *int64    `json:"size,omitempty"`
	ContentType   *string   `json:"content_type,omitempty"`
	UploadedAt    time.Time `json:"uploaded_at"`
	UserOwnsFile  bool      `json:"user_owns_file"`
	DownloadCount *int64    `json:"download_count,omitempty"`
}

type ListContentsRequest struct {
	FolderID        *uuid.UUID    `json:"folder_id"`
	Search          string        `json:"search"`
	MimeType        string        `json:"content_type"`
	UploadedAfter   *time.Time    `json:"uploaded_after"`
	UploadedBefore  *time.Time    `json:"uploaded_before"`
	OwnershipStatus int32         `json:"user_owns_file"`
	Limit           int32         `json:"limit"`
	Offset          int32         `json:"offset"`
	MinSize         sql.NullInt64 `json:"min_size"`
	MaxSize         sql.NullInt64 `json:"max_size"`
	SortBy          string        `json:"sort_by"`
	SortOrder       string        `json:"sort_order"`
}

type ListContentsResponse struct {
	Data       []ContentItem `json:"data"`
	TotalCount int64         `json:"totalCount"`
}

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

// Helper function to map the generated sqlc rows to ContentItem,
// which is the standardized way to represent Content (files / folders)
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

// Helper function to map the generated sqlc rows to ContentItem,
// which is the standardized way to represent Content (files / folders)
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

func (s *Service) IncrementDownloadCount(ctx context.Context, fileID uuid.UUID) error {
	return s.repo.IncrementDownloadCount(ctx, fileID)
}

func (s *Service) ListAllFiles(ctx context.Context, limit, offset int32, sortBy, sortOrder string) ([]sqlc.ListAllFilesRow, error) {
	return s.repo.ListAllFiles(ctx, sqlc.ListAllFilesParams{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
}

type ShareInfoResponse struct {
	ShareURL   string `json:"shareURL"`
	SharedWith []User `json:"sharedWith"`
	AllUsers   []User `json:"allUsers"`
}

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
	allUsersRows, err := s.usersRepo.ListOtherUsers(ctx, userID)
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
