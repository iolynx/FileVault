-- blobs table: content-addressed storage metadata
CREATE TABLE blobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sha256 TEXT UNIQUE NOT NULL,
  storage_path TEXT NOT NULL,
  size BIGINT NOT NULL,
  mime_type TEXT,
  refcount INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- files table: per-user file records referencing blobs
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

CREATE INDEX idx_blobs_sha256 ON blobs(sha256);
CREATE INDEX idx_files_owner ON files(owner_id);

