-- name: CreateBlob :one
INSERT INTO blobs (sha256, storage_path, size, mime_type, refcount)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteBlob :exec
DELETE FROM blobs
WHERE id = $1;

-- name: GetBlobBySha :one
SELECT * FROM blobs WHERE sha256 = $1 LIMIT 1;

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

-- name: DeleteBlobIfUnused :one
DELETE FROM blobs 
WHERE id = $1 AND refcount <= 0
RETURNING storage_path;

-- name: DeleteBlobsByStoragePaths :exec
DELETE FROM blobs
WHERE storage_path = ANY(sqlc.arg(storage_paths)::text[]) AND refcount <= 1;



-- name: CreateFile :one
INSERT INTO files (owner_id, blob_id, filename, declared_mime, size, folder_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListFilesByOwner :many
SELECT id, filename, size, declared_mime as content_type, uploaded_at, is_public, download_count
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
SELECT id, filename, size, declared_mime as mime_type, uploaded_at, is_public, download_count
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


-- name: UpdateFilename :one
UPDATE files
SET filename = $1
WHERE id = $2
RETURNING *;


-- name: UpdateFileFolder :exec
UPDATE files
SET folder_id = $1
WHERE id = $2;


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

-- name: DeleteAllSharesForFile :exec
DELETE FROM file_shares
WHERE file_id = $1;

-- name: AddSharesToFile :copyfrom
INSERT INTO file_shares (file_id, shared_with)
VALUES ($1, $2);

-- name: ListFilesForUser :many
SELECT DISTINCT 
    f.id,
    f.filename,
    f.size,
    f.declared_mime AS content_type,
    f.uploaded_at,
    (f.owner_id = sqlc.arg(user_id)) AS user_owns_file,
    f.download_count
FROM files f
LEFT JOIN file_shares fs ON f.id = fs.file_id
WHERE 
    (f.owner_id = sqlc.arg(user_id) OR fs.shared_with = sqlc.arg(user_id))
    AND (sqlc.arg(filename)::TEXT = '' OR f.filename ILIKE '%' || sqlc.arg(filename)::TEXT || '%')
    AND (sqlc.arg(mime_type)::TEXT = '' OR f.declared_mime = sqlc.arg(mime_type)::TEXT)
    AND (sqlc.arg(uploaded_after)::TIMESTAMPTZ IS NULL OR f.uploaded_at > sqlc.arg(uploaded_after)::TIMESTAMPTZ)
    AND (sqlc.arg(uploaded_before)::TIMESTAMPTZ IS NULL OR f.uploaded_at < sqlc.arg(uploaded_before)::TIMESTAMPTZ)
    AND (
        sqlc.arg(ownership_status)::int = 0
        OR (sqlc.arg(ownership_status)::int = 1 AND f.owner_id = sqlc.arg(user_id))
        OR (sqlc.arg(ownership_status)::int = 2 AND f.owner_id <> sqlc.arg(user_id))
    )
ORDER BY f.uploaded_at DESC
LIMIT $1 OFFSET $2;

-----------------------------

-- name: ListFolderContents :many
WITH folder_contents AS (
    SELECT 
        f.id,
        f.name AS filename,
        'folder' AS item_type,
        NULL::bigint AS size,
        NULL::text AS content_type,
        f.created_at AS uploaded_at,
        (f.owner_id = sqlc.arg(user_id)) AS user_owns_file,
        NULL::bigint AS download_count,
        NULL::uuid AS folder_id
    FROM folders f
    WHERE 
        f.owner_id = sqlc.arg(user_id)
        AND f.parent_folder_id = sqlc.arg(parent_folder_id)::UUID
        AND (sqlc.arg(search)::TEXT = '' OR f.name ILIKE '%' || sqlc.arg(search)::TEXT || '%')
        AND (sqlc.arg(mime_type)::TEXT = 'folder/folder' OR sqlc.arg(mime_type)::TEXT = '')

    UNION ALL

    SELECT
        f.id,
        f.filename,
        'file' AS item_type,
        f.size,
        f.declared_mime AS content_type,
        f.uploaded_at,
        (f.owner_id = sqlc.arg(user_id)) AS user_owns_file,
        f.download_count,
        f.folder_id
    FROM files f
    WHERE
        f.owner_id = sqlc.arg(user_id)
        AND f.folder_id = sqlc.arg(parent_folder_id)::UUID
        AND (sqlc.arg(search)::TEXT = '' OR f.filename ILIKE '%' || sqlc.arg(search)::TEXT || '%')
        AND (sqlc.arg(mime_type)::TEXT = '' OR f.declared_mime = sqlc.arg(mime_type)::TEXT)
        AND (sqlc.arg(uploaded_after)::TIMESTAMPTZ IS NULL OR f.uploaded_at > sqlc.arg(uploaded_after)::TIMESTAMPTZ)
        AND (sqlc.arg(uploaded_before)::TIMESTAMPTZ IS NULL OR f.uploaded_at < sqlc.arg(uploaded_before)::TIMESTAMPTZ)
        AND (sqlc.narg(min_size)::BIGINT IS NULL OR f.size >= sqlc.narg(min_size)::BIGINT)
        AND (sqlc.narg(max_size)::BIGINT IS NULL OR f.size <= sqlc.narg(max_size)::BIGINT)
) 
SELECT *, COUNT(*) OVER() AS total_count 
FROM folder_contents
ORDER BY 
    item_type DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'filename' AND sqlc.arg(sort_order)::text = 'asc' THEN filename END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'filename' AND sqlc.arg(sort_order)::text = 'desc' THEN filename END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'size' AND sqlc.arg(sort_order)::text = 'asc' THEN size END ASC NULLS FIRST,
    CASE WHEN sqlc.arg(sort_by)::text = 'size' AND sqlc.arg(sort_order)::text = 'desc' THEN size END DESC NULLS LAST,
    CASE WHEN sqlc.arg(sort_by)::text = 'uploaded_at' AND sqlc.arg(sort_order)::text = 'asc' THEN uploaded_at END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'uploaded_at' AND sqlc.arg(sort_order)::text = 'desc' THEN uploaded_at END DESC
LIMIT $1 OFFSET $2;

-- name: ListRootContents :many
WITH root_contents AS (
    SELECT
        f.id, f.name AS filename, 'folder' AS item_type, NULL::bigint AS size,
        NULL::text AS content_type, f.created_at AS uploaded_at,
        (f.owner_id = sqlc.arg(user_id)) AS user_owns_file,
        NULL::bigint AS download_count, NULL::uuid AS folder_id
    FROM folders f
    WHERE f.owner_id = sqlc.arg(user_id) AND f.parent_folder_id IS NULL
      AND (sqlc.arg(search)::TEXT = '' OR f.name ILIKE '%' || sqlc.arg(search)::TEXT || '%')
      AND (sqlc.arg(mime_type)::TEXT = 'folder/folder' OR sqlc.arg(mime_type)::TEXT = '')

    UNION ALL

    SELECT
        f.id, f.filename, 'file' AS item_type, f.size, f.declared_mime AS content_type,
        f.uploaded_at, (f.owner_id = sqlc.arg(user_id)) AS user_owns_file,
        f.download_count, f.folder_id
    FROM files f
    WHERE f.owner_id = sqlc.arg(user_id) AND f.folder_id IS NULL
        AND (sqlc.arg(search)::TEXT = '' OR f.filename ILIKE '%' || sqlc.arg(search)::TEXT || '%')
        AND (sqlc.arg(mime_type)::TEXT = '' OR f.declared_mime = sqlc.arg(mime_type)::TEXT)
        AND (sqlc.arg(uploaded_after)::TIMESTAMPTZ IS NULL OR f.uploaded_at > sqlc.arg(uploaded_after)::TIMESTAMPTZ)
        AND (sqlc.arg(uploaded_before)::TIMESTAMPTZ IS NULL OR f.uploaded_at < sqlc.arg(uploaded_before)::TIMESTAMPTZ)
        AND (sqlc.narg(min_size)::BIGINT IS NULL OR f.size >= sqlc.narg(min_size)::BIGINT)
        AND (sqlc.narg(max_size)::BIGINT IS NULL OR f.size <= sqlc.narg(max_size)::BIGINT)
        AND (
            sqlc.arg(ownership_status)::int = 0
            OR (sqlc.arg(ownership_status)::int = 1 AND f.owner_id = sqlc.arg(user_id))
            OR (sqlc.arg(ownership_status)::int = 2 AND f.owner_id <> sqlc.arg(user_id))
          )

    UNION ALL

    SELECT
        f.id, f.filename, 'file' AS item_type, f.size, f.declared_mime AS content_type,
        f.uploaded_at, (f.owner_id = sqlc.arg(user_id)) AS user_owns_file,
        f.download_count, NULL::uuid as folder_id
    FROM files f
    JOIN file_shares fs ON f.id = fs.file_id
    WHERE fs.shared_with = sqlc.arg(user_id)
      AND (sqlc.arg(search)::TEXT = '' OR f.filename ILIKE '%' || sqlc.arg(search)::TEXT || '%')
      AND (sqlc.arg(mime_type)::TEXT = '' OR f.declared_mime = sqlc.arg(mime_type)::TEXT)
      AND (sqlc.arg(uploaded_after)::TIMESTAMPTZ IS NULL OR f.uploaded_at > sqlc.arg(uploaded_after)::TIMESTAMPTZ)
      AND (sqlc.arg(uploaded_before)::TIMESTAMPTZ IS NULL OR f.uploaded_at < sqlc.arg(uploaded_before)::TIMESTAMPTZ)
      AND (sqlc.narg(min_size)::BIGINT IS NULL OR f.size >= sqlc.narg(min_size)::BIGINT)
      AND (sqlc.narg(max_size)::BIGINT IS NULL OR f.size <= sqlc.narg(max_size)::BIGINT)
) 
SELECT *, COUNT(*) OVER() AS total_count 
FROM root_contents
ORDER BY
    item_type DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'filename' AND sqlc.arg(sort_order)::text = 'asc' THEN filename END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'filename' AND sqlc.arg(sort_order)::text = 'desc' THEN filename END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'size' AND sqlc.arg(sort_order)::text = 'asc' THEN size END ASC NULLS FIRST,
    CASE WHEN sqlc.arg(sort_by)::text = 'size' AND sqlc.arg(sort_order)::text = 'desc' THEN size END DESC NULLS LAST,
    CASE WHEN sqlc.arg(sort_by)::text = 'uploaded_at' AND sqlc.arg(sort_order)::text = 'asc' THEN uploaded_at END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'uploaded_at' AND sqlc.arg(sort_order)::text = 'desc' THEN uploaded_at END DESC
LIMIT $1 OFFSET $2;


-- name: IncrementFileDownloadCount :exec
UPDATE files
SET download_count = download_count + 1
WHERE id = $1;

-- name: GetBlobIDsInFolderHierarchy :many
WITH RECURSIVE folder_hierarchy AS (
    -- This part is correct and finds all sub-folder IDs
    SELECT fo.id FROM folders fo WHERE fo.id = $1
    UNION ALL
    SELECT f.id FROM folders f
    INNER JOIN folder_hierarchy fh ON f.parent_folder_id = fh.id
)
SELECT DISTINCT f.blob_id
FROM files f
WHERE f.folder_id IN (SELECT id FROM folder_hierarchy);

-- name: ListAllFiles :many
SELECT
    f.id,
    f.filename,
    f.size,
    f.declared_mime,
    f.uploaded_at,
    f.download_count,
    f.owner_id,
    u.email as owner_email,
    COUNT(*) OVER() AS total_count
FROM
    files f
JOIN
    users u ON f.owner_id = u.id 
ORDER BY 
    CASE WHEN sqlc.arg(sort_order)::text = 'asc' THEN
        CASE sqlc.arg(sort_by)::text
            WHEN 'filename' THEN f.filename::text
            WHEN 'owner_email' THEN u.email::text
            -- Pad numbers to ensure correct alphabetical sorting
            WHEN 'size' THEN LPAD(f.size::text, 20, '0')
            WHEN 'download_count' THEN LPAD(f.download_count::text, 20, '0')
            -- Timestamps in ISO format sort correctly as text
            WHEN 'uploaded_at' THEN f.uploaded_at::text
            ELSE f.uploaded_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg(sort_order)::text = 'desc' THEN
        CASE sqlc.arg(sort_by)::text
            WHEN 'filename' THEN f.filename::text
            WHEN 'owner_email' THEN u.email::text
            WHEN 'size' THEN LPAD(f.size::text, 20, '0')
            WHEN 'download_count' THEN LPAD(f.download_count::text, 20, '0')
            WHEN 'uploaded_at' THEN f.uploaded_at::text
            ELSE f.uploaded_at::text
        END
    END DESC 
LIMIT $1 OFFSET $2;
