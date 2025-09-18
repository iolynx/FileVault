-- name: CreateUser :one
INSERT INTO users (email, name, password, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE id = $1
);

-- name: ListOtherUsers :many
SELECT id, email, name
FROM users
WHERE id <> $1
ORDER BY name;
