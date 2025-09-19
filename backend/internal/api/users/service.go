package users

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-iolynx/internal/db/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repository
	jwtSecret []byte
}

func NewService(repo *Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
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

	user, err := s.repo.CreateUser(ctx, email, name, string(passwordHash))
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
		return nil, err
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

type Me struct {
	ID                   int64  `json:"id"`
	Email                string `json:"email"`
	Name                 string `json:"name"`
	Role                 string `json:"role"`
	OriginalStorageBytes int64  `json:"original_storage_bytes"`
	DedupStorageBytes    int64  `json:"dedup_storage_bytes"`
}

func (s *Service) GetUserByID(ctx context.Context, userID int64) (Me, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return Me{}, err
	}

	return Me{
		ID:                   user.ID,
		Email:                user.Email,
		Name:                 user.Name,
		Role:                 user.Role,
		OriginalStorageBytes: user.OriginalStorageBytes,
		DedupStorageBytes:    user.DedupStorageBytes,
	}, nil

}
