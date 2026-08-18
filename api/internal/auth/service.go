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

type CustomClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
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
	clams := &CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "stock-paper-trading",
		},
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
