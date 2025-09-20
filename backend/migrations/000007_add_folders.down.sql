DROP INDEX IF EXISTS idx_files_folder_id;
DROP INDEX IF EXISTS idx_folders_owner_id_parent_id;

ALTER TABLE files
DROP COLUMN IF EXISTS folder_id;

DROP TABLE IF EXISTS folders;
