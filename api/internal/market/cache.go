package market

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MarketCache interface {
	GetLatestPrice(ctx context.Context, symbol string) (PriceEvent, error)
	SetLatestPrice(ctx context.Context, symbol string, price PriceEvent) error
	GetCandles(
		ctx context.Context,
		symbol string,
		interval Interval,
		outputSize int64,
	) ([]GetCandleResponse, error)
	SetCandles(
		ctx context.Context,
		symbol string,
		candles []GetCandleResponse,
		interval Interval,
		outputSize int64,
	) error
}

type marketCache struct {
	rdb *redis.Client
}

func NewMarketCache(rdb *redis.Client) MarketCache {
	return &marketCache{
		rdb: rdb,
	}
}

func getCachedKey(symbol string) string {
	return fmt.Sprintf("market:%s:price", symbol)
}

func getCandlesCachedKey(symbol string, interval Interval, outputSize int64) string {
	return fmt.Sprintf("market:%s:candles:%s:%d", symbol, interval, outputSize)
}

func (c *marketCache) GetLatestPrice(ctx context.Context, symbol string) (PriceEvent, error) {
	cachedKey := getCachedKey(symbol)

	cachedData, err := c.rdb.Get(ctx, cachedKey).Result()
	if err != nil {
		return PriceEvent{}, err
	}

	var price PriceEvent

	if err := json.Unmarshal([]byte(cachedData), &price); err != nil {
		return PriceEvent{}, err
	}

	return price, nil
}

func (c *marketCache) SetLatestPrice(ctx context.Context, symbol string, price PriceEvent) error {
	cachedKey := getCachedKey(symbol)

	cachedData, err := json.Marshal(price)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, cachedKey, cachedData, 60*time.Second).Err()
}

func (c *marketCache) GetCandles(
	ctx context.Context,
	symbol string,
	interval Interval,
	outputSize int64,
) ([]GetCandleResponse, error) {
	cachedKey := getCandlesCachedKey(symbol, interval, outputSize)

	cachedData, err := c.rdb.Get(ctx, cachedKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var candles []GetCandleResponse

	if err := json.Unmarshal([]byte(cachedData), &candles); err != nil {
		return nil, err
	}

	return candles, nil
}

func getCandlesTTL(interval Interval) time.Duration {
	switch interval {
	case "1min":
		return 30 * time.Second
	case "5min":
		return time.Minute
	case "1h":
		return 30 * time.Minute
	case "1day":
		return time.Hour
	default:
		return 5 * time.Minute
	}
}

func (c *marketCache) SetCandles(
	ctx context.Context,
	symbol string,
	candles []GetCandleResponse,
	interval Interval,
	outputSize int64,
) error {
	cachedKey := getCandlesCachedKey(symbol, interval, outputSize)

	cachedData, err := json.Marshal(candles)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, cachedKey, cachedData, getCandlesTTL(interval)).Err()
}
