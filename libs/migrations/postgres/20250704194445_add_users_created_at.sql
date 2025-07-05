-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE IF EXISTS audit_logs;

ALTER TABLE users ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE channels ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE users_stats ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE users_stats ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_stats_updated_at
	BEFORE UPDATE ON users_stats
	FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
