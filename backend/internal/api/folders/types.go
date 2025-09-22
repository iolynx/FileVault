package folders

import (
	"time"

	"github.com/google/uuid"
)

// CreateFolderRequest represents the JSON payload for creating a folder.
type CreateFolderRequest struct {
	Name           string     `json:"name"`
	ParentFolderID *uuid.UUID `json:"parent_folder_id"`
}

// UpdateFolderRequest represents the JSON payload for updating a folder.
type UpdateFolderRequest struct {
	Name string `json:"name"`
}

// Folder represents a folder in the system.
// It contains metadata about the folder, including its name, creation time,
// and optional parent folder reference (a value of nil refers to no parent, i.e. Root)
type Folder struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	CreatedAt      time.Time  `json:"created_at"`
	ParentFolderID *uuid.UUID `json:"parent_folder_id,omitempty"`
}

// FolderResponse is the standard response for folder operations.
// It embeds ItemType for compatibility with frontend component types.
type FolderResponse struct {
	ID           uuid.UUID `json:"id"`
	Filename     string    `json:"filename"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UserOwnsFile bool      `json:"user_owns_file"`
	ItemType     string    `json:"item_type"`
}

// UpdateFolderParentRequest represents a request to move a folder to a new parent folder.
// TargetFolderID is optional; if nil, the folder will be moved to the root level.
type UpdateFolderParentRequest struct {
	TargetFolderID *uuid.UUID `json:"target_folder_id"`
}
