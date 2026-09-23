package watchlist

import "time"

type AddWatchlistResponse struct {
	ID int64 `json:"id"`
}

type AddWatchlistItemRequest struct {
	Symbol string `json:"symbol" binding:"required"`
}

type GetWatchlistItemsQuery struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit,default=20"`
}

type WatchlistItemResponse struct {
	ID        int64     `json:"id"`
	Symbol    string    `json:"symbol"`
	Price     *float64  `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}
