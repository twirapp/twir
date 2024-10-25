-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_commands ADD COLUMN expires_at timestamptz;
ALTER TABLE channels_commands ADD COLUMN expired boolean default false;
ALTER TABLE channels_commands ADD COLUMN expired_in integer DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
