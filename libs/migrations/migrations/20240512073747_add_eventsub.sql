-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE eventsub_topics_condition_type_enum AS ENUM (
	'BROADCASTER_USER_ID',
	'USER_ID',
	'BROADCASTER_WITH_USER_ID',
	'BROADCASTER_WITH_MODERATOR_ID',
	'TO_BROADCASTER_ID'
);

CREATE TABLE "eventsub_topics"
(
	"id"             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	"topic"          VARCHAR(255)                        NOT NULL,
	"version"        VARCHAR(255)                        NOT NULL,
	"condition_type" eventsub_topics_condition_type_enum NOT NULL
);

CREATE TABLE "eventsub_subscriptions"
(
	"id"           UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
	"topic_id"     UUID         NOT NULL references eventsub_topics (id),
	"user_id"      text         NOT NULL references users (id),
	"status"       VARCHAR(255) NOT NULL,
	"created_at"   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"version"      VARCHAR(255) NOT NULL,
	"callback_url" VARCHAR(255) NOT NULL
);

CREATE INDEX "eventsub_topics_topic_idx" ON "eventsub_topics" ("topic");
CREATE INDEX "eventsub_subscriptions_user_id_idx" ON "eventsub_subscriptions" ("user_id");

INSERT INTO eventsub_topics (topic, version, condition_type)
VALUES ('channel.update', '2', 'BROADCASTER_USER_ID'),
			 ('stream.online', '1', 'BROADCASTER_USER_ID'),
			 ('stream.offline', '1', 'BROADCASTER_USER_ID'),
			 ('user.update', '1', 'USER_ID'),
			 ('channel.follow', '2', 'BROADCASTER_WITH_MODERATOR_ID'),
			 ('channel.moderator.add', '1', 'BROADCASTER_USER_ID'),
			 ('channel.moderator.remove', '1', 'BROADCASTER_USER_ID'),
			 ('channel.channel_points_custom_reward_redemption.add', '1', 'BROADCASTER_USER_ID'),
			 ('channel.channel_points_custom_reward_redemption.update', '1', 'BROADCASTER_USER_ID'),
			 ('channel.poll.begin', '1', 'BROADCASTER_USER_ID'),
			 ('channel.poll.progress', '1', 'BROADCASTER_USER_ID'),
			 ('channel.poll.end', '1', 'BROADCASTER_USER_ID'),
			 ('channel.prediction.begin', '1', 'BROADCASTER_USER_ID'),
			 ('channel.prediction.lock', '1', 'BROADCASTER_USER_ID'),
			 ('channel.prediction.progress', '1', 'BROADCASTER_USER_ID'),
			 ('channel.prediction.end', '1', 'BROADCASTER_USER_ID'),
			 ('channel.ban', '1', 'BROADCASTER_USER_ID'),
			 ('channel.subscribe', '1', 'BROADCASTER_USER_ID'),
			 ('channel.subscription.gift', '1', 'BROADCASTER_USER_ID'),
			 ('channel.subscription.message', '1', 'BROADCASTER_USER_ID'),
			 ('channel.raid', '1', 'TO_BROADCASTER_ID'),
			 ('channel.chat.clear', '1', 'BROADCASTER_WITH_USER_ID'),
			 ('channel.chat.clear_user_messages', '1', 'BROADCASTER_WITH_USER_ID'),
			 ('channel.chat.message_delete', '1', 'BROADCASTER_WITH_USER_ID'),
			 ('channel.chat.notification', '1', 'BROADCASTER_WITH_USER_ID'),
			 ('channel.chat.message', '1', 'BROADCASTER_WITH_USER_ID'),
			 ('channel.unban_request.create', '1', 'BROADCASTER_WITH_MODERATOR_ID'),
			 ('channel.unban_request.resolve', '1', 'BROADCASTER_WITH_MODERATOR_ID');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
