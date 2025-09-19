ALTER TABLE users
  ADD COLUMN original_storage_bytes BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN dedup_storage_bytes BIGINT NOT NULL DEFAULT 0;
