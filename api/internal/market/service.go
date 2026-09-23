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
	cache    MarketCache
}

func NewMarketService(provider MarketDataProvider, cache MarketCache) MarketService {
	return &marketService{
		provider: provider,
		cache:    cache,
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
	cachedCandles, err := s.cache.GetCandles(ctx, symbol, query.Interval, query.OutputSize)
	if err != nil {
		return nil, err
	}

	if cachedCandles != nil {
		return cachedCandles, nil
	}

	candles, err := s.provider.GetCandles(ctx, symbol, query)
	if err != nil {
		return nil, err
	}

	if err = s.cache.SetCandles(
		ctx,
		symbol,
		candles,
		query.Interval,
		query.OutputSize,
	); err != nil {
		return nil, err
	}

	return candles, nil
}

func (s *marketService) GetStockDetails(
	ctx context.Context,
	symbol string,
) (GetStockDetailsResponse, error) {
	return s.provider.GetStockDetails(ctx, symbol)
}
