package users

import (
	"context"
	"encoding/json"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"log"
	"net/http"
	"time"
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
		util.WriteError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userID, err := h.service.AuthenticateUser(context.Background(), req.Email, req.Password)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
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
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Login Succesfull",
	})
	w.Write([]byte("Logged in successfully."))
}
