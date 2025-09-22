package db

import (
	"context"
	"log"

	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect establishes a PostgreSQL connection pool using the provided DSN.
// It parses the DSN, creates a connection pool, verifies connectivity with a ping,
// and returns the initialized pool. Logs fatal errors if any step fails.
func Connect(dsn string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("unable to parse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("unable to ping DB: %v", err)
	}

	log.Println("Connected to DB")
	return pool
}
