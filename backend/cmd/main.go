package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/files"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/folders"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/config"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize DB
	pool := db.Connect(cfg.Database.URL)
	defer pool.Close()

	// Initialize Minio, Storage Handler
	store, err := storage.NewMinioStorage(cfg.Minio)
	if err != nil {
		log.Fatal("Failed to initialize storage", err)
	}

	// Initialize Users Repository, Service, Handler
	userRepo := users.NewRepository(pool)
	userService := users.NewService(userRepo, os.Getenv("JWT_SECRET"))
	userHandler := users.NewHandler(userService)

	// Initialize Files Repository, Service, Handler
	fileRepo := files.NewRepository(pool)
	fileService := files.NewService(fileRepo, store)
	fileHandler := files.NewFileHandler(fileService)

	// Initialize Folders Repository, Service, Handler
	folderRepo := folders.NewRepository(pool)
	foldersService := folders.NewService(folderRepo)
	foldersHandler := folders.NewHandler(foldersService)

	server := api.NewServer(userHandler, fileHandler, foldersHandler)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.Router))
}
