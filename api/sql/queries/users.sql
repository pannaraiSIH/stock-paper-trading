-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash
)
VALUES ($1, $2)
RETURNING *;

-- name: GerUserByEmail :one
SELECT * 
FROM users 
WHERE email = $1
LIMIT 1;