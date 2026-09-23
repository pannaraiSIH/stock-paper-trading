package watchlist

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
)

type WatchlistRepository interface {
	CreateWatchlist(
		ctx context.Context,
		userID int64,
	) (queries.Watchlist, error)

	GetWatchlistByUserID(
		ctx context.Context,
		userID int64,
	) (queries.Watchlist, error)

	AddWatchlistItem(
		ctx context.Context,
		param queries.AddWatchlistItemParams,
	) (queries.WatchlistItem, error)

	GetWatchlistItems(
		ctx context.Context,
		param queries.GetWatchlistItemsParams,
	) ([]queries.WatchlistItem, error)

	DeleteWatchlistItem(
		ctx context.Context,
		itemID int64,
	) error
}

type watchlistRepository struct {
	store *db.Store
}

func NewWatchlistRepository(store *db.Store) *watchlistRepository {
	return &watchlistRepository{
		store: store,
	}
}

func (r *watchlistRepository) CreateWatchlist(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return r.store.Queries.CreateWatchlist(ctx, userID)
}

func (r *watchlistRepository) GetWatchlistByUserID(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return r.store.GetWatchlistByUserID(ctx, userID)
}

func (r *watchlistRepository) AddWatchlistItem(
	ctx context.Context,
	param queries.AddWatchlistItemParams,
) (queries.WatchlistItem, error) {
	return r.store.AddWatchlistItem(ctx, param)
}

func (r *watchlistRepository) GetWatchlistItems(
	ctx context.Context,
	param queries.GetWatchlistItemsParams,
) ([]queries.WatchlistItem, error) {
	return r.store.GetWatchlistItems(ctx, param)
}

func (r *watchlistRepository) DeleteWatchlistItem(
	ctx context.Context,
	itemID int64,
) error {
	return r.store.DeleteWatchlistItem(ctx, itemID)
}
