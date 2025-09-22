package files

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// ListContentsRow is the common structure returned by the manual list queries.
// It combines both files and folders into a single structure to be returned to the frontend.
type ListContentsRow struct {
	ID            uuid.UUID
	Filename      string
	ItemType      string
	Size          sql.NullInt64
	ContentType   sql.NullString
	UploadedAt    time.Time
	UserOwnsFile  bool
	DownloadCount sql.NullInt64
	FolderID      uuid.NullUUID
}

// File represents a file entity in the system with metadata and ownership info.
type File struct {
	ID            uuid.UUID `json:"id"`
	Filename      string    `json:"filename"`
	Size          int64     `json:"size"`
	ContentType   string    `json:"content_type"`
	UploadedAt    time.Time `json:"uploaded_at"`
	UserOwnsFile  bool      `json:"user_owns_file"`
	DownloadCount *int64    `json:"download_count,omitempty"`
}

// FileResponse represents the API response for a file, including its type (file or folder)
// to simplify handling in frontend components.
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

// User represents a user who has access to files, including permissions.
type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Permission string `json:"permission"`
}

// ContentItem is a standardized representation of a content node,
// which may be a file or a folder, including metadata useful to the frontend.
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

// ListContentsRequest defines the filter, pagination, and sort
// parameters accepted when listing folder or root contents.
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

// ListContentsResponse wraps the list of ContentItems returned
// along with the total count for pagination support.
type ListContentsResponse struct {
	Data       []ContentItem `json:"data"`
	TotalCount int64         `json:"totalCount"`
}

// ShareInfoResponse represents sharing details for a file,
// including its share URL, current shared users, and all
// available users the file can be shared with.
// Primarily used to populate the Share Modal.
type ShareInfoResponse struct {
	ShareURL   string `json:"shareURL"`
	SharedWith []User `json:"sharedWith"`
	AllUsers   []User `json:"allUsers"`
}

// UpdateFileSharesRequest represents a request to update
// the users a file is shared with, replacing any existing shares.
type UpdateFileSharesRequest struct {
	FileID  uuid.UUID
	UserIDs []int64
}

// UpdateFilenameRequest represents a request to update
// the filename of a file belonging to the user.
type UpdateFilenameRequest struct {
	Filename string `json:"name"`
}

// ListAllFilesResponse represents metadata for a single file,
// including ownership, size, MIME type, and download count.
type ListAllFilesResponse struct {
	ID            string    `json:"id"`
	Filename      string    `json:"filename"`
	Size          int64     `json:"size"`
	DeclaredMime  string    `json:"declared_mime"`
	UploadedAt    time.Time `json:"uploaded_at"`
	DownloadCount *int64    `json:"download_count,omitempty"`
	OwnerID       int64     `json:"owner_id"`
	OwnerEmail    string    `json:"owner_email"`
}

// PaginatedFilesResponse represents a paginated response aggregating
// content data (files & folders), bundled with total count for pagination.
type PaginatedFilesResponse struct {
	Data       []ListAllFilesResponse `json:"data"`
	TotalCount int64                  `json:"totalCount"`
}

// MoveFileRequest represents the default request to
// move a file to a target folder, represented by its ID
type MoveFileRequest struct {
	TargetFolderID *uuid.UUID `json:"target_folder_id"`
}

// updateSharesPayload represents the JSON payload used to update
// the list of user IDs a file is shared with.
type updateSharesPayload struct {
	UserIDs []int64 `json:"user_ids"`
}
