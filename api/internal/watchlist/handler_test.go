package watchlist

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/stretchr/testify/assert"
)

type MockWatchlistService struct {
	CreateWatchlistFunc func(
		ctx context.Context,
		userID int64,
	) (queries.Watchlist, error)

	GetWatchlistByUserIDFunc func(
		ctx context.Context,
		userID int64,
	) (queries.Watchlist, error)

	AddWatchlistItemFunc func(
		ctx context.Context,
		watchlistID int64,
		symbol string,
	) (queries.WatchlistItem, error)

	GetWatchlistItemsFunc func(
		ctx context.Context,
		watchlistID int64,
		query GetWatchlistItemsQuery,
	) ([]queries.WatchlistItem, error)

	DeleteWatchlistItemFunc func(
		ctx context.Context,
		itemID int64,
	) error
}

func (m *MockWatchlistService) CreateWatchlist(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return m.CreateWatchlistFunc(ctx, userID)
}

func (m *MockWatchlistService) GetWatchlistByUserID(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return m.GetWatchlistByUserIDFunc(ctx, userID)
}

func (m *MockWatchlistService) AddWatchlistItem(
	ctx context.Context,
	watchlistID int64,
	symbol string,
) (queries.WatchlistItem, error) {
	return m.AddWatchlistItemFunc(ctx, watchlistID, symbol)
}

func (m *MockWatchlistService) GetWatchlistItems(
	ctx context.Context,
	watchlistID int64,
	query GetWatchlistItemsQuery,
) ([]queries.WatchlistItem, error) {
	return m.GetWatchlistItemsFunc(ctx, watchlistID, query)
}

func (m *MockWatchlistService) DeleteWatchlistItem(
	ctx context.Context,
	itemID int64,
) error {
	return m.DeleteWatchlistItemFunc(ctx, itemID)
}

type MockMarketCache struct {
	GetLatestPriceFunc func(ctx context.Context, symbol string) (market.PriceEvent, error)
	SetLatestPriceFunc func(ctx context.Context, symbol string, price market.PriceEvent) error
	GetCandlesFunc     func(
		ctx context.Context,
		symbol string,
		interval market.Interval,
		outputSize int64,
	) ([]market.GetCandleResponse, error)
	SetCandlesFunc func(
		ctx context.Context,
		symbol string,
		candles []market.GetCandleResponse,
		interval market.Interval,
		outputSize int64,
	) error
}

func (m *MockMarketCache) GetLatestPrice(
	ctx context.Context,
	symbol string,
) (market.PriceEvent, error) {
	return m.GetLatestPriceFunc(ctx, symbol)
}

func (m *MockMarketCache) SetLatestPrice(
	ctx context.Context,
	symbol string,
	price market.PriceEvent,
) error {
	return m.SetLatestPriceFunc(ctx, symbol, price)
}

func (m *MockMarketCache) GetCandles(
	ctx context.Context,
	symbol string,
	interval market.Interval,
	outputSize int64,
) ([]market.GetCandleResponse, error) {
	return m.GetCandlesFunc(ctx, symbol, interval, outputSize)
}

func (m *MockMarketCache) SetCandles(
	ctx context.Context,
	symbol string,
	candles []market.GetCandleResponse,
	interval market.Interval,
	outputSize int64,
) error {
	return m.SetCandlesFunc(ctx, symbol, candles, interval, outputSize)
}

func setupWatchlistTestRouter(
	service WatchlistService,
	marketCache market.MarketCache,
	userID any,
) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewWatchlistHandler(service, marketCache)

	r := gin.New()

	r.Use(func(ctx *gin.Context) {
		if userID != nil {
			ctx.Set("userID", userID)
		}
		ctx.Next()
	})

	r.GET("/watchlist", handler.GetWatchlist)
	r.POST("/watchlist", handler.CreateWatchlist)

	r.GET("/watchlist/items", handler.GetWatchlistItems)
	r.POST("/watchlist/items", handler.AddWatchlistItem)
	r.DELETE("/watchlist/items/:itemID", handler.DeleteWatchlistItem)

	return r
}

func makeRequest(
	r *gin.Engine,
	method string,
	path string,
	body ...string,
) *httptest.ResponseRecorder {
	var reader io.Reader

	if len(body) > 0 {
		reader = strings.NewReader(body[0])
	}

	req := httptest.NewRequest(method, path, reader)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestCreateWatchlist(t *testing.T) {
	tests := []struct {
		name         string
		userID       any
		watchlist    queries.Watchlist
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "unauthorized",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid user ID",
			userID:       "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			userID:       int64(1),
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "success",
			userID: int64(1),
			watchlist: queries.Watchlist{
				ID: 1,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockWatchlistService{
				CreateWatchlistFunc: func(ctx context.Context, userID int64) (queries.Watchlist, error) {
					return tt.watchlist, tt.serviceErr
				},
			}

			marketCache := &MockMarketCache{}

			r := setupWatchlistTestRouter(service, marketCache, tt.userID)

			w := makeRequest(r, "POST", "/watchlist")

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestGetWatchlist(t *testing.T) {
	tests := []struct {
		name         string
		userID       any
		watchlist    queries.Watchlist
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "unauthorized",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "service error",
			userID:       int64(1),
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "success",
			userID: int64(1),
			watchlist: queries.Watchlist{
				ID: 1,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockWatchlistService{
				GetWatchlistByUserIDFunc: func(ctx context.Context, userID int64) (queries.Watchlist, error) {
					return tt.watchlist, tt.serviceErr
				},
			}

			marketCache := &MockMarketCache{}

			r := setupWatchlistTestRouter(service, marketCache, tt.userID)

			w := makeRequest(r, "GET", "/watchlist")

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestAddWatchlistItem(t *testing.T) {
	tests := []struct {
		name          string
		userID        any
		watchlistID   int64
		body          string
		watchlistItem queries.WatchlistItem
		serviceErr    error
		expectedCode  int
	}{
		{
			name:         "unauthorized",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid body request",
			userID:       int64(1),
			body:         `{symbol: ""}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			userID:       int64(1),
			watchlistID:  1,
			body:         `{"symbol": "AAPL"}`,
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "watchlist item already exists",
			userID:       int64(1),
			watchlistID:  1,
			body:         `{"symbol": "AAPL"}`,
			serviceErr:   ErrWatchlistItemAlreadyExists,
			expectedCode: http.StatusConflict,
		},
		{
			name:        "success",
			userID:      int64(1),
			watchlistID: 1,
			body:        `{"symbol": "AAPL"}`,
			watchlistItem: queries.WatchlistItem{
				ID: 1,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockWatchlistService{
				GetWatchlistByUserIDFunc: func(ctx context.Context, userID int64) (queries.Watchlist, error) {
					return queries.Watchlist{
						ID: tt.watchlistID,
					}, nil
				},

				AddWatchlistItemFunc: func(ctx context.Context, watchlistID int64, symbol string) (queries.WatchlistItem, error) {
					return tt.watchlistItem, tt.serviceErr
				},
			}

			marketCache := &MockMarketCache{
				GetLatestPriceFunc: func(ctx context.Context, symbol string) (market.PriceEvent, error) {
					return market.PriceEvent{}, errors.New("cache miss")
				},
			}

			r := setupWatchlistTestRouter(service, marketCache, tt.userID)

			w := makeRequest(r, "POST", "/watchlist/items", tt.body)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestGetWatchlistItems(t *testing.T) {
	tests := []struct {
		name         string
		userID       any
		watchlistID  int64
		items        []queries.WatchlistItem
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "unauthorized",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "service error",
			userID:       int64(1),
			watchlistID:  1,
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "success",
			userID:      int64(1),
			watchlistID: 1,
			items: []queries.WatchlistItem{
				{
					ID:     1,
					Symbol: "AAPL",
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockWatchlistService{
				GetWatchlistByUserIDFunc: func(ctx context.Context, userID int64) (queries.Watchlist, error) {
					return queries.Watchlist{
						ID: tt.watchlistID,
					}, nil
				},

				GetWatchlistItemsFunc: func(ctx context.Context, watchlistID int64, query GetWatchlistItemsQuery) ([]queries.WatchlistItem, error) {
					return tt.items, tt.serviceErr
				},
			}

			marketCache := &MockMarketCache{
				GetLatestPriceFunc: func(ctx context.Context, symbol string) (market.PriceEvent, error) {
					return market.PriceEvent{}, errors.New("cache miss")
				},
			}

			r := setupWatchlistTestRouter(service, marketCache, tt.userID)

			w := makeRequest(r, "GET", "/watchlist/items")

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestDeleteWatchlistItem(t *testing.T) {
	tests := []struct {
		name         string
		userID       int64
		itemID       string
		items        []queries.WatchlistItem
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "invalid item ID",
			userID:       1,
			itemID:       "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			userID:       1,
			itemID:       "1",
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "success",
			userID: 1,
			itemID: "1",
			items: []queries.WatchlistItem{
				{
					ID:     1,
					Symbol: "AAPL",
				},
			},
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockWatchlistService{
				DeleteWatchlistItemFunc: func(ctx context.Context, itemID int64) error {
					return tt.serviceErr
				},
			}

			marketCache := &MockMarketCache{}

			r := setupWatchlistTestRouter(service, marketCache, tt.userID)

			w := makeRequest(r, "DELETE", "/watchlist/items/"+tt.itemID)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
