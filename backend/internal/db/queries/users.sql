-- name: CreateUser :one
INSERT INTO users (email, name, password, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

