package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, ErrInvalidRequestBody.Error())
		return
	}

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			response.Conflict(c, err.Error())
		default:
			response.InternalServerError(c, "failed to create user")
		}
		return
	}

	response.Success(c, http.StatusCreated, CreateUserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginUserRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, ErrInvalidRequestBody.Error())
		return
	}

	accessToken, err := h.service.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "failed to login")
		}
		return
	}

	response.Success(c, http.StatusOK, LoginUserResponse{
		AccessToken: accessToken,
	})
}
