ALTER TABLE users DROP COLUMN IF EXISTS original_storage_bytes;
ALTER TABLE users DROP COLUMN IF EXISTS dedup_storage_bytes;

ALTER TABLE users ADD COLUMN storage_used BIGINT NOT NULL DEFAULT 0;

CREATE INDEX idx_users_storage_used ON users(storage_used);
