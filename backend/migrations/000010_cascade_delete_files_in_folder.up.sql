ALTER TABLE files DROP CONSTRAINT IF EXISTS files_folder_id_fkey;

ALTER TABLE files
ADD CONSTRAINT files_folder_id_fkey
FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE;
