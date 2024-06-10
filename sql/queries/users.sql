-- name: CreateUser :one
INSERT INTO accounts (user_id, username, email, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM accounts WHERE user_id = $1;