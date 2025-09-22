package users

import (
	"context"
	"errors"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
)

// Repository handles database operations related to users
type Repository struct {
	queries *sqlc.Queries
}

// NewRepository creates a new Repository instance with the provided database queries.
// This Repository can be used to perform User related database operations.
func NewRepository(db *sqlc.Queries) *Repository {
	return &Repository{
		queries: db,
	}
}

// CreateUser creates a new user record in the database with the provided email, name, password hash,
// and default storage quota. Returns the created user or an error if the operation fails.
func (r *Repository) CreateUser(ctx context.Context, email, name string, passwordHash string, defaultStorageQuota int64) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		Name:         name,
		Password:     passwordHash,
		StorageQuota: defaultStorageQuota,
	})
}

// GetUserByEmail retrieves a user by their email address.
// Returns a pointer to the user if found, or an error if no user exists with that email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// ListOtherUsers returns all users except the one with the specified userID.
// Used for displaying potential recipients for sharing files.
func (r *Repository) ListOtherUsers(ctx context.Context, userID int64) ([]sqlc.ListOtherUsersRow, error) {
	return r.queries.ListOtherUsers(ctx, userID)
}

// GetUserByID retrieves a user by their unique ID.
// Returns the user or an error if no user exists with the provided ID.
func (r *Repository) GetUserByID(ctx context.Context, userID int64) (sqlc.User, error) {
	return r.queries.GetUserByID(ctx, userID)
}

// GetDeduplicatedUsage returns the total storage usage for a user,
// accounting for deduplication of stored blobs. Returns the usage in bytes or an error.
func (r *Repository) GetDeduplicatedUsage(ctx context.Context, userID int64) (int64, error) {
	return r.queries.GetDeduplicatedUsage(ctx, userID)
}
