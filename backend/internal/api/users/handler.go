package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apphandler"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/go-chi/chi/v5"
)

// Handler provides HTTP route handlers for user-related endpoints.
// It delegates business logic to the underlying Service.
type Handler struct {
	service *Service
}

// NewHandler creates a new Handler instance with the provided Service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the user-related routes (auth, users) on the router.
// Currently includes: /auth/me and /users.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/auth/me", apphandler.MakeHTTPHandler(h.Me))
	r.Get("/users", apphandler.MakeHTTPHandler(h.GetOtherUsers))
}

// Signup handles user registration requests.
// It decodes the request body, calls the service layer to create the user,
// and responds with HTTP 201 on success.
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

// Login handles user login requests.
// It authenticates the user, generates a JWT on success, and sets it as an HTTP-only cookie.
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

// Logout handles user logout requests.
// It clears the JWT cookie, effectively logging the user out.
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

// GetOtherUsers handles the /users endpoint.
// It returns a list of all users except the currently authenticated user.
func (h *Handler) GetOtherUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return apierror.NewUnauthorizedError()
	}

	otherUsers, err := h.service.ListOtherUsers(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, otherUsers)
}

// Me handles the /auth/me endpoint.
// It returns information about the currently authenticated user, including role and storage stats.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) error {
	user, err := h.service.GetMe(r.Context())
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, user)
}
