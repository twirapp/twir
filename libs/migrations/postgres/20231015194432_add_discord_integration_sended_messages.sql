-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO "integrations" ("service") VALUES ('DISCORD');
CREATE TABLE discord_sended_notifications (
	id TEXT PRIMARY KEY DEFAULT uuid_generate_v4(),
	guild_id TEXT NOT NULL,
	message_id TEXT NOT NULL,
	channel_id TEXT references channels(id) NOT NULL,
	discord_channel_id TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT now(),
	updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
