-- name: CreateUser :one
INSERT INTO users (email, name, password, created_at, storage_quota)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE id = $1
);

-- name: ListOtherUsers :many
SELECT id, email, name
FROM users
WHERE id <> $1
ORDER BY name;
