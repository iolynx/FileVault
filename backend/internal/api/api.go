package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/files"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/folders"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/middleware"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(userHandler *users.Handler, fileHandler *files.FileHandler, folderHandler *folders.Handler) *Server {
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

		fileHandler.RegisterRoutes(r)
		folderHandler.RegisterRoutes(r)
		userHandler.RegisterRoutes(r)
	})

	return &Server{Router: r}
}
