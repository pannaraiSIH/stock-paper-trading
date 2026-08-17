package main

import (
	"log"

	"github.com/pannaraiSIH/stock-paper-trading/internal/config"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/handler"
	"github.com/pannaraiSIH/stock-paper-trading/internal/router"
)

func main() {
	cfg := config.Load()

	store, err := db.NewStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connection Postgres", err)
	}
	defer store.Close()

	healthHandler := handler.NewHealthHandler(store)

	r := router.SetupRouter(healthHandler)

	r.Run()
}
