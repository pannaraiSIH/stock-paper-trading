package watchlist

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
)

type WatchlistService interface {
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
		watchlistID int64,
		symbol string,
	) (queries.WatchlistItem, error)

	GetWatchlistItems(
		ctx context.Context,
		watchlistID int64,
	) ([]queries.WatchlistItem, error)

	DeleteWatchlistItem(
		ctx context.Context,
		itemID int64,
	) error
}

type watchlistService struct {
	repository WatchlistRepository
}

func NewWatchlistService(repository WatchlistRepository) WatchlistService {
	return &watchlistService{
		repository: repository,
	}
}

func (s *watchlistService) CreateWatchlist(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	watchlist, err := s.repository.CreateWatchlist(ctx, userID)
	if err != nil {
		if db.IsUniqueViolation(err, "watchlists_user_id_key") {
			existing, err := s.repository.GetWatchlistByUserID(ctx, userID)

			if err != nil {
				return queries.Watchlist{}, err
			}

			return existing, nil
		}

		return queries.Watchlist{}, err
	}

	return watchlist, nil
}

func (s *watchlistService) GetWatchlistByUserID(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return s.repository.GetWatchlistByUserID(ctx, userID)
}

func (s *watchlistService) AddWatchlistItem(
	ctx context.Context,
	watchlistID int64,
	symbol string,
) (queries.WatchlistItem, error) {
	watchlistItem, err := s.repository.AddWatchlistItem(ctx, queries.AddWatchlistItemParams{
		WatchlistID: watchlistID,
		Symbol:      symbol,
	})
	if err != nil {
		if db.IsUniqueViolation(err, "uq_watchlist_items_symbol") {
			return queries.WatchlistItem{}, ErrWatchlistItemAlreadyExists
		}

		return queries.WatchlistItem{}, err
	}

	return watchlistItem, nil
}

func (s *watchlistService) GetWatchlistItems(
	ctx context.Context,
	watchlistID int64,
) ([]queries.WatchlistItem, error) {
	return s.repository.GetWatchlistItems(ctx, watchlistID)
}

func (s *watchlistService) DeleteWatchlistItem(
	ctx context.Context,
	itemID int64,
) error {
	return s.repository.DeleteWatchlistItem(ctx, itemID)
}
