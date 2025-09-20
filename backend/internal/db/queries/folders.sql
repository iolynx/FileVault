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
RETURNING *;

-- name: DeleteFolder :exec
DELETE FROM folders
WHERE id = $1;
