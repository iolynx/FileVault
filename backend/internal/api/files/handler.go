package files

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
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

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		util.WriteError(w, http.StatusBadRequest, "Could not parse form")
		return
	}

	log.Print("Parsed multipart form")

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		util.WriteError(w, http.StatusBadRequest, "No files were uploaded")
		return
	}

	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "userID missing")
	}

	for _, header := range files {
		log.Printf("Processing file: %s", header.Filename)
		file, err := header.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v", header.Filename, err)
			util.WriteError(w, http.StatusInternalServerError, "Failed to process file")
			return
		}
		// Use defer inside the loop to ensure each file is closed
		defer file.Close()

		// Call UploadFile service for each individual file
		_, err = h.service.UploadFile(r.Context(), ownerID, file, header)
		if err != nil {
			log.Printf("Upload failed for file %s: %v", header.Filename, err)
			util.WriteError(w, http.StatusInternalServerError, "Upload failed for one or more files")
			return
		}
	}

	log.Printf("Successfully uploaded %d files", len(files))
	util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Successfully uploaded %d file(s)", len(files)),
	})
}

func (h *FileHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		http.Error(w, "missing file UUID", http.StatusBadRequest)
		return
	}

	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusInternalServerError, "failed to get owner id")
		return
	}

	url, err := h.service.GetFileURL(context.Background(), uuid.MustParse(fileID), int64(ownerID))
	if err != nil {
		log.Print("failed to get file url: ", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to get file url")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"url": url,
	})
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "id")

	// ownerIDStr = r.Context().Value(middleware.UserIDKey).(string)
	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	file, err := h.service.repo.GetFileByUUID(context.Background(), uuid.MustParse(fileID))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "File not found")
	}

	// TODO: move auth checks to service layer
	if file.OwnerID != ownerID {
		util.WriteError(w, http.StatusForbidden, "Not Allowed")
	}

	blobReader, err := h.service.GetBlobReader(r.Context(), file)
	if err != nil {
		log.Print("Error while reading file: ", err)
		util.WriteError(w, http.StatusInternalServerError, "Cannot read file")
		return
	}
	defer blobReader.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, blobReader)
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ownerID, _ := userctx.GetUserID(r.Context())
	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	log.Printf("Obtaining Files for \nownerID: %d, search: %s, limitStr: %d, offsetStr: %d", ownerID, search, limit, offset)

	files, err := h.service.ListFiles(context.Background(), ownerID, search, int32(limit), int32(offset))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to fetch files")
	}
	log.Print("Fetched files")

	json.NewEncoder(w).Encode(files)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ownerID, _ := userctx.GetUserID(r.Context())

	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		http.Error(w, "missing file UUID", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteFile(context.Background(), uuid.MustParse(fileID), int64(ownerID))
	if err != nil {
		log.Print("error deleting file:", err)
		util.WriteError(w, http.StatusInternalServerError, "error deleting file")
		return
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "File Deleted"})
}

type UpdateFilenameRequest struct {
	Filename string `json:"filename"`
}

func (h *FileHandler) UpdateFilename(w http.ResponseWriter, r *http.Request) {
	ownerID, _ := userctx.GetUserID(r.Context())

	id := chi.URLParam(r, "id")
	var req UpdateFilenameRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Filename == "" {
		http.Error(w, "Filename cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to rename file to: %s", req.Filename)

	fileID := uuid.MustParse(id)
	err = h.service.UpdateFilename(context.Background(), req.Filename, fileID, ownerID)
	if err != nil {
		log.Print("Error while renaming file: ", err)
		util.WriteError(w, http.StatusInternalServerError, "Error renaming file")
		return
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "File Renamed"})
}

type ShareFileRequest struct {
	TargetUserID string `json:"target_user_id"`
}

func (h *FileHandler) ShareFile(w http.ResponseWriter, r *http.Request) {
	ownerID, _ := userctx.GetUserID(r.Context())

	fileID, _ := uuid.Parse(chi.URLParam(r, "id"))
	var req ShareFileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TargetUserID == "" {
		http.Error(w, "UserID cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to share file %s to: %s", fileID, req.TargetUserID)

	// Convert targetUserID to an int64
	targetUserID, err := strconv.Atoi(req.TargetUserID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.ShareFile(r.Context(), fileID, ownerID, int64(targetUserID)); err != nil {
		util.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Shared file"})
}
