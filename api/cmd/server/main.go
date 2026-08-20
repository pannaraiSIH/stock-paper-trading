package main

import (
	"log"

	auth "github.com/pannaraiSIH/stock-paper-trading/internal/auth"
	"github.com/pannaraiSIH/stock-paper-trading/internal/client/twelvedata"
	"github.com/pannaraiSIH/stock-paper-trading/internal/config"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/health"
	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/pannaraiSIH/stock-paper-trading/internal/router"
)

func main() {
	cfg := config.Load()

	store, err := db.NewStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connection Postgres: %v", err)
	}
	defer store.Close()

	twelveDataClient, err := twelvedata.NewTwelveDataClient(cfg.TwelveDataApiKey)
	if err != nil {
		log.Fatalf("failed to create Twelve Data client: %v", err)
	}

	healthHandler := health.NewHealthHandler(store)

	authRepository := auth.NewAuthRepository(store)
	authService := auth.NewAuthService(authRepository, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	marketService := market.NewMarketService(twelveDataClient)
	markHandler := market.NewMarketHandler(marketService)

	r := router.SetupRouter(
		healthHandler,
		authHandler,
		markHandler,
		cfg.JWTSecret,
	)

	r.Run()
}
