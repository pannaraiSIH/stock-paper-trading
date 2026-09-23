package watchlist

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
)

type WatchlistHandler struct {
	service     WatchlistService
	marketCache market.MarketCache
}

func NewWatchlistHandler(service WatchlistService, marketCache market.MarketCache) *WatchlistHandler {
	return &WatchlistHandler{
		service:     service,
		marketCache: marketCache,
	}
}

func (h *WatchlistHandler) getUserID(ctx *gin.Context) (int64, bool) {
	value, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "unauthorized")
		return 0, false
	}

	userID, ok := value.(int64)
	if !ok {
		response.BadRequest(ctx, "invalid userID")
		return 0, false
	}

	return userID, true
}

func (h *WatchlistHandler) CreateWatchlist(ctx *gin.Context) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return
	}

	watchlist, err := h.service.CreateWatchlist(ctx, userID)
	if err != nil {
		response.InternalServerError(ctx, "failed to create watchlist")
		return
	}

	response.Success(ctx, http.StatusCreated, AddWatchlistResponse{
		ID: watchlist.ID,
	})
}

func (h *WatchlistHandler) getWatchlistID(ctx *gin.Context) (int64, bool) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return 0, false
	}

	watchlist, err := h.service.GetWatchlistByUserID(ctx, userID)
	if err != nil {
		response.InternalServerError(ctx, "failed to get watchlist")
		return 0, false
	}

	return watchlist.ID, true
}

func (h *WatchlistHandler) GetWatchlist(ctx *gin.Context) {
	watchlistID, ok := h.getWatchlistID(ctx)
	if !ok {
		return
	}

	response.Success(ctx, http.StatusOK, AddWatchlistResponse{
		ID: watchlistID,
	})
}

func (h *WatchlistHandler) toWatchlistItemResponse(
	ctx *gin.Context,
	item queries.WatchlistItem,
) WatchlistItemResponse {
	var price *float64

	priceEvent, err := h.marketCache.GetLatestPrice(ctx, item.Symbol)
	if err == nil {
		price = &priceEvent.Price
	}

	return WatchlistItemResponse{
		ID:        item.ID,
		Symbol:    item.Symbol,
		Price:     price,
		CreatedAt: item.CreatedAt,
	}
}

func (h *WatchlistHandler) AddWatchlistItem(ctx *gin.Context) {
	watchlistID, ok := h.getWatchlistID(ctx)
	if !ok {
		return
	}

	var req AddWatchlistItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid body request")
		return
	}

	item, err := h.service.AddWatchlistItem(ctx, watchlistID, req.Symbol)
	if err != nil {
		if errors.Is(err, ErrWatchlistItemAlreadyExists) {
			response.Conflict(ctx, ErrWatchlistItemAlreadyExists.Error())
			return
		}

		response.InternalServerError(ctx, "failed to add watchlist item")
		return
	}

	result := h.toWatchlistItemResponse(ctx, item)

	response.Success(ctx, http.StatusCreated, result)
}

func (h *WatchlistHandler) GetWatchlistItems(ctx *gin.Context) {
	watchlistID, ok := h.getWatchlistID(ctx)
	if !ok {
		return
	}

	var query GetWatchlistItemsQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(ctx, "invalid query params")
		return
	}

	watchlistItems, err := h.service.GetWatchlistItems(ctx, watchlistID, query)
	if err != nil {
		response.InternalServerError(ctx, "failed to get watchlist items")
		return
	}

	result := make([]WatchlistItemResponse, 0, len(watchlistItems))

	for _, item := range watchlistItems {
		result = append(result, h.toWatchlistItemResponse(ctx, item))
	}

	response.Success(ctx, http.StatusOK, result)
}

func (h *WatchlistHandler) DeleteWatchlistItem(ctx *gin.Context) {
	itemID, err := strconv.ParseInt(ctx.Param("itemID"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "invalid watchlist item ID")
		return
	}

	err = h.service.DeleteWatchlistItem(ctx, itemID)
	if err != nil {
		response.InternalServerError(ctx, "failed to delete watchlist item")
		return
	}

	ctx.Status(http.StatusNoContent)
}
