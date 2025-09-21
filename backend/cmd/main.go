package main

import (
	"context"
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
	"github.com/redis/go-redis/v9"
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

	// Initialize Redis
	redisOpts := &redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	redisClient := redis.NewClient(redisOpts)
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// Initialize Users Repository, Service, Handler
	userRepo := users.NewRepository(pool)
	userService := users.NewService(userRepo, os.Getenv("JWT_SECRET"), cfg)
	userHandler := users.NewHandler(userService)

	// Initialize Files Repository, Service, Handler
	fileRepo := files.NewRepository(pool)
	fileService := files.NewService(fileRepo, store)
	fileHandler := files.NewFileHandler(fileService)

	// Initialize Folders Repository, Service, Handler
	folderRepo := folders.NewRepository(pool)
	folderService := folders.NewService(folderRepo)
	folderHandler := folders.NewHandler(folderService)

	server := api.NewServer(userHandler, fileHandler, folderHandler, redisClient, cfg)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.Router))
}
