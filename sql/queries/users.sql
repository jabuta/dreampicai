-- name: CreateUser :one
INSERT INTO accounts (user_id, username, email, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM accounts WHERE user_id = $1;

-- name: UpdateUser :one
UPDATE accounts
SET username = $2, email = $3, updated_at = $4
WHERE user_id = $1
RETURNING *;
