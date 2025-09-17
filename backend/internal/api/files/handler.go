package files

import (
	"context"
	"io"
	"log"
	"net/http"

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
		log.Fatal("failed to parse form")
		util.WriteError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	log.Print("parsed")
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Fatal("Failed to read file: ", err)
		util.WriteError(w, http.StatusBadRequest, "Failed to read file:")
		return
	}
	defer file.Close()
	log.Print("read file")

	ownerID, ok := userctx.GetUserID(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "userID missing")
	}

	_, err = h.service.UploadFile(context.Background(), ownerID, file, header)
	if err != nil {
		log.Fatal("upload failed: ", err)
		util.WriteError(w, http.StatusInternalServerError, "Upload Failed")
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "file uploaded successfully",
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
		util.WriteError(w, http.StatusInternalServerError, "file not found")
	}

	if file.OwnerID != ownerID {
		util.WriteError(w, http.StatusForbidden, "not allowed")
	}

	blobReader, err := h.service.GetBlobReader(r.Context(), file)
	if err != nil {
		log.Print("Error while reading file: ", err)
		util.WriteError(w, http.StatusInternalServerError, "cannot read file")
		return
	}
	defer blobReader.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, blobReader)
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ownerID, _ := userctx.GetUserID(r.Context())
	// the middleware _should_ prune everything up to this point, but if it doesnt,
	// we have to check for an invalid ownerID but for now this is good enough

	_, err := h.service.ListFiles(context.Background(), ownerID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "error")
	}

	// add files to this later
	util.WriteJSON(w, http.StatusOK, map[string]string{
		"files": "files",
	})

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
	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "file deleted"})
}
