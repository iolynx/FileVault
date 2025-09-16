package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(userHandler *users.Handler) *Server {
	r := chi.NewRouter()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	r.Use(corsOptions.Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// User routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", userHandler.Signup)
		r.Post("/login", userHandler.Login)
	})

	return &Server{Router: r}
}
