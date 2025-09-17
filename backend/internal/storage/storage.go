package storage

import (
	"context"
	"io"
)

type Storage interface {
	UploadBlob(ctx context.Context, r io.Reader, fileName string, size int64, contentType string) (string, error)
	GetBlob(ctx context.Context, fileName string) (io.ReadCloser, error)
	GetBlobURL(ctx context.Context, fileName string) (string, error)
	DeleteBlob(ctx context.Context, fileName string) error
}
