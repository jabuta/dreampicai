-- name: CreateUser :one
INSERT INTO accounts (user_id, username, email, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
