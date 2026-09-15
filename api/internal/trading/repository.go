package trading

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/pannaraiSIH/stock-paper-trading/internal/util"
	"github.com/shopspring/decimal"
)

type TradingRepository interface {
	CreateAccount(
		ctx context.Context,
		userID int64,
	) (queries.Account, error)

	GetAccountByUserID(
		ctx context.Context,
		userID int64,
	) (queries.Account, error)

	ExecuteOrder(
		ctx context.Context,
		params queries.CreateOrderParams,
		userID int64,
	) (queries.Order, error)

	GetOrderByID(
		ctx context.Context,
		orderID int64,
	) (queries.Order, error)

	GetOrders(
		ctx context.Context,
		userID int64,
		params queries.GetOrdersParams,
	) ([]queries.Order, error)

	GetPositions(
		ctx context.Context,
		userID int64,
		params queries.GetPositionsParams,
	) ([]queries.Position, error)
}

type tradingRepository struct {
	store *db.Store
}

func NewTradingRepository(store *db.Store) *tradingRepository {
	return &tradingRepository{
		store: store,
	}
}

func (r *tradingRepository) CreateAccount(
	ctx context.Context,
	userID int64,
) (queries.Account, error) {
	return r.store.CreateAccount(ctx, userID)
}

func (r *tradingRepository) GetAccountByUserID(
	ctx context.Context,
	userID int64,
) (queries.Account, error) {
	return r.store.GetAccountByUserID(ctx, userID)
}

func (r *tradingRepository) createPendingOrder(
	ctx context.Context,
	accountID int64,
	params queries.CreateOrderParams,
) (queries.Order, error) {
	return r.store.CreateOrder(ctx, queries.CreateOrderParams{
		AccountID:      accountID,
		Symbol:         params.Symbol,
		Side:           params.Side,
		Quantity:       params.Quantity,
		ExecutionPrice: params.ExecutionPrice,
		TotalValue:     params.TotalValue,
		Status:         "pending",
	})
}

func (r *tradingRepository) executeSellOrder(
	ctx context.Context,
	qxt *queries.Queries,
	params queries.CreateOrderParams,
	accountID int64,
	cashBalance pgtype.Numeric,
) error {
	position, err := qxt.GetPositionByAccountAndSymbolForUpdate(ctx,
		queries.GetPositionByAccountAndSymbolForUpdateParams{
			AccountID: accountID,
			Symbol:    params.Symbol,
		})

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInsufficientShares
	}

	if err != nil {
		return err
	}

	if position.Quantity < params.Quantity {
		return ErrInsufficientShares
	}

	newQuantity := position.Quantity - params.Quantity

	if newQuantity == 0 {
		if _, err := qxt.DeletePosition(ctx, position.ID); err != nil {
			return err
		}
	} else {
		_, err = qxt.UpdatePosition(ctx, queries.UpdatePositionParams{
			ID:       position.ID,
			Quantity: newQuantity,
		})
		if err != nil {
			return err
		}
	}

	totalValue := util.ToDecimal(params.TotalValue)
	newBalance := util.ToDecimal(cashBalance).Add(totalValue)

	_, err = qxt.UpdateAccountBalance(ctx, queries.UpdateAccountBalanceParams{
		ID:          accountID,
		CashBalance: util.ToNumeric(newBalance),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *tradingRepository) executeBuyOrder(
	ctx context.Context,
	qxt *queries.Queries,
	params queries.CreateOrderParams,
	accountID int64,
	cashBalance pgtype.Numeric,
) error {
	cashBalanceDecimal := util.ToDecimal(cashBalance)
	totalValue := util.ToDecimal(params.TotalValue)

	if cashBalanceDecimal.Cmp(totalValue) < 0 {
		return ErrInsufficientBalance
	}

	newBalance := cashBalanceDecimal.Sub(totalValue)

	if _, err := qxt.UpdateAccountBalance(ctx, queries.UpdateAccountBalanceParams{
		ID:          accountID,
		CashBalance: util.ToNumeric(newBalance),
	}); err != nil {
		return err
	}

	position, err := qxt.GetPositionByAccountAndSymbolForUpdate(ctx,
		queries.GetPositionByAccountAndSymbolForUpdateParams{
			AccountID: accountID,
			Symbol:    params.Symbol,
		})
	if errors.Is(err, pgx.ErrNoRows) {
		_, err = qxt.CreatePosition(ctx, queries.CreatePositionParams{
			AccountID:    accountID,
			Symbol:       params.Symbol,
			Quantity:     params.Quantity,
			AveragePrice: params.ExecutionPrice,
		})

	} else if err == nil {
		averagePrice := util.ToDecimal(position.AveragePrice)
		quantity := decimal.NewFromInt(int64(position.Quantity))
		totalPrice := averagePrice.Mul(quantity)
		buyPrice := util.ToDecimal(params.ExecutionPrice)
		buyQuantity := decimal.NewFromInt(int64(params.Quantity))
		buyTotalPrice := buyPrice.Mul(buyQuantity)
		newAveragePrice := totalPrice.Add(buyTotalPrice).Div(quantity.Add(buyQuantity))

		_, err = qxt.UpdatePosition(ctx, queries.UpdatePositionParams{
			ID:           position.ID,
			Quantity:     position.Quantity + params.Quantity,
			AveragePrice: util.ToNumeric(newAveragePrice),
		})
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *tradingRepository) ExecuteOrder(
	ctx context.Context,
	params queries.CreateOrderParams,
	userID int64,
) (queries.Order, error) {
	account, err := r.store.GetAccountByUserID(ctx, userID)
	if err != nil {
		return queries.Order{}, err
	}

	pendingOrder, err := r.createPendingOrder(ctx, account.ID, params)
	if err != nil {
		return queries.Order{}, err
	}

	pool := r.store.Pool

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return queries.Order{}, err
	}
	defer tx.Rollback(ctx)

	qxt := r.store.Queries.WithTx(tx)

	lockedAccount, err := qxt.GetAccountByUserIDForUpdate(ctx, userID)
	if err != nil {
		return queries.Order{}, err
	}

	switch params.Side {
	case "buy":
		err = r.executeBuyOrder(
			ctx,
			qxt,
			params,
			lockedAccount.ID,
			lockedAccount.CashBalance,
		)

	case "sell":
		err = r.executeSellOrder(
			ctx,
			qxt,
			params,
			lockedAccount.ID,
			lockedAccount.CashBalance,
		)

	default:
		err = ErrInvalidOrderSide
	}

	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Printf("failed to rollback %v", rollbackErr)
		}

		if errors.Is(err, ErrInsufficientBalance) ||
			errors.Is(err, ErrInsufficientShares) ||
			errors.Is(err, ErrInvalidOrderSide) {
			rejectedOrder, updateErr := r.store.UpdateOrder(ctx, queries.UpdateOrderParams{
				ID:     pendingOrder.ID,
				Status: "rejected",
			})
			if updateErr != nil {
				return queries.Order{}, updateErr
			}

			return rejectedOrder, err
		}

		return queries.Order{}, err
	}

	executedOrder, err := qxt.UpdateOrder(ctx, queries.UpdateOrderParams{
		ID:     pendingOrder.ID,
		Status: "executed",
	})
	if err != nil {
		return queries.Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return queries.Order{}, err
	}

	return executedOrder, nil
}

func (r *tradingRepository) GetOrderByID(
	ctx context.Context,
	orderID int64,
) (queries.Order, error) {
	return r.store.GetOrderByID(ctx, orderID)
}

func (r *tradingRepository) GetOrders(
	ctx context.Context,
	userID int64,
	params queries.GetOrdersParams,
) ([]queries.Order, error) {
	account, err := r.store.GetAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.store.GetOrders(ctx, queries.GetOrdersParams{
		AccountID:    account.ID,
		SideFilter:   params.SideFilter,
		StatusFilter: params.StatusFilter,
		LimitCount:   params.LimitCount,
		OffsetCount:  params.OffsetCount,
	})
}

func (r *tradingRepository) GetPositions(
	ctx context.Context,
	userID int64,
	params queries.GetPositionsParams,
) ([]queries.Position, error) {
	account, err := r.store.GetAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.store.GetPositions(ctx, queries.GetPositionsParams{
		AccountID:    account.ID,
		SymbolFilter: params.SymbolFilter,
		LimitCount:   params.LimitCount,
		OffsetCount:  params.OffsetCount,
	})
}
