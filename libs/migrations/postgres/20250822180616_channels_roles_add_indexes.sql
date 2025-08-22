-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX IF NOT EXISTS channels_roles_channel_id_idx ON channels_roles("channelId");
CREATE UNIQUE INDEX IF NOT EXISTS channels_roles_roleid_userid_idx ON channels_roles_users("roleId", "userId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
