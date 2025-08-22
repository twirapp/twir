-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX IF NOT EXISTS idx_users_online_channel_id ON users_online("channelId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
