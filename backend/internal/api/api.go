package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/files"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/middleware"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(userHandler *users.Handler, fileHandler *files.FileHandler) *Server {
	r := chi.NewRouter()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	// TODO: add cors to middleware
	r.Use(corsOptions.Handler)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/signup", userHandler.Signup)
		r.Post("/auth/login", userHandler.Login)
		r.Post("/auth/logout", userHandler.Logout)
	})

	// TODO: pass the secret using config

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))

		r.Get("/auth/me", userHandler.Me)

		r.Post("/files/upload", fileHandler.Upload)

		r.Get("/files", fileHandler.ListFiles)
		//r.Get("/files/me", fileHandler.ListOwnFiles) // list of files owned by user
		r.Get("/files/url/{id}", fileHandler.GetURL)

		r.Get("/files/{id}", fileHandler.DownloadFile)
		r.Patch("/files/{id}", fileHandler.UpdateFilename)
		r.Delete("/files/{id}", fileHandler.DeleteFile)

		r.Post("/files/{id}/share", fileHandler.ShareFile)
		r.Get("/files/{id}/shares", fileHandler.GetFileShares) //get list of users with access to file
		r.Delete("/files/{id}/share/{userid}", fileHandler.UnshareFile)

		r.Get("/users", userHandler.GetOtherUsers)

	})

	return &Server{Router: r}
}
