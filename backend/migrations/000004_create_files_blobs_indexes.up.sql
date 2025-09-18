CREATE INDEX IF NOT EXISTS idx_blobs_sha256 ON blobs(sha256);
CREATE INDEX IF NOT EXISTS idx_files_owner ON files(owner_id);
CREATE INDEX IF NOT EXISTS idx_files_owner_filename
ON files(owner_id, filename);
