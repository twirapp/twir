-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE "channels_moderation_settings";
DROP TABLE "channels_moderation_warnings";
DROP TYPE "channels_moderation_settings_type_enum";

CREATE TYPE "channels_moderation_settings_type_enum" AS ENUM (
	'links',
	'deny_list',
	'symbols',
	'long_message',
	'caps',
	'emotes',
	'language'
);

CREATE TABLE "channels_moderation_settings" (
		"id" uuid NOT NULL default uuid_generate_v4(),
		"channel_id" text NOT NULL references "channels" ("id") ON DELETE CASCADE,
		"type" channels_moderation_settings_type_enum NOT NULL,
		"created_at" timestamp NOT NULL default now(),
		"updated_at" timestamp NOT NULL default now(),
		"enabled" boolean NOT NULL default false,
		"max_warnings" integer NOT NULL default 0,
		"ban_time" integer NOT NULL default 0,
		"ban_message" varchar(500) NOT NULL default '',
		"warning_message" varchar(500) NOT NULL default '',
		"check_clips" boolean NOT NULL default false,
		"trigger_length" integer NOT NULL default 0,
		"max_percentage" integer NOT NULL default 0,
		"deny_list" text[] default '{}',
		"denied_chat_languages" text[] default '{}',
		"excluded_roles" text[] default '{}',
		CONSTRAINT "channels_moderation_settings_pkey" PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
