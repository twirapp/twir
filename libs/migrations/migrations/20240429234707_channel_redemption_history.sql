-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channel_redemptions_history" (
	id uuid primary key unique,
	channel_id text not null,
	user_id text not null,
	reward_id uuid not null,
	reward_title text not null,
	reward_prompt text,
	reward_cost integer not null,
	redeemed_at timestamp not null,
	FOREIGN KEY (channel_id) REFERENCES channels(id)
);

CREATE INDEX "channel_redemptions_history_channel_id_idx" ON "channel_redemptions_history" (channel_id);
CREATE INDEX "channel_redemptions_history_user_id_idx" ON "channel_redemptions_history" (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
