-- name: CreateBlob :one
INSERT INTO blobs (sha256, storage_path, size, mime_type, refcount)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteBlob :exec
DELETE FROM blobs
WHERE id = $1;

-- name: GetBlobBySha :one
SELECT * FROM blobs WHERE sha256 = $1;

-- name: GetBlobByID :one
SELECT * FROM blobs WHERE id = $1;

-- name: UpdateBlobRefcount :exec
UPDATE blobs SET refcount = refcount + $2 WHERE id = $1;

-- name: UserOwnsBlob :one
SELECT 1 FROM files WHERE owner_id = $1 AND blob_id = $2 LIMIT 1;

-- name: IncrementBlobRefcount :one
UPDATE blobs SET refcount = refcount + 1
WHERE id = $1
RETURNING refcount;

-- name: DecrementBlobRefcount :one
UPDATE blobs SET refcount = refcount - 1
WHERE id = $1
RETURNING refcount;

-- name: DeleteBlobIfUnused :exec
DELETE FROM blobs WHERE id = $1 AND refcount <= 0;



-- name: CreateFile :one
INSERT INTO files (owner_id, blob_id, filename, declared_mime, size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListFilesByOwner :many
SELECT *
FROM files
WHERE owner_id = $1
ORDER BY uploaded_at DESC;

-- name: GetFileByUUID :one
SELECT *
FROM files f
WHERE f.id = $1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;
