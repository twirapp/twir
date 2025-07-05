-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO eventsub_topics (topic, version, condition_type)
VALUES ('channel.channel_points_custom_reward.add', '1', 'BROADCASTER_USER_ID'),
			 ('channel.channel_points_custom_reward.remove', '1', 'BROADCASTER_USER_ID'),
			 ('channel.channel_points_custom_reward.update', '1', 'BROADCASTER_USER_ID');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
