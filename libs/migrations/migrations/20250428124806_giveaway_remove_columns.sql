-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_giveaways DROP COLUMN ended_at;
ALTER TABLE channels_giveaways DROP COLUMN archived_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
