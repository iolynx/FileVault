package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/admin"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apphandler"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/files"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/folders"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/middleware"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/config"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(
	cfg *config.Config,
	userHandler *users.Handler,
	fileHandler *files.FileHandler,
	folderHandler *folders.Handler,
	adminHandler *admin.Handler,
	redisClient *redis.Client,
	repo *sqlc.Queries,
) *Server {
	r := chi.NewRouter()

	// TODO: move corsoptions to env vars
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

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

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.Server.JWTSecret))

		rateLimitWindow := time.Duration(cfg.Server.RateLimitWindowSeconds) * time.Second
		r.Use(middleware.RateLimiter(redisClient, cfg.Server.RateLimit, rateLimitWindow))

		fileHandler.RegisterRoutes(r)
		folderHandler.RegisterRoutes(r)
		userHandler.RegisterRoutes(r)
	})

	// Admin Routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.Server.JWTSecret))
		r.Use(middleware.AdminMiddleware(repo))

		r.Get("/files", apphandler.MakeHTTPHandler(fileHandler.ListAllFiles))
		adminHandler.RegisterRoutes(r)
	})
	return &Server{Router: r}
}
