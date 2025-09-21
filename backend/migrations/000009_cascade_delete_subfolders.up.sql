ALTER TABLE folders DROP CONSTRAINT IF EXISTS folders_parent_folder_id_fkey;

ALTER TABLE folders
ADD CONSTRAINT folders_parent_folder_id_fkey
FOREIGN KEY (parent_folder_id) REFERENCES folders(id) ON DELETE CASCADE;
