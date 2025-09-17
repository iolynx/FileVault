package users

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
		"user_id": strconv.FormatInt(userID, 10),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
