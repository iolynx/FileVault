package files

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	repo    *Repository
	storage storage.Storage
}

type File struct {
	ID           uuid.UUID `json:"id"`
	Filename     string    `json:"filename"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UserOwnsFile bool      `json:"user_owns_file"`
}

type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Permission string `json:"permission"`
}

func NewService(repo *Repository, storage storage.Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

func (s *Service) UploadFile(ctx context.Context, ownerID int64, file multipart.File, header *multipart.FileHeader) (sqlc.File, error) {
	// Compute hash (sha256)
	hasher := sha256.New()
	buf, err := io.ReadAll(io.TeeReader(file, hasher))
	if err != nil {
		return sqlc.File{}, err
	}
	sha := hex.EncodeToString(hasher.Sum(nil))

	// Check if blob exists
	blob, err := s.repo.GetBlobBySha(ctx, sha)
	if err == nil {
		// Existing blob: update refcount
		log.Print("blob already exists, updating refcount")
		if _, err := s.repo.IncrementBlobRefcount(ctx, blob.ID); err != nil {
			return sqlc.File{}, err
		}

		// Update Storage Value for User
		err = s.repo.IncrementUserStorage(ctx, ownerID, int(blob.Size), 0)
		if err != nil {
			return sqlc.File{}, err
		}

		// Create file record pointing to existing blob
		return s.repo.CreateFile(ctx, ownerID, blob.ID, header.Filename, header.Header.Get("Content-Type"), blob.Size)
	}

	// New blob: upload to storage with this objectKey
	objectKey := fmt.Sprintf("%s_%s", sha, header.Filename)
	reader := bytes.NewReader(buf)
	_, err = s.storage.UploadBlob(ctx, reader, objectKey, int64(len(buf)), header.Header.Get("Content-Type"))
	if err != nil {
		return sqlc.File{}, err
	}
	log.Print("uploaded blob")

	// Create blob record
	newBlob, err := s.repo.CreateBlob(ctx, sha, objectKey, header.Header.Get("Content-Type"), int64(len(buf)))
	if err != nil {
		return sqlc.File{}, err
	}
	log.Print("Created Blob record in db")

	// Update storage value for user
	s.repo.IncrementUserStorage(ctx, ownerID, int(newBlob.Size), int(newBlob.Size))

	// Create file record referencing blob
	return s.repo.CreateFile(ctx, ownerID, newBlob.ID, header.Filename, header.Header.Get("Content-Type"), int64(len(buf)))
}

// ListFiles returns a list of Files (Special Object Type that does not contain all fields) owned by a particular User
func (s *Service) ListFilesByOwner(ctx context.Context, ownerID int64, search string, limit, offset int32) ([]File, error) {
	fileRows, err := s.repo.ListFilesByOwner(ctx, ownerID, search, limit, offset)
	if err != nil {
		return []File{}, err
	}

	files := make([]File, 0, len(fileRows))
	for _, r := range fileRows {
		files = append(files, File{
			ID:           r.ID,
			Filename:     r.Filename,
			Size:         r.Size,
			ContentType:  r.ContentType.String,
			UploadedAt:   r.UploadedAt.Time,
			UserOwnsFile: true,
		})
	}

	return files, nil
}

func (s *Service) GetFileURL(ctx context.Context, fileID uuid.UUID, userID int64) (string, error) {
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return "", err
	}

	// for now, the file is only viewable by its owner
	if file.OwnerID != userID {
		return "", errors.New("Unauthorized")
	}

	blob, err := s.repo.GetBlobByID(ctx, file.BlobID)
	if err != nil {
		return "", errors.New("unable to fetch blob")
	}

	return s.storage.GetBlobURL(ctx, blob.StoragePath)
}

func (s *Service) DeleteFile(ctx context.Context, fileID uuid.UUID, userID int64) error {
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return err
	}

	// ownership check
	if file.OwnerID != userID {
		return errors.New("Cannot delete file owned by another user")
	}

	// delete the file record
	if err := s.repo.queries.DeleteFile(ctx, fileID); err != nil {
		return err
	}

	// decrement blob refcount
	refCount, err := s.repo.DecrementBlobRefCount(ctx, file.BlobID)
	if err != nil {
		return err
	}

	// Unreferenced blob: delete blob, db row
	if refCount == 0 {
		log.Print("refcount is 0, deleting blob")
		blob, err := s.repo.GetBlobByID(ctx, file.BlobID)
		if err != nil {
			return err
		}
		if err := s.storage.DeleteBlob(ctx, blob.StoragePath); err != nil {
			return err
		}
		if err := s.repo.queries.DeleteBlobIfUnused(ctx, file.BlobID); err != nil {
			return err
		}

		// Decrement storage value for user
		err = s.repo.DecrementUserStorage(ctx, userID, int(blob.Size), int(blob.Size))
		if err != nil {
			return err
		}
	} else {
		// Decrement storage value for user
		err = s.repo.DecrementUserStorage(ctx, userID, int(file.Size), 0)
		if err != nil {
			return err
		}
	}

	log.Print("deleted file")
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

func (s *Service) UpdateFilename(ctx context.Context, newFilename string, fileID uuid.UUID, ownerID int64) error {
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return err
	}

	if file.OwnerID != ownerID {
		return errors.New("Unauthorized")
	}

	return s.repo.queries.UpdateFilename(ctx, sqlc.UpdateFilenameParams{
		Filename: newFilename,
		ID:       fileID,
	})
}

func (s *Service) ShareFile(ctx context.Context, fileID uuid.UUID, ownerID, targetUserID int64) error {
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return err
	}

	if file.OwnerID != ownerID {
		return errors.New("Not Authorized")
	}

	userExists, _ := s.repo.DoesUserExist(ctx, targetUserID)
	if !userExists {
		return errors.New("User does not exist")
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

func (s *Service) RemoveFileShare(ctx context.Context, fileID uuid.UUID, ownerID, sharedWith int64) error {
	file, err := s.repo.GetFileByUUID(ctx, fileID)
	if err != nil {
		return err
	}

	if file.OwnerID != ownerID {
		return errors.New("not authorized to unshare this file")
	}

	return s.repo.DeleteFileShare(ctx, fileID, sharedWith)
}

func (s *Service) ListUsersWithAccesToFile(ctx context.Context, fileID uuid.UUID, ownerID int64) ([]User, error) {
	file, _ := s.repo.GetFileByUUID(ctx, fileID)

	if file.OwnerID != ownerID {
		return nil, errors.New("Unauthorized")
	}

	userRows, err := s.repo.ListUsersWithAccessToFile(ctx, fileID)
	if err != nil {
		return nil, err
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

func (s *Service) ListFilesForUser(ctx context.Context, userID int64, search string, limit, offset int32) ([]File, error) {
	fileRows, err := s.repo.ListFilesForUser(ctx, userID, search, limit, offset)
	if err != nil {
		return nil, err
	}

	files := make([]File, 0, len(fileRows))
	for _, r := range fileRows {
		files = append(files, File{
			ID:           r.ID,
			Filename:     r.Filename,
			Size:         r.Size,
			ContentType:  r.ContentType.String,
			UploadedAt:   r.UploadedAt.Time,
			UserOwnsFile: r.UserOwnsFile,
		})
	}

	return files, nil
}
