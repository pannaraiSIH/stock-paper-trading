package market

import (
	"context"
)

type Interval string

const (
	Interval1Min Interval = "1min"
	Interval5Min Interval = "5min"
	Interval1H   Interval = "1h"
	Interval1Day Interval = "1day"
)

type MarketDataProvider interface {
	SearchStocks(
		ctx context.Context,
		query SearchStocksQuery,
	) ([]SearchStocksResponse, error)

	GetCandles(
		ctx context.Context,
		symbol string,
		query GetCandlesQuery,
	) ([]GetCandleResponse, error)

	GetStockDetails(
		ctx context.Context,
		symbol string,
	) (GetStockDetailsResponse, error)
}
