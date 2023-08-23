-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE table "channels_alerts"
(
	id         uuid default uuid_generate_v4() not null primary key,
	channel_id text                            not null
		constraint "channels_alerts_channel_id" references channels ("id")
			on update cascade
			on delete cascade,
	audio_id   uuid                            null
		constraint "channels_alerts_audio_id" references channels_files ("id")
			on update cascade
			on delete SET NULL,
	name       text                            not null
);
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_ALERTS';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_ALERTS';

ALTER table "channels_commands"
	ADD COLUMN "alert_id" uuid null
		constraint "channels_commands_alert_id" references channels_alerts ("id")
			on update cascade
			on delete set null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE
FROM pg_enum
WHERE enumlabel = 'channels_roles_permissions_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'VIEW_ALERTS');

DELETE
FROM pg_enum
WHERE enumlabel = 'channels_roles_permissions_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'MANAGE_ALERTS');
-- +goose StatementEnd
