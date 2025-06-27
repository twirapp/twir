-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE channels_commands_prefix (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		channel_id TEXT NOT NULL,
		prefix TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_channels_commands_prefix_channel_id ON channels_commands_prefix (channel_id);

ALTER TABLE channels_commands ADD COLUMN prefix TEXT;

ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_BOT_SETTINGS';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_BOT_SETTINGS';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
