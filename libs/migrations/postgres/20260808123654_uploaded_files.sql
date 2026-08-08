-- +goose Up
-- +goose StatementBegin
CREATE TABLE uploaded_files (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    public_id TEXT NOT NULL,
    uploaded_by_user_id UUID,
    file_name TEXT,
    mime_type TEXT NOT NULL,
    extension TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    s3_key TEXT NOT NULL,
    delete_key TEXT NOT NULL,
    user_agent TEXT,
    user_ip cidr,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (uploaded_by_user_id) REFERENCES users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX uploaded_files_public_id_idx ON uploaded_files(public_id);
CREATE INDEX uploaded_files_uploaded_by_user_id_idx ON uploaded_files(uploaded_by_user_id);
CREATE INDEX uploaded_files_expires_at_idx ON uploaded_files(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE uploaded_files;
-- +goose StatementEnd
