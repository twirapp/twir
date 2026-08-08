-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at
	ON users (created_at);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_users_created_at;
