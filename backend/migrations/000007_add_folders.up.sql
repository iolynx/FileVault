CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_folder_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE files
ADD COLUMN folder_id UUID REFERENCES folders(id) ON DELETE SET NULL;

CREATE INDEX idx_folders_owner_id_parent_id ON folders(owner_id, parent_folder_id);
CREATE INDEX idx_files_folder_id ON files(folder_id);

COMMENT ON COLUMN folders.parent_folder_id IS 'NULL indicates a root folder';
