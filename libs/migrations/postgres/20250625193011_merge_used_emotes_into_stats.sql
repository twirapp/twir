-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE users_stats ADD COLUMN "emotes" int4 NOT NULL DEFAULT 0;

UPDATE users_stats
SET "emotes" = (
	SELECT COUNT(*)
	FROM channels_emotes_usages
	WHERE "channelId" = users_stats."channelId" AND "userId" = users_stats."userId"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
