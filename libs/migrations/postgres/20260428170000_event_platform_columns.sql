-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_events
	ADD COLUMN IF NOT EXISTS platforms platform[] NOT NULL DEFAULT '{}'::platform[];

ALTER TABLE channels_events_list
	ADD COLUMN IF NOT EXISTS platform platform NOT NULL DEFAULT 'twitch';

ALTER TABLE channels_info_history
	ADD COLUMN IF NOT EXISTS platform platform NOT NULL DEFAULT 'twitch';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_info_history
	DROP COLUMN IF EXISTS platform;

ALTER TABLE channels_events_list
	DROP COLUMN IF EXISTS platform;

ALTER TABLE channels_events
	DROP COLUMN IF EXISTS platforms;
-- +goose StatementEnd
