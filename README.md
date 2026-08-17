# Stock Paper Trading Platform

A simple stock market paper trading platform built to practice backend development with Go.

Users can create an account, search stocks, view stock price charts, create a watchlist, and simulate buying and selling stocks using virtual money.

## Features

- Register and login
- JWT authentication
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
- JWT
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
│   ├── auth/
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

1. User registration and login
2. JWT authentication
3. Stock search
4. Stock detail and chart
5. Real-time price updates
6. Watchlist
7. Paper buy/sell orders
8. Portfolio tracking
9. Order history

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
