package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func Success[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, APIResponse[T]{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func ServiceUnavailable(c *gin.Context, message string) {
	Error(c, http.StatusServiceUnavailable, message)
}
