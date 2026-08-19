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
		query string,
		outputSize int64,
	) ([]StockSearchResponse, error)
	GetCandles(ctx context.Context,
		symbol string,
		interval Interval,
		outputSize int64,
	) ([]GetCandleResponse, error)
	GetStockDetails(
		ctx context.Context,
		symbol string,
	) (GetStockDetailsResponse, error)
}
