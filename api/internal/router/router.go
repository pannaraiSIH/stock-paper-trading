package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/auth"
	"github.com/pannaraiSIH/stock-paper-trading/internal/health"
)

func SetupRouter(healthHandler *health.HealthHandler, authHandler *auth.AuthHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", healthHandler.Health)

	api.POST("/auth/register", authHandler.CreateUser)
	api.POST("/auth/login", authHandler.Login)

	return r
}
