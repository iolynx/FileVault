// Package folders provides database operations for file and folder management,
// including fetching, creating, and moving folders.
package folders

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apphandler"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler is the HTTP handler for folder-related endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new folder Handler with the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all folder-related HTTP routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/folders", apphandler.MakeHTTPHandler(h.CreateFolder))
	r.Patch("/folders/{folderId}", apphandler.MakeHTTPHandler(h.UpdateFolder))
	r.Delete("/folders/{folderId}", apphandler.MakeHTTPHandler(h.DeleteFolder))
	r.Patch("/folders/{id}/move", apphandler.MakeHTTPHandler(h.MoveFolder))
	r.Get("/folders/{id}", apphandler.MakeHTTPHandler(h.GetSelectableFolders))
	r.Get("/folders/", apphandler.MakeHTTPHandler(h.GetSelectableFolders))
}

// CreateFolder handles POST /folders.
// It creates a new folder for the authenticated user.
func (h *Handler) CreateFolder(w http.ResponseWriter, r *http.Request) error {
	var req CreateFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.NewBadRequestError("Invalid request body")
	}

	folder, err := h.service.CreateFolder(r.Context(), req)
	if err != nil {
		log.Printf("Error while trying to create folder: %s", err)
		return err
	}

	return util.WriteJSON(w, http.StatusCreated, folder)
}

// UpdateFolder handles PATCH /folders/{folderId}.
// It updates the folder's metadata, such as its name.
func (h *Handler) UpdateFolder(w http.ResponseWriter, r *http.Request) error {
	folderIDStr := chi.URLParam(r, "folderId")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return apierror.NewBadRequestError("Invalid folder ID")
	}

	var req UpdateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.NewBadRequestError("Invalid request body")
	}

	updatedFolder, err := h.service.UpdateFolder(r.Context(), folderID, req)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, updatedFolder)
}

// DeleteFolder handles DELETE /folders/{folderId}.
// It deletes the specified folder along with its contents,
// including files and subfolders, recursively.
func (h *Handler) DeleteFolder(w http.ResponseWriter, r *http.Request) error {
	folderIDStr := chi.URLParam(r, "folderId")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return apierror.NewBadRequestError("Invalid folder ID")
	}

	err = h.service.DeleteFolder(r.Context(), folderID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// MoveFolder handles PATCH /folders/{id}/move.
// It updates the parent folder of a folder, effectively moving it.
func (h *Handler) MoveFolder(w http.ResponseWriter, r *http.Request) error {
	folderIDStr := chi.URLParam(r, "id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return apierror.NewBadRequestError("Invalid fileID")
	}

	var req UpdateFolderParentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.NewBadRequestError("Invalid request body")
	}

	return h.service.UpdateFolderParent(r.Context(), folderID, req)
}

// GetSelectableFolders handles GET /folders/{id} and GET /folders/.
// It returns folders that the authenticated user move their folder to.
// This includes the set of all folders the user owns, with the exception
// of the folder itself and its parent (if any).
// The GET/folders/ endpoint is used when the parent folder is null (root).
func (h *Handler) GetSelectableFolders(w http.ResponseWriter, r *http.Request) error {
	folderIDStr := chi.URLParam(r, "id")
	var folderID *uuid.UUID
	if folderIDStr != "" {
		parsed, err := uuid.Parse(folderIDStr)
		if err != nil {
			return apierror.NewBadRequestError("Invalid FolderID")
		}
		folderID = &parsed
	}

	folders, err := h.service.GetSelectableFolders(r.Context(), folderID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, folders)
}
