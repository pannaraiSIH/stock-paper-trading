package trading

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/pannaraiSIH/stock-paper-trading/internal/util"
	"github.com/shopspring/decimal"
)

type TradingService struct {
	repository  TradingRepository
	marketCache market.MarketCache
}

func NewTradingService(repository TradingRepository, marketCache market.MarketCache) *TradingService {
	return &TradingService{
		repository:  repository,
		marketCache: marketCache,
	}
}

func (s *TradingService) CreateAccount(
	ctx context.Context,
	userID int64,
) (queries.Account, error) {
	account, err := s.repository.CreateAccount(ctx, userID)

	if err != nil {
		if db.IsUniqueViolation(err, "accounts_user_id_key") {
			return queries.Account{}, ErrAccountAlreadyExists
		}

		return queries.Account{}, err
	}

	return account, nil

}

func (s *TradingService) GetAccount(
	ctx context.Context,
	userID int64,
) (queries.Account, error) {
	return s.repository.GetAccountByUserID(ctx, userID)
}

func (s *TradingService) CreateOrder(
	ctx context.Context,
	req CreateOrderRequest,
	userID int64,
) (queries.Order, error) {
	priceEvent, err := s.marketCache.GetLatestPrice(ctx, req.Symbol)
	if err != nil {
		return queries.Order{}, ErrPriceNotFound
	}

	price := decimal.NewFromFloat(priceEvent.Price)
	totalValue := price.Mul(decimal.NewFromInt(int64(req.Quantity)))

	order, err := s.repository.ExecuteOrder(ctx, queries.CreateOrderParams{
		Symbol:         req.Symbol,
		Side:           req.Side,
		Quantity:       req.Quantity,
		ExecutionPrice: util.ToNumeric(price),
		TotalValue:     util.ToNumeric(totalValue),
	}, userID)

	if err != nil {
		if db.IsConstraintViolation(err, "chk_orders_side") {
			return queries.Order{}, ErrInvalidOrderSide
		}

		return queries.Order{}, err
	}

	return order, nil
}

func (s *TradingService) GetOrders(
	ctx context.Context,
	userID int64,
	query GetOrdersQuery,
) ([]queries.Order, error) {
	return s.repository.GetOrders(ctx, userID, queries.GetOrdersParams{
		SideFilter:   query.Side,
		StatusFilter: query.Status,
		LimitCount:   query.Limit,
		OffsetCount:  query.Offset,
	})
}

func (s *TradingService) GetOrderByID(
	ctx context.Context,
	orderID int64,
) (queries.Order, error) {
	return s.repository.GetOrderByID(ctx, orderID)
}

func (s *TradingService) GetPositions(
	ctx context.Context,
	userID int64,
	query GetPositionsQuery,
) ([]queries.Position, error) {
	return s.repository.GetPositions(ctx, userID, queries.GetPositionsParams{
		SymbolFilter: query.Symbol,
		LimitCount:   query.Limit,
		OffsetCount:  query.Offset,
	})
}
