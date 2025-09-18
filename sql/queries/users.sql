-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUserStorageUsage :exec
-- CORRECTED: We now name the arguments using sqlc.arg() for clarity.
UPDATE users
SET storage_used_bytes = storage_used_bytes + sqlc.arg(amount)
WHERE id = sqlc.arg(id);
