-- name: CreateUser :one
-- This query inserts a new user into the database and returns the newly created user row.
INSERT INTO users (
  email,
  password_hash
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUserByEmail :one
-- This query retrieves a single user from the database by their email address.
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
-- This query retrieves a single user from the database by their ID.
SELECT * FROM users
WHERE id = $1 LIMIT 1;