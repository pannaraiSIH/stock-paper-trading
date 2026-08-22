package market

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MarketCache interface {
	GetLatestPrice(ctx context.Context, symbol string) (PriceEvent, error)
	SetLatestPrice(ctx context.Context, symbol string, price PriceEvent) error
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
