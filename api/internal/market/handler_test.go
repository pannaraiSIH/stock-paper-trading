package market

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockMarketService struct {
	SearchStocksFunc func(
		ctx context.Context,
		query SearchStocksQuery,
	) ([]SearchStocksResponse, error)

	GetCandlesFunc func(
		ctx context.Context,
		symbol string,
		query GetCandlesQuery,
	) ([]GetCandleResponse, error)

	GetStockDetailsFunc func(
		ctx context.Context,
		symbol string,
	) (GetStockDetailsResponse, error)
}

func (m *MockMarketService) SearchStocks(
	ctx context.Context,
	query SearchStocksQuery,
) ([]SearchStocksResponse, error) {
	return m.SearchStocksFunc(ctx, query)
}

func (m *MockMarketService) GetCandles(
	ctx context.Context,
	symbol string,
	query GetCandlesQuery,
) ([]GetCandleResponse, error) {
	return m.GetCandlesFunc(ctx, symbol, query)
}

func (m *MockMarketService) GetStockDetails(
	ctx context.Context,
	symbol string,
) (GetStockDetailsResponse, error) {
	return m.GetStockDetailsFunc(ctx, symbol)
}

func setupMarketTestRouter(service MarketService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	handler := NewMarketHandler(service)

	// protected := r.Group("")
	r.GET("/market/stocks", handler.SearchStocks)
	r.GET("/market/stocks/:symbol", handler.GetStockDetails)
	r.GET("/market/stocks/:symbol/candles", handler.GetCandles)

	return r
}

func makeRequest(
	r *gin.Engine,
	path string,
	method string,
	body ...string,
) *httptest.ResponseRecorder {
	var reader io.Reader

	if len(body) > 0 {
		reader = strings.NewReader(body[0])
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestSearchStocks(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		stocks       []SearchStocksResponse
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "invalid query parameters",
			query:        "?q=aapl",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			query:        "?query=aapl&outputSize=10",
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "stock not found",
			query:        "?query=aaplhh&outputSize=10",
			expectedCode: http.StatusOK,
		},
		{
			name:  "success",
			query: "?query=aapl&outputSize=10",
			stocks: []SearchStocksResponse{
				{
					Symbol:   "AAA",
					Name:     "Sunshine",
					Exchange: "ABC",
					Currency: "ARS",
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service :=
				&MockMarketService{
					SearchStocksFunc: func(ctx context.Context, query SearchStocksQuery) ([]SearchStocksResponse, error) {
						return tt.stocks, tt.serviceErr
					},
				}

			r := setupMarketTestRouter(service)

			w := makeRequest(r, "/market/stocks"+tt.query, http.MethodGet)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestGetCandles(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		candles      []GetCandleResponse
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "invalid interval",
			query:        "?interval=year&outputSize=10",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			query:        "?interval=1h&outputSize=10",
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "not found candles",
			query:        "?interval=1h&outputSize=10",
			serviceErr:   ErrStockNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "success",
			query: "?interval=1h&outputSize=10",
			candles: []GetCandleResponse{
				{
					Datetime: time.Now().String(),
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := &MockMarketService{
				GetCandlesFunc: func(ctx context.Context, symbol string, query GetCandlesQuery) ([]GetCandleResponse, error) {
					return tt.candles, tt.serviceErr
				},
			}

			r := setupMarketTestRouter(service)

			w := makeRequest(r, "/market/stocks/aapl/candles"+tt.query, http.MethodGet)

			assert.Equal(t, tt.expectedCode, w.Code)

			if len(tt.candles) > 0 {
				var res response.APIResponse[[]GetCandleResponse]

				err := json.Unmarshal(w.Body.Bytes(), &res)
				require.NoError(t, err)

				require.NotEmpty(t, res.Data)

				assert.Equal(t, tt.candles[0].Datetime, res.Data[0].Datetime)
			}
		})
	}
}

func TestGetStockDetails(t *testing.T) {
	tests := []struct {
		name         string
		details      GetStockDetailsResponse
		serviceErr   error
		expectedCode int
	}{
		{
			name:         "service error",
			serviceErr:   errors.New("service error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "stock detail not found",
			serviceErr:   ErrMarketProviderUnavailable,
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name: "success",
			details: GetStockDetailsResponse{
				Name: "test company",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockMarketService{
				GetStockDetailsFunc: func(ctx context.Context, symbol string) (GetStockDetailsResponse, error) {
					return tt.details, tt.serviceErr
				},
			}

			r := setupMarketTestRouter(service)

			w := makeRequest(r, "/market/stocks/aapl", http.MethodGet)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.details != (GetStockDetailsResponse{}) {
				var res response.APIResponse[GetStockDetailsResponse]

				err := json.Unmarshal(w.Body.Bytes(), &res)
				require.NoError(t, err)

				assert.Equal(t, tt.details.Name, res.Data.Name)
			}
		})
	}
}
