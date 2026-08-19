package market

import (
	"context"
)

type MarketService struct {
	provider MarketDataProvider
}

func NewMarketService(provider MarketDataProvider) *MarketService {
	return &MarketService{
		provider: provider,
	}
}

func (s *MarketService) SearchStocks(
	ctx context.Context,
	query SearchStocksQuery,
) ([]StockSearchResponse, error) {
	return s.provider.SearchStocks(ctx, query.Query, query.OutputSize)
}

func (s *MarketService) GetCandles(
	ctx context.Context,
	symbol string,
	query GetCandlesQuery,
) ([]GetCandleResponse, error) {
	return s.provider.GetCandles(ctx, symbol, query.Interval, query.OutputSize)
}
