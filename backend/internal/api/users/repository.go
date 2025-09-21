package users

import (
	"context"
	"errors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db *sqlc.Queries) *Repository {
	return &Repository{
		queries: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, email, name string, passwordHash string, defaultStorageQuota int64) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		Name:         name,
		Password:     passwordHash,
		StorageQuota: defaultStorageQuota,
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
