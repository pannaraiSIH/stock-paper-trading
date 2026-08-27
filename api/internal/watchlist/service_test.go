package watchlist

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

type MockWatchlistRepository struct {
	CreateWatchlistFunc func(
		ctx context.Context,
		userID int64,
	) (queries.Watchlist, error)

	GetWatchlistByUserIDFunc func(
		ctx context.Context,
		userID int64,
	) (queries.Watchlist, error)

	AddWatchlistItemFunc func(
		ctx context.Context,
		param queries.AddWatchlistItemParams,
	) (queries.WatchlistItem, error)

	GetWatchlistItemsFunc func(
		ctx context.Context,
		watchlistID int64,
	) ([]queries.WatchlistItem, error)

	DeleteWatchlistItemFunc func(
		ctx context.Context,
		itemID int64,
	) error
}

func (m *MockWatchlistRepository) CreateWatchlist(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return m.CreateWatchlistFunc(ctx, userID)
}

func (m *MockWatchlistRepository) GetWatchlistByUserID(
	ctx context.Context,
	userID int64,
) (queries.Watchlist, error) {
	return m.GetWatchlistByUserIDFunc(ctx, userID)
}

func (m *MockWatchlistRepository) AddWatchlistItem(
	ctx context.Context,
	param queries.AddWatchlistItemParams,
) (queries.WatchlistItem, error) {
	return m.AddWatchlistItemFunc(ctx, param)
}

func (m *MockWatchlistRepository) GetWatchlistItems(
	ctx context.Context,
	watchlistID int64,
) ([]queries.WatchlistItem, error) {
	return m.GetWatchlistItemsFunc(ctx, watchlistID)
}

func (m *MockWatchlistRepository) DeleteWatchlistItem(
	ctx context.Context,
	itemID int64,
) error {
	return m.DeleteWatchlistItemFunc(ctx, itemID)
}

func TestCreateWatchlistService(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		watchlist   queries.Watchlist
		repoErr     error
		expectedErr error
	}{
		{
			name:        "repository error",
			userID:      1,
			repoErr:     errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
		{
			name:   "success",
			userID: 1,
			watchlist: queries.Watchlist{
				ID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockWatchlistRepository{
				CreateWatchlistFunc: func(ctx context.Context, userID int64) (queries.Watchlist, error) {
					return tt.watchlist, tt.repoErr
				},
			}

			service := NewWatchlistService(repo)

			watchlist, err := service.CreateWatchlist(t.Context(), tt.userID)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, watchlist.ID)
		})
	}
}

func TestGetWatchlistService(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		watchlist   queries.Watchlist
		repoErr     error
		expectedErr error
	}{
		{
			name:        "repository error",
			userID:      1,
			repoErr:     errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
		{
			name:   "success",
			userID: 1,
			watchlist: queries.Watchlist{
				ID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockWatchlistRepository{
				GetWatchlistByUserIDFunc: func(ctx context.Context, userID int64) (queries.Watchlist, error) {
					return tt.watchlist, tt.repoErr
				},
			}

			service := NewWatchlistService(repo)

			watchlist, err := service.GetWatchlistByUserID(t.Context(), tt.userID)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, watchlist.ID)
		})
	}
}

func TestAddWatchlistItemService(t *testing.T) {
	tests := []struct {
		name        string
		symbol      string
		watchlistID int64
		item        queries.WatchlistItem
		repoErr     error
		expectedErr error
	}{
		{
			name:        "repository error",
			symbol:      "AAPL",
			watchlistID: 1,
			repoErr:     errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
		{
			name:        "watchlist item already exists ",
			symbol:      "AAPL",
			watchlistID: 1,
			repoErr: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "uq_watchlist_items_symbol",
			},
			expectedErr: ErrWatchlistItemAlreadyExists,
		},
		{
			name:        "success",
			symbol:      "AAPL",
			watchlistID: 1,
			item: queries.WatchlistItem{
				ID:     1,
				Symbol: "AAPL",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockWatchlistRepository{
				AddWatchlistItemFunc: func(ctx context.Context, param queries.AddWatchlistItemParams) (queries.WatchlistItem, error) {
					assert.Equal(t, tt.symbol, param.Symbol)

					return tt.item, tt.repoErr
				},
			}

			service := NewWatchlistService(repo)

			watchlist, err := service.AddWatchlistItem(t.Context(), tt.watchlistID, tt.symbol)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, watchlist.ID)
			assert.Equal(t, tt.symbol, watchlist.Symbol)
		})
	}
}

func TestGetWatchlistItemService(t *testing.T) {
	tests := []struct {
		name        string
		watchlistID int64
		items       []queries.WatchlistItem
		repoErr     error
		expectedErr error
	}{
		{
			name:        "repository error",
			watchlistID: 1,
			repoErr:     errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
		{
			name:        "success",
			watchlistID: 1,
			items: []queries.WatchlistItem{
				{
					ID:     1,
					Symbol: "AAPL",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockWatchlistRepository{
				GetWatchlistItemsFunc: func(ctx context.Context, watchlistID int64) ([]queries.WatchlistItem, error) {
					return tt.items, tt.repoErr
				},
			}

			service := NewWatchlistService(repo)

			items, err := service.GetWatchlistItems(t.Context(), tt.watchlistID)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}

			for _, item := range items {
				assert.NoError(t, err)
				assert.NotZero(t, item.ID)
			}

		})
	}
}

func TestDeleteWatchlistItemService(t *testing.T) {
	tests := []struct {
		name        string
		itemID      int64
		repoErr     error
		expectedErr error
	}{
		{
			name:        "repository error",
			itemID:      1,
			repoErr:     errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
		{
			name:   "success",
			itemID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockWatchlistRepository{
				DeleteWatchlistItemFunc: func(ctx context.Context, itemID int64) error {
					return tt.repoErr
				},
			}

			service := NewWatchlistService(repo)

			err := service.DeleteWatchlistItem(t.Context(), tt.itemID)

			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}

			assert.Nil(t, err)
		})
	}
}
