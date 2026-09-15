package trading

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AccountResponse struct {
	ID          int64          `json:"id"`
	CashBalance pgtype.Numeric `json:"cashBalance"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type CreateOrderRequest struct {
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
	Quantity int32  `json:"quantity"`
}

type GetOrdersQuery struct {
	Status string `form:"status"`
	Side   string `form:"side"`
	Offset int32  `form:"offset"`
	Limit  int32  `form:"limit"`
}

type OrderResponse struct {
	ID             int64          `json:"id"`
	Symbol         string         `json:"symbol"`
	Side           string         `json:"side"`
	ExecutionPrice pgtype.Numeric `json:"executionPrice"`
	TotalValue     pgtype.Numeric `json:"totalValue"`
	Status         string         `json:"status"`
}

type GetPositionsQuery struct {
	Symbol string `form:"symbol"`
	Offset int32  `form:"offset"`
	Limit  int32  `form:"limit"`
}

type PositionResponse struct {
	ID           int64          `json:"id"`
	Symbol       string         `json:"symbol"`
	Quantity     int            `json:"quantity"`
	AveragePrice pgtype.Numeric `json:"averagePrice"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
