package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockAuthService struct {
	CreateUserFunc func(ctx context.Context, req CreateUserRequest) (queries.User, error)
	LoginFunc      func(ctx context.Context, req LoginUserRequest) (string, error)
}

func (m *MockAuthService) CreateUser(ctx context.Context, req CreateUserRequest) (queries.User, error) {
	return m.CreateUserFunc(ctx, req)
}

func (m *MockAuthService) Login(ctx context.Context, req LoginUserRequest) (string, error) {
	return m.LoginFunc(ctx, req)
}

func setupAuthTestRouter(service AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	handler := NewAuthHandler(service)

	r.POST("/register", handler.CreateUser)
	r.POST("/login", handler.Login)

	return r
}

func makeRequest(
	r *gin.Engine,
	method string,
	path string,
	body ...string,
) *httptest.ResponseRecorder {
	var reader io.Reader

	if len(body) > 0 {
		reader = strings.NewReader(body[0])
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		serviceErr   error
		expectedCode int
	}{
		{
			name: "invalid request body",
			body: `{
				"password": "pas"
			}`,
			serviceErr:   ErrInvalidRequestBody,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			body: `{
				"email": "dev",
				"password": "password-123"
			}`,
			serviceErr:   ErrInvalidRequestBody,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			body: `{
				"email": "dev@gmail.com",
				"password": "pas"
			}`,
			serviceErr:   ErrInvalidRequestBody,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "email already exists",
			body: `{
				"email": "dev@gmail.com",
				"password": "password-123"
			}`,
			serviceErr:   ErrEmailAlreadyExists,
			expectedCode: http.StatusConflict,
		},
		{
			name: "success",
			body: `{
				"email": "dev@gmail.com",
				"password": "password-123"
			}`,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				CreateUserFunc: func(ctx context.Context, req CreateUserRequest) (queries.User, error) {
					return queries.User{}, tt.serviceErr
				},
			}

			r := setupAuthTestRouter(mockService)

			w := makeRequest(r, http.MethodPost, "/register", tt.body)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}

}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		serviceErr    error
		expectedCode  int
		expectedToken bool
	}{
		{
			name: "invalid request body",
			body: `{
				"password": "password-123"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: `{
				"email": "dev@gmail.com",
				"password": "password-123"
			}`,
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "invalid credentials",
			body: `{
				"email": "dev@gmail.com",
				"password": "password-123"
			}`,
			serviceErr:   ErrInvalidCredentials,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "success",
			body: `{
				"email": "dev@gmail.com",
				"password": "password-123"
			}`,
			expectedCode:  http.StatusOK,
			expectedToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				LoginFunc: func(ctx context.Context, req LoginUserRequest) (string, error) {
					if tt.serviceErr != nil {
						return "", tt.serviceErr
					}

					return "access-token", tt.serviceErr
				},
			}

			r := setupAuthTestRouter(mockService)

			w := makeRequest(r, http.MethodPost, "/login", tt.body)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedToken {
				var res response.APIResponse[LoginUserResponse]

				err := json.Unmarshal(w.Body.Bytes(), &res)
				require.NoError(t, err)

				assert.NotEmpty(t, res.Data.AccessToken)
			}
		})
	}
}
