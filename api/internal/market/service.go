package market

import (
	"context"
)

type MarketService interface {
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

type marketService struct {
	provider MarketDataProvider
}

func NewMarketService(provider MarketDataProvider) MarketService {
	return &marketService{
		provider: provider,
	}
}

func (s *marketService) SearchStocks(
	ctx context.Context,
	query SearchStocksQuery,
) ([]SearchStocksResponse, error) {
	return s.provider.SearchStocks(ctx, query)
}

func (s *marketService) GetCandles(
	ctx context.Context,
	symbol string,
	query GetCandlesQuery,
) ([]GetCandleResponse, error) {
	return s.provider.GetCandles(ctx, symbol, query)
}

func (s *marketService) GetStockDetails(
	ctx context.Context,
	symbol string,
) (GetStockDetailsResponse, error) {
	return s.provider.GetStockDetails(ctx, symbol)
}
