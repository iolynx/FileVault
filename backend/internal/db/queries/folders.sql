-- name: CreateFolder :one
INSERT INTO folders (
    name,
    owner_id,
    parent_folder_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetFolderByID :one
SELECT * FROM folders
WHERE id = $1;

-- name: UpdateFolder :one
UPDATE folders
SET 
    name = $2,
    parent_folder_id = $3
WHERE 
    id = $1
RETURNING id, name as filename, owner_id, parent_folder_id, created_at;

-- name: DeleteFolder :exec
DELETE FROM folders
WHERE id = $1;

-- name: UpdateFolderParentFolder :exec
UPDATE folders
SET parent_folder_id = $1
WHERE id = $2
RETURNING *;

-- name: ListSelectableFolders :many
WITH RECURSIVE forbidden_folders AS (
    SELECT id FROM folders WHERE id = sqlc.narg('current_folder_id')::uuid

    UNION ALL

    -- find all children of the folders already in our set.
    SELECT f.id
    FROM folders f
    INNER JOIN forbidden_folders ff ON f.parent_folder_id = ff.id
)
SELECT
    f.id,
    f.name,
    f.created_at,
    f.parent_folder_id
FROM folders f
WHERE
    f.owner_id = $1
    -- exclude all folders that are in the forbidden list
    AND f.id NOT IN (SELECT id FROM forbidden_folders)
ORDER BY
    f.created_at DESC;
