ALTER TABLE users
  DROP COLUMN IF EXISTS original_storage_bytes,
  DROP COLUMN IF EXISTS dedup_storage_bytes;
