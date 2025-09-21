package files

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// ListContentsRow is the common structure returned by our manual list queries.
// It matches the columns from your UNION queries.
type ListContentsRow struct {
	ID            uuid.UUID
	Filename      string
	ItemType      string
	Size          sql.NullInt64 // Use sql types for nullable db columns
	ContentType   sql.NullString
	UploadedAt    time.Time
	UserOwnsFile  bool
	DownloadCount sql.NullInt64
	FolderID      uuid.NullUUID
}
