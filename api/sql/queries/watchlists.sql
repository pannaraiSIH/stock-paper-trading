-- name: CreateWatchlist :one
INSERT INTO watchlists (
    user_id
)
VALUES ($1)
RETURNING *;

-- name: GetWatchlistByUserID :one
SELECT * 
FROM watchlists 
WHERE user_id = $1
LIMIT 1;

-- name: AddWatchlistItem :one
INSERT INTO watchlist_items (
    watchlist_id,
    symbol
)
VALUES ($1, $2)
RETURNING *;

-- name: GetWatchlistItems :many
SELECT *
FROM watchlist_items 
WHERE watchlist_id = $1
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count)::int 
OFFSET sqlc.arg(offset_count)::int;

-- name: DeleteWatchlistItem :exec
DELETE FROM watchlist_items
WHERE id = $1;
