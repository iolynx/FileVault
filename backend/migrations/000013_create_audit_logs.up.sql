CREATE TYPE audit_action AS ENUM (
    'USER_REGISTERED',
    'USER_LOGGED_IN',
    'FILE_UPLOADED',
    'FILE_DOWNLOADED',
    'FILE_RENAMED',
    'FILE_DELETED'
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action audit_action NOT NULL,
    target_id UUID,
    details JSONB, -- A flexible column to store any extra, machine-readable details.
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
