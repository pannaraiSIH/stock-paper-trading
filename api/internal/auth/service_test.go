package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthRepository struct {
	CreateUserFunc     func(ctx context.Context, params queries.CreateUserParams) (queries.User, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (queries.User, error)
}

func (m *MockAuthRepository) CreateUser(
	ctx context.Context,
	params queries.CreateUserParams,
) (queries.User, error) {
	return m.CreateUserFunc(ctx, params)
}

func (m *MockAuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (queries.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

func TestCreateUserService(t *testing.T) {
	tests := []struct {
		name        string
		user        queries.User
		repoErr     error
		expectedErr error
	}{
		{
			name:        "repository error",
			repoErr:     errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
		{
			name: "email already exists",
			repoErr: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "users_email_key",
			},
			expectedErr: ErrEmailAlreadyExists,
		},
		{
			name: "success",
			user: queries.User{
				ID:    1,
				Email: "dev@gmail.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			req := CreateUserRequest{
				Email:    "dev@gmail.com",
				Password: "password-123",
			}

			mockRepository := &MockAuthRepository{
				CreateUserFunc: func(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
					assert.Equal(t, req.Email, params.Email)

					err := bcrypt.CompareHashAndPassword([]byte(params.PasswordHash), []byte(req.Password))
					assert.NoError(t, err)

					return tt.user, tt.repoErr
				},
			}
			service := NewAuthService(mockRepository)

			user, err := service.CreateUser(ctx, req)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, user.ID)
			assert.Equal(t, req.Email, user.Email)
		})
	}

}

func GenPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
}

func TestLoginUser(t *testing.T) {
	dbErr := errors.New("database error")

	validHashAnother, err := GenPasswordHash("another-password")
	require.NoError(t, err)

	validHash, err := GenPasswordHash("password-123")
	require.NoError(t, err)

	tests := []struct {
		name          string
		user          queries.User
		repoErr       error
		expectedErr   error
		expectedToken bool
	}{
		{
			name:        "repository error",
			repoErr:     dbErr,
			expectedErr: dbErr,
		},
		{
			name:        "not found user",
			repoErr:     pgx.ErrNoRows,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			user: queries.User{
				Email:        "dev@gmail.com",
				PasswordHash: string(validHashAnother),
			},
			expectedErr: ErrInvalidCredentials,
		},
		{
			name: "success",
			user: queries.User{
				ID:           1,
				Email:        "dev@gmail.com",
				PasswordHash: string(validHash),
			},
			expectedToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LoginUserRequest{
				Email:    "dev@gmail.com",
				Password: "password-123",
			}

			mockRepository := &MockAuthRepository{
				GetUserByEmailFunc: func(ctx context.Context, email string) (queries.User, error) {
					assert.Equal(t, req.Email, email)

					return tt.user, tt.repoErr
				},
			}
			service := NewAuthService(mockRepository)

			accessToken, err := service.Login(t.Context(), req)

			if tt.expectedToken {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				return
			}

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Empty(t, accessToken)
		})
	}
}
