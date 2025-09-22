package users

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/api/apierror"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/config"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/userctx"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo                *Repository
	jwtSecret           []byte
	defaultStorageQuota int64
}

func NewService(repo *Repository, jwtSecret string, cfg *config.Config) *Service {
	return &Service{
		repo:                repo,
		jwtSecret:           []byte(jwtSecret),
		defaultStorageQuota: cfg.Server.DefaultStorageQuota,
	}
}

func (s *Service) Signup(ctx context.Context, email, name string, password string) (sqlc.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("Failed to hash password: %w", err)
	}

	_, err = s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return sqlc.User{}, fmt.Errorf("User already exists")
	}

	user, err := s.repo.CreateUser(ctx, email, name, string(passwordHash), s.defaultStorageQuota)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("Failed to create user: %w", err)
	}
	return user, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (int64, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("Log In Failed: No Such User")
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Log In Failed: Invalid Credentials")
		return 0, err
	}

	log.Println("User Authenticated")
	return user.ID, nil
}

func (s *Service) GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *Service) ListOtherUsers(ctx context.Context, userID int64) ([]User, error) {
	otherUsersRow, err := s.repo.ListOtherUsers(ctx, userID)
	if err != nil {
		return nil, apierror.NewInternalServerError("error while fetching users")
	}

	otherUsers := make([]User, 0, len(otherUsersRow))
	for _, r := range otherUsersRow {
		otherUsers = append(otherUsers, User{
			ID:    r.ID,
			Name:  r.Name,
			Email: r.Email,
		})
	}
	return otherUsers, nil
}

type MeResponse struct {
	ID                     int64   `json:"id"`
	Email                  string  `json:"email"`
	Name                   string  `json:"name"`
	Role                   string  `json:"role"`
	StorageUsedBytes       int64   `json:"storage_used_bytes"`       // "Original storage usage"
	DeduplicatedUsageBytes int64   `json:"deduplicated_usage_bytes"` // "Total storage used (deduplicated)"
	StorageQuotaBytes      int64   `json:"storage_quota_bytes"`
	SavingsBytes           int64   `json:"savings_bytes"`
	SavingsPercentage      float64 `json:"savings_percentage"`
}

// GetMe is the complete service for the /auth/me endpoint
// it returns the User info along with role, and storage statistics.
func (s *Service) GetMe(ctx context.Context) (MeResponse, error) {
	userID, ok := userctx.GetUserID(ctx)
	if !ok {
		return MeResponse{}, apierror.NewUnauthorizedError()
	}

	// Fetch the User record from the database.
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return MeResponse{}, apierror.NewInternalServerError("could not retrieve user data")
	}

	// The 'storage_used' column becomes the "Original Storage Usage"
	originalUsage := user.StorageUsed

	// we calculate "Deduplicated Storage Usage" with our new query
	deduplicatedUsage, err := s.repo.GetDeduplicatedUsage(ctx, userID)
	if err != nil {
		return MeResponse{}, apierror.NewInternalServerError("could not calculate storage stats")
	}

	savingsBytes := originalUsage - deduplicatedUsage
	savingsPercentage := 0.0
	if originalUsage > 0 {
		savingsPercentage = (float64(savingsBytes) / float64(originalUsage)) * 100
	}

	return MeResponse{
		ID:                     user.ID,
		Email:                  user.Email,
		Name:                   user.Name,
		Role:                   user.Role,
		StorageUsedBytes:       originalUsage,
		DeduplicatedUsageBytes: deduplicatedUsage,
		StorageQuotaBytes:      user.StorageQuota,
		SavingsBytes:           savingsBytes,
		SavingsPercentage:      savingsPercentage,
	}, nil
}
