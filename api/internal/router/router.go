package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/handler"
)

func SetupRouter(healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", healthHandler.Health)

	return r
}
