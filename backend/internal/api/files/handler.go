package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apphandler"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FileHandler struct {
	service *Service
}

func NewFileHandler(service *Service) *FileHandler {
	return &FileHandler{service: service}
}

func (h *FileHandler) RegisterRoutes(r chi.Router) {
	r.Post("/files/upload", apphandler.MakeHTTPHandler(h.Upload))

	r.Get("/files", apphandler.MakeHTTPHandler(h.ListContents))
	r.Get("/files/url/{id}", apphandler.MakeHTTPHandler(h.GetURL))

	r.Get("/files/{id}", apphandler.MakeHTTPHandler(h.DownloadFile))
	r.Patch("/files/{id}", apphandler.MakeHTTPHandler(h.UpdateFilename))
	r.Delete("/files/{id}", apphandler.MakeHTTPHandler(h.DeleteFile))
	r.Patch("/files/{id}/move", apphandler.MakeHTTPHandler(h.MoveFile))

	r.Get("/files/{id}/share-info", apphandler.MakeHTTPHandler(h.GetShareInfo))
	r.Delete("/files/{id}/share/{userid}", apphandler.MakeHTTPHandler(h.UnshareFile))
	r.Post("/files/{id}/share", apphandler.MakeHTTPHandler(h.ShareFile))
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		return apierror.NewBadRequestError("Could not parse form")
	}

	log.Print("Parsed multipart form")

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		return apierror.NewBadRequestError("No files uploaded")
	}

	var folderID *uuid.UUID
	folderIDStr := r.FormValue("folder_id")
	if folderIDStr != "" {
		if parsedUUID, err := uuid.Parse(folderIDStr); err == nil {
			folderID = &parsedUUID
		}
	}

	for _, header := range files {
		log.Printf("Processing file: %s", header.Filename)
		file, err := header.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v", header.Filename, err)
			return apierror.NewInternalServerError("Failed to process file")
		}
		// Use defer inside the loop to ensure each file is closed
		defer file.Close()

		// Call UploadFile service for each individual file
		_, err = h.service.UploadFile(r.Context(), file, header, folderID)
		if err != nil {
			log.Printf("Upload failed for file %s: %v", header.Filename, err)
			return err
		}
	}

	log.Printf("Successfully uploaded %d files", len(files))
	return util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Successfully uploaded %d file(s)", len(files)),
	})
}

func (h *FileHandler) GetURL(w http.ResponseWriter, r *http.Request) error {
	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		return apierror.NewBadRequestError("Missing file UUID")
	}

	url, err := h.service.GetFileURL(r.Context(), uuid.MustParse(fileID))
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{
		"url": url,
	})
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	fileID := uuid.MustParse(chi.URLParam(r, "id"))

	blobReader, filename, err := h.service.DownloadFile(ctx, fileID)
	if err != nil {
		log.Printf("Error while reading blob: %s", err)
		return apierror.NewInternalServerError("Cannot read file")
	}
	defer blobReader.Close()

	// increment download count for file
	h.service.IncrementDownloadCount(ctx, fileID)

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, blobReader)
	return err
}

func (h *FileHandler) ListContents(w http.ResponseWriter, r *http.Request) error {
	req := ListContentsRequest{
		Search:   r.URL.Query().Get("search"),
		MimeType: r.URL.Query().Get("content_type"),
		Limit:    util.ParseInt32OrDefault(r.URL.Query().Get("limit"), 20),
		Offset:   util.ParseInt32OrDefault(r.URL.Query().Get("offset"), 0),
	}

	if folderID := r.URL.Query().Get("folder_id"); folderID != "" {
		f := uuid.MustParse(folderID)
		req.FolderID = &f
	}

	if before := r.URL.Query().Get("uploaded_before"); before != "" {
		if t, err := time.Parse(time.RFC3339, before); err == nil {
			req.UploadedBefore = &t
		}
	}

	if after := r.URL.Query().Get("uploaded_after"); after != "" {
		if t, err := time.Parse(time.RFC3339, after); err == nil {
			req.UploadedAfter = &t
		}
	}

	if ownershipStatus := r.URL.Query().Get("user_owns_file"); ownershipStatus != "" {
		req.OwnershipStatus = util.ParseInt32OrDefault(ownershipStatus, 0)
	}

	if minSizeStr := r.URL.Query().Get("min_size"); minSizeStr != "" {
		if val, err := strconv.ParseInt(minSizeStr, 10, 64); err == nil {
			req.MinSize = sql.NullInt64{Int64: val, Valid: true}
		} else {
			return apierror.NewBadRequestError("Invalid min_size parameter")
		}
	}

	if maxSizeStr := r.URL.Query().Get("max_size"); maxSizeStr != "" {
		if val, err := strconv.ParseInt(maxSizeStr, 10, 64); err == nil {
			req.MaxSize = sql.NullInt64{Int64: val, Valid: true}
		} else {
			return apierror.NewBadRequestError("Invalid max_size parameter")
		}
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	// whitelisting columns to prevent users from sorting by arbitrary/unindexed columns
	allowedSortColumns := map[string]bool{
		"filename":    true,
		"size":        true,
		"uploaded_at": true,
	}
	if !allowedSortColumns[sortBy] {
		sortBy = "uploaded_at" // Default sort column
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	log.Printf("Sorting with SortBy: %s, SortOrder: %s", sortBy, sortOrder)
	req.SortBy = sortBy
	req.SortOrder = sortOrder

	contents, err := h.service.ListContents(r.Context(), req)
	if err != nil {
		log.Printf("Error while trying to retrieve contents: %s", err)
		return apierror.NewInternalServerError(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	return util.WriteJSON(w, http.StatusOK, contents)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) error {
	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		return apierror.NewBadRequestError("Missing file UUID")
	}

	err := h.service.DeleteFile(r.Context(), uuid.MustParse(fileID))
	if err != nil {
		log.Print("error deleting file:", err)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

type UpdateFilenameRequest struct {
	Filename string `json:"name"`
}

func (h *FileHandler) UpdateFilename(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var req UpdateFilenameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return apierror.NewBadRequestError("Invalid request body")
	}

	if req.Filename == "" {
		return apierror.NewBadRequestError("Filename cannot be empty")
	}

	log.Printf("Received request to rename file to: %s", req.Filename)

	fileID := uuid.MustParse(id)
	file, err := h.service.UpdateFilename(r.Context(), req.Filename, fileID)
	if err != nil {
		log.Print("Error while renaming file: ", err.Error())
		return err
	}
	return util.WriteJSON(w, http.StatusOK, file)
}

type ShareFileRequest struct {
	TargetUserID int64 `json:"target_user_id"`
}

func (h *FileHandler) ShareFile(w http.ResponseWriter, r *http.Request) error {
	fileID, _ := uuid.Parse(chi.URLParam(r, "id"))

	var req ShareFileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return apierror.NewBadRequestError("Invalid Request Body")
	}

	log.Printf("Received request to share file %s to: %d", fileID, req.TargetUserID)
	if err := h.service.ShareFile(r.Context(), fileID, req.TargetUserID); err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Shared file successfully"})
}

func (h *FileHandler) UnshareFile(w http.ResponseWriter, r *http.Request) error {
	fileID, _ := uuid.Parse(chi.URLParam(r, "id"))
	targetUserID, err := strconv.Atoi(chi.URLParam(r, "userid"))
	if err != nil {
		return apierror.NewBadRequestError("Invalid UserID")
	}

	log.Printf("Received request to unshare file %s from: %d", fileID, targetUserID)
	if err := h.service.RemoveFileShare(r.Context(), fileID, int64(targetUserID)); err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Unshared file successfully"})
}

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

type PaginatedFilesResponse struct {
	Data       []ListAllFilesResponse `json:"data"`
	TotalCount int64                  `json:"totalCount"`
}

func (h *FileHandler) ListAllFiles(w http.ResponseWriter, r *http.Request) error {
	limit := util.ParseInt32OrDefault(r.URL.Query().Get("limit"), 50)
	offset := util.ParseInt32OrDefault(r.URL.Query().Get("offset"), 0)

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	// whitelisting columns to prevent users from sorting by arbitrary/unindexed columns
	allowedSortColumns := map[string]bool{
		"filename":    true,
		"size":        true,
		"uploaded_at": true,
	}
	if !allowedSortColumns[sortBy] {
		sortBy = "uploaded_at" // Default sort column
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	log.Printf("Sorting with SortBy: %s, SortOrder: %s", sortBy, sortOrder)

	fileRows, err := h.service.ListAllFiles(r.Context(), limit, offset, sortBy, sortOrder)
	if err != nil {
		return apierror.NewInternalServerError("Failed to list all files")
	}

	files := make([]ListAllFilesResponse, 0, len(fileRows))
	for _, r := range fileRows {
		files = append(files, ListAllFilesResponse{
			ID:            r.ID.String(),
			OwnerID:       r.OwnerID,
			OwnerEmail:    r.OwnerEmail,
			Filename:      r.Filename,
			Size:          r.Size,
			DeclaredMime:  r.DeclaredMime.String,
			UploadedAt:    r.UploadedAt.Time,
			DownloadCount: &r.DownloadCount.Int64,
		})
	}

	totalCount := int64(0)
	if len(files) > 0 {
		totalCount = fileRows[0].TotalCount
	}

	response := PaginatedFilesResponse{
		Data:       files,
		TotalCount: totalCount,
	}

	return util.WriteJSON(w, http.StatusOK, response)

}

// GetShareInfo handles the HTTP request for file sharing data.
// It parses the file ID from the URL, calls the file service to retrieve the
// consolidated sharing information, and writes the result as a JSON response.
func (h *FileHandler) GetShareInfo(w http.ResponseWriter, r *http.Request) error {
	fileIDStr := chi.URLParam(r, "id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		return apierror.NewBadRequestError("Invalid fileID")
	}

	shareInfo, err := h.service.GetShareInfo(r.Context(), fileID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, shareInfo)
}

type MoveFileRequest struct {
	TargetFolderID *uuid.UUID `json:"target_folder_id"`
}

func (h *FileHandler) MoveFile(w http.ResponseWriter, r *http.Request) error {
	fileIDStr := chi.URLParam(r, "id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		return apierror.NewBadRequestError("Invalid fileID")
	}

	var req MoveFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.NewBadRequestError("Invalid request body")
	}

	log.Println(req)

	return h.service.MoveFile(r.Context(), fileID, req)
}
