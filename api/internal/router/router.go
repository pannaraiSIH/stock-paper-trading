package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/auth"
	"github.com/pannaraiSIH/stock-paper-trading/internal/health"
	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
)

func SetupRouter(
	healthHandler *health.HealthHandler,
	authHandler *auth.AuthHandler,
	marketHandler *market.MarketHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", healthHandler.Health)

	api.POST("/auth/register", authHandler.CreateUser)
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(auth.AuthMiddleware(jwtSecret))

	protected.GET("/market/stocks", marketHandler.SearchStocks)
	protected.GET("/market/stocks/:symbol", marketHandler.GetStockDetails)
	protected.GET("/market/stocks/:symbol/candles", marketHandler.GetCandles)

	return r
}
