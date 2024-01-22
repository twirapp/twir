-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_integrations_seventv" (
		"id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
		"channel_id" text NOT NULL,
		"reward_id_for_add_emote" uuid,
		"reward_id_for_remove_emote" uuid,
		CONSTRAINT "channels_integrations_seventv_channel_id_fk" FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
