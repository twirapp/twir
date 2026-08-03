-- +goose Up
-- +goose StatementBegin
ALTER TABLE notifications
	ADD COLUMN discord_message_id TEXT,
	ADD COLUMN discord_channel_id TEXT,
	ADD COLUMN discord_attachment_keys TEXT[] NOT NULL DEFAULT '{}',
	ADD COLUMN updated_at TIMESTAMPTZ;

CREATE UNIQUE INDEX notifications_discord_message_id_idx
	ON notifications (discord_message_id)
	WHERE discord_message_id IS NOT NULL;

ALTER TABLE users_viewed_notifications
	DROP CONSTRAINT "FK_f5d19d90314d14d636752e2888b",
	ADD CONSTRAINT "FK_f5d19d90314d14d636752e2888b"
		FOREIGN KEY ("notificationId")
		REFERENCES notifications (id)
		ON UPDATE CASCADE
		ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users_viewed_notifications
	DROP CONSTRAINT "FK_f5d19d90314d14d636752e2888b",
	ADD CONSTRAINT "FK_f5d19d90314d14d636752e2888b"
		FOREIGN KEY ("notificationId")
		REFERENCES notifications (id)
		ON UPDATE CASCADE
		ON DELETE RESTRICT;

DROP INDEX notifications_discord_message_id_idx;

ALTER TABLE notifications
	DROP COLUMN updated_at,
	DROP COLUMN discord_attachment_keys,
	DROP COLUMN discord_channel_id,
	DROP COLUMN discord_message_id;
-- +goose StatementEnd
