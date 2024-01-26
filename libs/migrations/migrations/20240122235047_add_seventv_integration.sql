-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_integrations_seventv"
(
	"id"                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	"channel_id"                      text             NOT NULL,
	"reward_id_for_add_emote"         text,
	"reward_id_for_remove_emote"      text,
	"delete_emotes_only_added_by_app" boolean          NOT NULL DEFAULT true,
	"added_emotes"                    text[]           NOT NULL DEFAULT '{}',
	CONSTRAINT "channels_integrations_seventv_channel_id_fk" FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE
);

CREATE INDEX "channels_integrations_seventv_channel_id_idx" ON "channels_integrations_seventv" ("channel_id");
CREATE INDEX "channels_integrations_seventv_reward_id_for_add_emote_idx" ON "channels_integrations_seventv" ("reward_id_for_add_emote");
CREATE INDEX "channels_integrations_seventv_reward_id_for_remove_emote_idx" ON "channels_integrations_seventv" ("reward_id_for_remove_emote");

ALTER TYPE channels_events_operations_type_enum ADD VALUE 'SEVENTV_ADD_EMOTE';
ALTER TYPE channels_events_operations_type_enum ADD VALUE 'SEVENTV_REMOVE_EMOTE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
