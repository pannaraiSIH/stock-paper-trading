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

-- name: GetPositionByAccountAndSymbol :one
SELECT *
FROM positions
WHERE
    account_id = $1
    AND symbol = $2
LIMIT 1;