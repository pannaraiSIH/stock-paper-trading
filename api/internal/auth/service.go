package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository *AuthRepository
	jwtSecret  string
}

func NewAuthService(repository *AuthRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, req CreateUserRequest) (queries.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return queries.User{}, err
	}

	user, err := s.repository.CreateUser(ctx, queries.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		if db.IsUniqueViolation(err, "users_email_key") {
			return queries.User{}, ErrEmailAlreadyExists
		}
		return queries.User{}, err
	}

	return user, nil
}

func GenerateAccessToken(user queries.User, secret []byte) (string, error) {
	clams := jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clams)

	return token.SignedString(secret)
}

func (s *AuthService) Login(ctx context.Context, req LoginUserRequest) (string, error) {
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return GenerateAccessToken(user, []byte(s.jwtSecret))
}
