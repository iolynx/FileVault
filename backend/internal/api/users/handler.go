package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type signupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// not returning User for now
	_, err := h.service.Signup(context.Background(), req.Email, req.Name, req.Password)
	if err != nil {
		log.Printf("Sign Up Error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Sign Up Error: %s", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	log.Printf("Received request to Login user")
	userID, err := h.service.AuthenticateUser(context.Background(), req.Email, req.Password)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		util.WriteError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	token, err := h.service.GenerateToken(userID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Domain:   "",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	util.WriteJSON(w, http.StatusOK, "Login successful")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to Logout user")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Domain:   "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

}

func (h *Handler) GetOtherUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Missing User ID")
		return
	}

	otherUsers, err := h.service.ListOtherUsers(ctx, userID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Error: %s", err))
		return
	}

	util.WriteJSON(w, http.StatusOK, otherUsers)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Missing User ID")
		return
	}

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Error while fetching user info")
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}
