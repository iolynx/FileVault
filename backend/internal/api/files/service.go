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

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	repo    *Repository
	storage storage.Storage
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

	// Create file record referencing blob
	return s.repo.CreateFile(ctx, ownerID, newBlob.ID, header.Filename, header.Header.Get("Content-Type"), int64(len(buf)))
}

func (s *Service) ListFiles(ctx context.Context, ownerID int64) ([]sqlc.File, error) {
	return s.repo.ListFilesByOwner(ctx, ownerID)
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
		return errors.New("cannot delete file owned by another user")
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
