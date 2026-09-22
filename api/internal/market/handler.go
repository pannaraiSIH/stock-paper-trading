package market

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
)

type MarketHandler struct {
	service MarketService
}

func NewMarketHandler(service MarketService) *MarketHandler {
	return &MarketHandler{
		service: service,
	}
}

func (h *MarketHandler) SearchStocks(c *gin.Context) {
	var query SearchStocksQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, ErrInvalidQueryParameters.Error())
		return
	}

	stocks, err := h.service.SearchStocks(c, query)
	if err != nil {
		handleMarketError(c, err, "failed to search stocks")
		return
	}

	response.Success(c, http.StatusOK, stocks)
}

func (h *MarketHandler) GetCandles(c *gin.Context) {
	var query GetCandlesQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, ErrInvalidQueryParameters.Error())
		return
	}

	candles, err := h.service.GetCandles(c, query.Symbol, query)
	if err != nil {
		handleMarketError(c, err, "failed to get candles")
		return
	}

	response.Success(c, http.StatusOK, candles)
}

func (h *MarketHandler) GetStockDetails(c *gin.Context) {
	var query GetStockDetailsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, ErrInvalidQueryParameters.Error())
		return
	}

	detail, err := h.service.GetStockDetails(c, query.Symbol)
	if err != nil {
		handleMarketError(c, err, "failed to get stock details")
		return
	}

	response.Success(c, http.StatusOK, detail)
}

func handleMarketError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, ErrStockNotFound):
		response.NotFound(c, ErrStockNotFound.Error())
	case errors.Is(err, ErrMarketProviderUnavailable):
		response.ServiceUnavailable(c, ErrMarketProviderUnavailable.Error())
	default:
		response.InternalServerError(c, fallbackMessage)
	}
}
