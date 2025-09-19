package users

import (
	"context"
	"errors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: sqlc.New(pool),
	}
}

func (r *Repository) CreateUser(ctx context.Context, email, name string, passwordHash string) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:    email,
		Name:     name,
		Password: passwordHash,
	})
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *Repository) ListOtherUsers(ctx context.Context, userID int64) ([]sqlc.ListOtherUsersRow, error) {
	return r.queries.ListOtherUsers(ctx, userID)
}

func (r *Repository) GetUserByID(ctx context.Context, userID int64) (sqlc.User, error) {
	return r.queries.GetUserByID(ctx, userID)
}
