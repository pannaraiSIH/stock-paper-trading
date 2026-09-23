package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	auth "github.com/pannaraiSIH/stock-paper-trading/internal/auth"
	"github.com/pannaraiSIH/stock-paper-trading/internal/client/twelvedata"
	"github.com/pannaraiSIH/stock-paper-trading/internal/config"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/health"
	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/pannaraiSIH/stock-paper-trading/internal/redis"
	"github.com/pannaraiSIH/stock-paper-trading/internal/router"
	"github.com/pannaraiSIH/stock-paper-trading/internal/trading"
	"github.com/pannaraiSIH/stock-paper-trading/internal/watchlist"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, err := db.NewStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect Postgres: %v", err)
	}
	defer store.Close()

	rdb, err := redis.NewRedisClient(context.Background(), cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("failed to connect Redis: %v", err)
	}
	defer rdb.Close()

	marketDataClient, err := twelvedata.NewTwelveDataClient(cfg.TwelveDataApiKey)
	if err != nil {
		log.Fatalf("failed to create Twelve Data client: %v", err)
	}

	realtimeClient, err := twelvedata.NewTwelveDataWebsocket(cfg.TwelveDataApiKey)
	if err != nil {
		log.Fatalf("failed to create Twelve Data websocket: %v", err)
	}
	defer realtimeClient.Disconnect()

	hub := market.NewHub()
	marketCache := market.NewMarketCache(rdb)

	worker := market.NewMarketWorker(realtimeClient, hub, marketCache)

	go worker.Run(ctx)

	healthHandler := health.NewHealthHandler(store)

	authRepository := auth.NewAuthRepository(store)
	authService := auth.NewAuthService(authRepository, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	marketService := market.NewMarketService(marketDataClient, marketCache)
	markHandler := market.NewMarketHandler(marketService)

	watchlistRepository := watchlist.NewWatchlistRepository(store)
	watchlistService := watchlist.NewWatchlistService(watchlistRepository)
	watchlistHandler := watchlist.NewWatchlistHandler(watchlistService, marketCache)

	tradingRepository := trading.NewTradingRepository(store)
	tradingService := trading.NewTradingService(tradingRepository, marketCache)
	tradingHandler := trading.NewTradingHandler(*tradingService)

	r := router.SetupRouter(
		healthHandler,
		authHandler,
		markHandler,
		hub,
		watchlistHandler,
		tradingHandler,
		cfg.JWTSecret,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("starting server on :%s", cfg.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
