CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    original_storage_bytes BIGINT NOT NULL DEFAULT 0,
    dedup_storage_bytes BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE blobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sha256 TEXT UNIQUE NOT NULL,
  storage_path TEXT NOT NULL,
  size BIGINT NOT NULL,
  mime_type TEXT,
  refcount INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  blob_id UUID NOT NULL REFERENCES blobs(id) ON DELETE RESTRICT,
  filename TEXT NOT NULL,
  declared_mime TEXT,
  size BIGINT NOT NULL,
  uploaded_at TIMESTAMPTZ DEFAULT now(),
  is_public BOOLEAN DEFAULT FALSE,
  public_token UUID,
  download_count BIGINT DEFAULT 0
);

CREATE TABLE file_shares (
    id BIGSERIAL PRIMARY KEY,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    shared_with BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission TEXT NOT NULL DEFAULT 'read',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(file_id, shared_with)
);


CREATE INDEX idx_blobs_sha256 ON blobs(sha256);
CREATE INDEX idx_files_owner ON files(owner_id);
CREATE INDEX idx_files_owner_filename
ON files(owner_id, filename);

