package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all the application configuration settings.
type Config struct {
	Server   ServerConfig
	Database DBConfig
	Minio    MinioConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string
}

// DBConfig holds database connection settings.
type DBConfig struct {
	URL string
}

// MinioConfig holds MinIO storage settings.
type MinioConfig struct {
	Endpoint string
	Access   string
	Secret   string
	Bucket   string
	Secure   bool
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() (*Config, error) {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Validate required environment variables
	if os.Getenv("PORT") == "" ||
		os.Getenv("MINIO_ENDPOINT") == "" ||
		os.Getenv("MINIO_ACCESS") == "" ||
		os.Getenv("MINIO_SECRET") == "" ||
		os.Getenv("MINIO_BUCKET") == "" ||
		os.Getenv("DB_USER") == "" ||
		os.Getenv("DB_PASSWORD") == "" ||
		os.Getenv("DB_HOST") == "" ||
		os.Getenv("DB_PORT") == "" ||
		os.Getenv("DB_NAME") == "" {
		return nil, errors.New("error: missing one or more required environment variables")
	}

	//Load Database URL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Load MinIO settings
	minioSecureStr := os.Getenv("MINIO_SECURE")
	minioSecure, err := strconv.ParseBool(minioSecureStr)
	if err != nil {
		return nil, errors.New("invalid value for MINIO_SECURE")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Database: DBConfig{
			URL: dsn,
		},
		Minio: MinioConfig{
			Endpoint: os.Getenv("MINIO_ENDPOINT"),
			Access:   os.Getenv("MINIO_ACCESS"),
			Secret:   os.Getenv("MINIO_SECRET"),
			Bucket:   os.Getenv("MINIO_BUCKET"),
			Secure:   minioSecure,
		},
	}

	return cfg, nil
}
