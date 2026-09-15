-- name: CreateAccount :one
INSERT INTO accounts (
    user_id
)
VALUES ($1)
RETURNING *;

-- name: GetAccountByUserID :one
SELECT *
FROM accounts 
WHERE user_id = $1
LIMIT 1;

-- name: GetAccountByUserIDForUpdate :one
SELECT *
FROM accounts 
WHERE user_id = $1
LIMIT 1
FOR UPDATE;

-- name: UpdateAccountBalance :one 
UPDATE accounts 
SET
    cash_balance = $2
WHERE id = $1
RETURNING *;


-- name: CreateOrder :one
INSERT INTO orders (
    account_id,
    symbol,
    side,
    quantity,
    execution_price,
    total_value,
    status
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders 
SET 
    execution_price = $2, 
    total_value = $3, 
    status = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1
LIMIT 1;

-- name: GetOrdersByAccountID :many
SELECT * 
FROM orders
WHERE account_id = $1
ORDER BY created_at DESC;

-- name: GetOrders :many
SELECT *
FROM orders
WHERE account_id = sqlc.arg(account_id)
  AND (
    COALESCE(sqlc.arg(status_filter)::text, '') = ''
    OR status = sqlc.arg(status_filter)::text
  )
  AND (
    COALESCE(sqlc.arg(side_filter)::text, '') = ''
    OR side = sqlc.arg(side_filter)::text
  )
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count)::int
OFFSET sqlc.arg(offset_count)::int;


-- name: CreatePosition :one
INSERT INTO positions (
    account_id,
    symbol,
    quantity,
    average_price
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePosition :one
UPDATE positions 
SET 
    quantity = $2, 
    average_price = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetPositionsByAccountID :many
SELECT * 
FROM positions
WHERE account_id = $1
ORDER BY symbol;

-- name: GetPositions :many 
SELECT *
FROM positions
WHERE account_id = sqlc.arg(account_id)
  AND (
    COALESCE(sqlc.arg(symbol_filter)::text, '') = ''
    OR symbol = sqlc.arg(symbol_filter)::text
  )
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count)::int
OFFSET sqlc.arg(offset_count)::int;

-- name: GetPositionByAccountAndSymbolForUpdate :one
SELECT *
FROM positions
WHERE
    account_id = $1
    AND symbol = $2
LIMIT 1
FOR UPDATE;

-- name: DeletePosition :one
DELETE FROM positions
WHERE id = $1
RETURNING *;
