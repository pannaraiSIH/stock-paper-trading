package main

import (
	"log"

	auth "github.com/pannaraiSIH/stock-paper-trading/internal/auth"
	"github.com/pannaraiSIH/stock-paper-trading/internal/config"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/health"
	"github.com/pannaraiSIH/stock-paper-trading/internal/router"
)

func main() {
	cfg := config.Load()

	store, err := db.NewStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connection Postgres", err)
	}
	defer store.Close()

	healthHandler := health.NewHealthHandler(store)

	authRepository := auth.NewAuthRepository(store)
	authService := auth.NewAuthService(authRepository)
	authHandler := auth.NewAuthHandler(authService)

	r := router.SetupRouter(healthHandler, authHandler)

	r.Run()
}
