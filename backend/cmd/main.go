package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/files"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/folders"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/config"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize DB
	pool := db.Connect(cfg.Database.URL)
	defer pool.Close()

	dbRepo := sqlc.New(pool)

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
	userRepo := users.NewRepository(dbRepo)
	userService := users.NewService(userRepo, os.Getenv("JWT_SECRET"), cfg)
	userHandler := users.NewHandler(userService)

	// Initialize Folders Repository, Service, Handler
	folderRepo := folders.NewRepository(dbRepo)
	folderService := folders.NewService(folderRepo, store)
	folderHandler := folders.NewHandler(folderService)

	// Initialize Files Repository, Service, Handler
	fileRepo := files.NewRepository(pool) // Initializing with pool to enable transactions
	fileService := files.NewService(fileRepo, userRepo, folderRepo, store)
	fileHandler := files.NewFileHandler(fileService)

	server := api.NewServer(cfg, userHandler, fileHandler, folderHandler, redisClient, dbRepo)

	log.Printf("Server listening on :%s", cfg.Server.Port)
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
