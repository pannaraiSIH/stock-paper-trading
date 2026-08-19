package market

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
)

type MarketHandler struct {
	service *MarketService
}

func NewMarketHandler(service *MarketService) *MarketHandler {
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
		response.InternalServerError(c, "failed to search stocks")
		return
	}

	response.Success(c, http.StatusOK, stocks)
}

func (h *MarketHandler) GetStockDetails(c *gin.Context) {
	symbol := c.Param("symbol")

	detail, err := h.service.GetStockDetails(c, symbol)
	if err != nil {
		response.InternalServerError(c, "failed to get stock details")
		return
	}

	response.Success(c, http.StatusOK, detail)
}

func (h *MarketHandler) GetCandles(c *gin.Context) {
	symbol := c.Param("symbol")

	var query GetCandlesQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Println(err)
		response.BadRequest(c, ErrInvalidQueryParameters.Error())
		return
	}

	candles, err := h.service.GetCandles(c, symbol, query)
	if err != nil {
		response.InternalServerError(c, "failed to get candles")
		return
	}

	response.Success(c, http.StatusOK, candles)
}
