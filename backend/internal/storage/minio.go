package storage

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioStorage(cfg config.MinioConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Access, cfg.Secret, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Created bucket %s", cfg.Bucket)
	}

	return &MinioStorage{Client: client, BucketName: cfg.Bucket}, nil
}

func (m *MinioStorage) UploadBlob(ctx context.Context, r io.Reader, fileName string, size int64, contentType string) (string, error) {
	log.Print("Uploading: ", contentType, fileName, size)
	if contentType == "" {
		contentType = "application/octet-stream"
		size = -1
	}

	info, err := m.Client.PutObject(ctx, m.BucketName, fileName, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	log.Printf("Uploaded %d bytes\n", info.Size)

	if err != nil {
		return "", err
	}
	return fileName, nil
}

// Get the public URL for the object. Here fileName represents the object name and not the actual file's name
func (m *MinioStorage) GetBlobURL(ctx context.Context, fileName string) (string, error) {
	expiry := time.Minute * 15
	url, err := m.Client.PresignedGetObject(ctx, m.BucketName, fileName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *MinioStorage) GetBlob(ctx context.Context, fileName string) (io.ReadCloser, error) {

	obj, err := m.Client.GetObject(ctx, m.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	if _, err := obj.Stat(); err != nil {
		return nil, err
	}
	log.Print("done")

	return obj, nil
}

func (m *MinioStorage) DeleteBlob(ctx context.Context, fileName string) error {
	return m.Client.RemoveObject(ctx, m.BucketName, fileName, minio.RemoveObjectOptions{})
}
