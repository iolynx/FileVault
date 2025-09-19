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
SELECT id, filename, size, declared_mime as content_type, uploaded_at
FROM files
WHERE owner_id = $1 AND ($2 = '' OR filename ILIKE '%' || $2 || '%')
ORDER BY uploaded_at DESC
LIMIT $3 OFFSET $4;

-- name: GetFileByUUID :one
SELECT *
FROM files f
WHERE f.id = $1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;

-- name: GetFilesForUser :many
SELECT id, filename, size, declared_mime as mime_type, uploaded_at, is_public
FROM files
WHERE owner_id = $1
AND (sqlc.narg('search')::TEXT IS NULL OR filename ILIKE '%' || sqlc.narg('search')::TEXT || '%')
ORDER BY uploaded_at DESC
LIMIT $2
OFFSET $3;

-- name: GetFilesForUserCount :one
SELECT count(*)
FROM files
WHERE owner_id = $1
AND (sqlc.narg('search')::TEXT IS NULL OR filename ILIKE '%' || sqlc.narg('search')::TEXT || '%');


-- name: UpdateFilename :exec
UPDATE files
SET filename = $1
WHERE id = $2;

-- name: ListFilesSharedWithUser :many
-- name: ListFilesSharedWithUser :many
SELECT f.*
FROM files f
JOIN file_shares fs ON f.id = fs.file_id
WHERE fs.shared_with = $1
  AND ($2 = '' OR f.filename ILIKE '%' || $2 || '%')
ORDER BY f.uploaded_at DESC
LIMIT $3 OFFSET $4;

-- name: UserHasAccess :one
SELECT EXISTS (
  SELECT 1
  FROM files f
  LEFT JOIN file_shares fs
    ON f.id = fs.file_id AND fs.shared_with = $1
  WHERE f.id = $2
    AND (f.owner_id = $1 OR fs.shared_with = $1)
);

-- name: ListUsersWithAccessToFile :many
SELECT u.id, u.name, u.email, fs.permission
FROM file_shares fs
JOIN users u ON u.id = fs.shared_with
WHERE fs.file_id = $1;

-- name: CreateFileShare :one
INSERT INTO file_shares (file_id, shared_with)
VALUES ($1, $2)
RETURNING id, file_id, shared_with, permission, created_at;

-- name: DeleteFileShare :exec
DELETE FROM file_shares
WHERE file_id = $1 AND shared_with = $2;


-- name: ListFilesForUser :many
SELECT DISTINCT f.id,
       f.owner_id,
       f.filename,
       f.size,
       f.declared_mime AS content_type,
       f.uploaded_at,
       (f.owner_id = $1) AS user_owns_file
FROM files f
LEFT JOIN file_shares fs ON f.id = fs.file_id
WHERE (f.owner_id = $1 OR fs.shared_with = $1)
  AND ($2 = '' OR f.filename ILIKE '%' || $2 || '%')
ORDER BY f.uploaded_at DESC
LIMIT $3 OFFSET $4;


-- name: IncrementUserStorage :exec
UPDATE users
SET original_storage_bytes = original_storage_bytes + $2,
    dedup_storage_bytes = dedup_storage_bytes + $3
WHERE id = $1;

-- name: DecrementUserStorage :exec
UPDATE users
SET original_storage_bytes = GREATEST(original_storage_bytes - $2, 0),
    dedup_storage_bytes = GREATEST(dedup_storage_bytes - $3, 0)
WHERE id = $1;
