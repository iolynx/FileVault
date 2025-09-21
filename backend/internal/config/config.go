package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/util"
	"github.com/joho/godotenv"
)

// Config holds all the application configuration settings.
type Config struct {
	Server   ServerConfig
	Database DBConfig
	Minio    MinioConfig
	Redis    RedisConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port                   string
	DefaultStorageQuota    int64
	RateLimit              int
	RateLimitWindowSeconds int
	JWTSecret              string
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

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	requiredVars := []string{
		"PORT", "MINIO_ENDPOINT", "MINIO_ACCESS", "MINIO_SECRET",
		"MINIO_BUCKET", "DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT",
		"DB_NAME", "REDIS_ADDR", "DEFAULT_STORAGE_QUOTA", "API_RATE_LIMIT",
		"API_RATE_LIMIT_WINDOW_SECONDS", "JWT_SECRET",
	}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf("error: missing required environment variable: %s", v)
		}
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

	defaultQuota, err := strconv.ParseInt(os.Getenv("DEFAULT_STORAGE_QUOTA"), 10, 64)
	if err != nil {
		return nil, errors.New("invalid value for DEFAULT_STORAGE_QUOTA")
	}
	RateLimit, err := strconv.Atoi(os.Getenv("API_RATE_LIMIT"))
	if err != nil {
		return nil, errors.New("invalid value for API_RATE_LIMIT")
	}
	RateLimitWindowSeconds, err := strconv.Atoi(os.Getenv("API_RATE_LIMIT_WINDOW_SECONDS"))
	if err != nil {
		return nil, errors.New("invalid value for API_RATE_LIMIT_WINDOW_SECONDS")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:                   os.Getenv("PORT"),
			DefaultStorageQuota:    defaultQuota,
			RateLimit:              RateLimit,
			RateLimitWindowSeconds: RateLimitWindowSeconds,
			JWTSecret:              os.Getenv("JWT_SECRET"),
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
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       util.ParseIntOrDefault(os.Getenv("REDIS_DB"), 0),
		},
	}

	return cfg, nil
}
