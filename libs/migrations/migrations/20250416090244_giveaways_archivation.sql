-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_giveaways ADD COLUMN archived_at TIMESTAMPTZ;
ALTER TABLE channels_giveaways ADD COLUMN is_archived BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
