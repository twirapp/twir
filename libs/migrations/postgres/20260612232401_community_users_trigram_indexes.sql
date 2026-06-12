-- +goose NO TRANSACTION
-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_login_trgm
	ON users USING gin (login gin_trgm_ops);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_display_name_trgm
	ON users USING gin (display_name gin_trgm_ops);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_users_login_trgm;
DROP INDEX CONCURRENTLY IF EXISTS idx_users_display_name_trgm;
