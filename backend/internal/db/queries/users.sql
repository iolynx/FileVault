-- name: CreateUser :one
INSERT INTO users (email, name, password, created_at, storage_quota)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ListOtherUsers :many
SELECT id, email, name
FROM users
WHERE id <> $1
ORDER BY name;

-- name: GetDeduplicatedUsage :one
SELECT COALESCE(SUM(b.size), 0)::BIGINT
FROM blobs b
WHERE b.id IN (
    SELECT DISTINCT blob_id FROM files WHERE owner_id = $1
);
