-- +goose Up
-- +goose NO TRANSACTION
ALTER TYPE platform ADD VALUE IF NOT EXISTS 'youtube';

-- +goose Down
-- +goose NO TRANSACTION
-- PostgreSQL does not support removing values from an enum type; left as a no-op.
SELECT 1;
