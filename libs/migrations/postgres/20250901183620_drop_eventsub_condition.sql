-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE eventsub_topics DROP COLUMN condition_type;

INSERT INTO eventsub_topics (topic, version)
VALUES ('channel.raid', '1'),
			 ('channel.channel_points_custom_reward.add', '1'),
			 ('channel.channel_points_custom_reward.remove', '1'),
			 ('channel.channel_points_custom_reward.update', '1'),
			 ('channel.vip.add', '1'),
			 ('channel.vip.remove', '1')
ON CONFLICT (topic) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
