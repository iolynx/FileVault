package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/users"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool := db.Connect(dsn)
	defer pool.Close()

	userRepo := users.NewRepository(pool)
	userService := users.NewService(userRepo, os.Getenv("JWT_SECRET"))
	userHandler := users.NewHandler(userService)

	server := api.NewServer(userHandler)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.Router))
}
