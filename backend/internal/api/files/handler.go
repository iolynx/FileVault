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
	fileID := uuid.MustParse(chi.URLParam(r, "id"))

	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Missing UserID")
		return
	}

	file, err := h.service.repo.GetFileByUUID(context.Background(), fileID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "File not found")
	}

	log.Printf("received request from user %d to download file %s", ownerID, fileID)

	// Check if the user owns the file / is shared the file
	userHasAccess, err := h.service.repo.UserHasAccess(r.Context(), ownerID, fileID)
	if !userHasAccess || err != nil {
		log.Printf("no access")
		util.WriteError(w, http.StatusForbidden, "Not Allowed")
		return
	}

	blobReader, err := h.service.GetBlobReader(r.Context(), file)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Cannot read file")
		return
	}
	defer blobReader.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, blobReader)
}

type ListFilesResponse struct {
	Owned  []File `json:"owned"`
	Shared []File `json:"shared"`
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

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

	sharedFiles, err := h.service.ListFilesSharedWithUser(r.Context(), ownerID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to fetch shared files")
	}

	log.Print("Fetched files")
	resp := ListFilesResponse{
		Owned:  files,
		Shared: sharedFiles,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		http.Error(w, "missing file UUID", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteFile(context.Background(), uuid.MustParse(fileID), int64(ownerID))
	if err != nil {
		log.Print("error deleting file:", err)
		util.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete file: %s", err))
		return
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "File Deleted"})
}

type UpdateFilenameRequest struct {
	Filename string `json:"filename"`
}

func (h *FileHandler) UpdateFilename(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

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
		util.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error renaming file: %s", err))
		return
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "File Renamed"})
}

type ShareFileRequest struct {
	TargetUserID int64 `json:"target_user_id"`
}

func (h *FileHandler) ShareFile(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileID, _ := uuid.Parse(chi.URLParam(r, "id"))

	var req ShareFileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Invalid Request Body: %s", err))
		return
	}
	if req.TargetUserID <= 0 {
		util.WriteError(w, http.StatusBadRequest, "Invalid UserID")
		return
	}

	// Convert targetUserID to an int64
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if int64(req.TargetUserID) == ownerID {
		util.WriteError(w, http.StatusBadRequest, "Target User ID cannot be Your ID")
		return
	}

	log.Printf("Received request to share file %s to: %s", fileID, req.TargetUserID)
	if err := h.service.ShareFile(r.Context(), fileID, ownerID, req.TargetUserID); err != nil {
		util.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Shared file successfully"})
}

func (h *FileHandler) GetFileShares(w http.ResponseWriter, r *http.Request) {
	userID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
	}

	fileID, _ := uuid.Parse(chi.URLParam(r, "id"))

	users, err := h.service.ListUsersWithAccesToFile(r.Context(), fileID, userID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, users)
}

func (h *FileHandler) UnshareFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ownerID, ok := userctx.GetUserID(ctx)
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileID, _ := uuid.Parse(chi.URLParam(r, "id"))
	targetUserID, err := strconv.Atoi(chi.URLParam(r, "userid"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid UserID")
		return
	}

	log.Printf("Received request to unshare file %s from: %d", fileID, targetUserID)

	if err := h.service.RemoveFileShare(r.Context(), fileID, ownerID, int64(targetUserID)); err != nil {
		util.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Unshared file successfully"})
}
