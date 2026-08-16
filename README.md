# Stock Paper Trading Platform

A simple stock market paper trading platform built to practice backend development with Go.

Users can search stocks, view stock price charts, create a watchlist, and simulate buying and selling stocks using virtual money.

## Features

- Search stocks
- View stock details
- View historical stock charts
- Real-time stock price updates
- Watchlist
- Paper trading
- Portfolio
- Order history

## Tech Stack

### Backend

- Go
- Gin
- PostgreSQL
- pgx
- sqlc
- Redis
- WebSocket

### Frontend

- Next.js
- TypeScript
- Tailwind CSS
- TradingView Lightweight Charts

### Infrastructure

- Docker
- Docker Compose
- GitHub Actions

## Project Structure

```text
.
├── cmd/
│   ├── api/
│   └── market-worker/
├── internal/
│   ├── user/
│   ├── market/
│   ├── candle/
│   ├── watchlist/
│   ├── order/
│   ├── portfolio/
│   └── websocket/
├── migrations/
├── sql/
├── docker-compose.yml
└── go.mod
```

## Getting Started

### 1. Start dependencies

```bash
docker compose up -d
```

### 2. Run the API

```bash
go run ./cmd/api
```

### 3. Run the market worker

```bash
go run ./cmd/market-worker
```

## MVP Scope

The MVP includes:

1. Stock search
2. Stock detail and chart
3. Real-time price updates
4. Watchlist
5. Paper buy/sell orders
6. Portfolio tracking
7. Order history

Real-money trading and brokerage integration are not included.

## Purpose

This project is mainly built to practice:

- Go API development
- Concurrency and goroutines
- WebSocket communication
- PostgreSQL transactions
- Redis caching
- Background workers
- External API integration
- Testing
