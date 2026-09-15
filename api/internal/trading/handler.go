package trading

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
)

type TradingHandler struct {
	service TradingService
}

func NewTradingHandler(service TradingService) *TradingHandler {
	return &TradingHandler{
		service: service,
	}
}

func (h *TradingHandler) getUserID(ctx *gin.Context) (int64, bool) {
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

func (h *TradingHandler) CreateAccount(ctx *gin.Context) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return
	}

	account, err := h.service.CreateAccount(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrAccountAlreadyExists) {
			response.Conflict(ctx, ErrAccountAlreadyExists.Error())
			return
		}

		response.InternalServerError(ctx, "failed to create account")
		return
	}

	response.Success(ctx, http.StatusCreated, AccountResponse{
		ID:          account.ID,
		CashBalance: account.CashBalance,
		UpdatedAt:   account.UpdatedAt,
	})
}

func (h *TradingHandler) GetAccount(ctx *gin.Context) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return
	}

	account, err := h.service.GetAccount(ctx, userID)
	if err != nil {
		response.InternalServerError(ctx, "failed to get account")
		return
	}

	response.Success(ctx, http.StatusOK, AccountResponse{
		ID:          account.ID,
		CashBalance: account.CashBalance,
		UpdatedAt:   account.UpdatedAt,
	})
}

func (h *TradingHandler) CreateOrder(ctx *gin.Context) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return
	}

	var req CreateOrderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid body request")
		return
	}

	order, err := h.service.CreateOrder(ctx, req, userID)
	if err != nil {
		if errors.Is(err, ErrInsufficientBalance) {
			response.Conflict(ctx, ErrInsufficientBalance.Error())
			return
		}

		if errors.Is(err, ErrInsufficientShares) {
			response.Conflict(ctx, ErrInsufficientShares.Error())
			return
		}

		if errors.Is(err, ErrInvalidOrderSide) {
			response.BadRequest(ctx, ErrInvalidOrderSide.Error())
			return
		}

		if errors.Is(err, ErrPriceNotFound) {
			response.InternalServerError(ctx, ErrPriceNotFound.Error())
			return
		}

		response.InternalServerError(ctx, "failed to create order")
		return
	}

	response.Success(ctx, http.StatusCreated, OrderResponse{
		ID:             order.ID,
		Symbol:         order.Symbol,
		Side:           order.Side,
		ExecutionPrice: order.ExecutionPrice,
		TotalValue:     order.TotalValue,
		Status:         order.Status,
	})
}
func (h *TradingHandler) GetOrders(ctx *gin.Context) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return
	}

	var query GetOrdersQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(ctx, "invalid query parameters")
		return
	}

	orders, err := h.service.GetOrders(ctx, userID, query)
	if err != nil {
		response.InternalServerError(ctx, "failed to get orders")
		return
	}

	var data []OrderResponse

	for _, order := range orders {
		data = append(data, OrderResponse{
			ID:             order.ID,
			Symbol:         order.Symbol,
			Side:           order.Side,
			ExecutionPrice: order.ExecutionPrice,
			TotalValue:     order.TotalValue,
			Status:         order.Status,
		})
	}

	response.Success(ctx, http.StatusOK, data)
}

func (h *TradingHandler) GetOrderByID(ctx *gin.Context) {
	orderID, err := strconv.ParseInt(ctx.Param("orderID"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "invalid order id")
		return
	}

	order, err := h.service.GetOrderByID(ctx, orderID)
	if err != nil {
		response.InternalServerError(ctx, "failed to get order")
		return
	}

	response.Success(ctx, http.StatusCreated, OrderResponse{
		ID:             order.ID,
		Symbol:         order.Symbol,
		Side:           order.Side,
		ExecutionPrice: order.ExecutionPrice,
		TotalValue:     order.TotalValue,
		Status:         order.Status,
	})
}

func (h *TradingHandler) GetPositions(ctx *gin.Context) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return
	}

	var query GetPositionsQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(ctx, "invalid query parameters")
		return
	}

	positions, err := h.service.GetPositions(ctx, userID, query)
	if err != nil {
		response.InternalServerError(ctx, "failed to get positions")
		return
	}

	var data []PositionResponse

	for _, position := range positions {
		data = append(data, PositionResponse{
			ID:           position.ID,
			Symbol:       position.Symbol,
			Quantity:     int(position.Quantity),
			AveragePrice: position.AveragePrice,
			UpdatedAt:    position.UpdatedAt,
		})
	}

	response.Success(ctx, http.StatusOK, data)
}
