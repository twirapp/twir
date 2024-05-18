-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE "channels_games_voteban_voting_mode" AS ENUM ('chat', 'twitch_polls');

CREATE TABLE "channels_games_voteban"
(
	"id"                         UUID PRIMARY KEY                            DEFAULT gen_random_uuid(),
	"channel_id"                 TEXT                               NOT NULL references "channels" ("id") ON DELETE CASCADE,
	"enabled"                    BOOL                               NOT NULL,
	"timeout_seconds"            INT2                               NOT NULL,
	"timeout_moderators"         BOOL                               NOT NULL,
	"ban_message"                varchar(500)                       NOT NULL,
	"ban_message_moderators"     varchar(500)                       NOT NULL,
	"survive_message"            varchar(500)                       NOT NULL,
	"survive_message_moderators" varchar(500)                       NOT NULL,
	"init_message"               varchar(500)                       NOT NULL,
	"needed_votes"               INT                                NOT NULL,
	"vote_duration"              INT                                NOT NULL,
	"voting_mode"                channels_games_voteban_voting_mode NOT NULL DEFAULT 'chat',
	"chat_votes_words_positive"  TEXT[]                             NOT NULL default '{}',
	"chat_votes_words_negative"  TEXT[]                             NOT NULL default '{}'
);

CREATE UNIQUE INDEX "channels_games_voteban_channel_id_index" ON "channels_games_voteban" ("channel_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
