CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    storage_quota BIGINT NOT NULL DEFAULT 10000000,
    storage_used BIGINT NOT NULL DEFAULT 0
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
  download_count BIGINT DEFAULT 0,
  folder_id UUID REFERENCES folders(id) ON DELETE SET NULL
);

CREATE TABLE file_shares (
    id BIGSERIAL PRIMARY KEY,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    shared_with BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission TEXT NOT NULL DEFAULT 'read',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(file_id, shared_with)
);

CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_folder_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action audit_action NOT NULL,
    target_id UUID,
    details JSONB, -- A flexible column to store any extra, machine-readable details.
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE audit_action AS ENUM (
    'USER_REGISTERED',
    'USER_LOGGED_IN',
    'FILE_UPLOADED',
    'FILE_DOWNLOADED',
    'FILE_RENAMED',
    'FILE_DELETED'
);

CREATE INDEX idx_blobs_sha256 ON blobs(sha256);
CREATE INDEX idx_files_owner ON files(owner_id);
CREATE INDEX idx_files_owner_filename ON files(owner_id, filename);
CREATE INDEX idx_folders_owner_id_parent_id ON folders(owner_id, parent_folder_id);
CREATE INDEX idx_files_folder_id ON files(folder_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
