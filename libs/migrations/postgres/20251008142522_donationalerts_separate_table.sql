-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS channels_integrations_donationalerts (
	id SERIAL PRIMARY KEY,
	public_id ulid NOT NULL DEFAULT gen_ulid(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	access_token TEXT,
	refresh_token TEXT,
	username TEXT,
	avatar TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE(channel_id)
);

INSERT INTO channels_integrations_donationalerts (
	channel_id,
	access_token,
	refresh_token,
	username,
	avatar,
	created_at,
	updated_at
)
SELECT
	ci."channelId",
	ci."accessToken",
	ci."refreshToken",
	ci.data->'userName'->>'userName' AS username,
	ci.data->>'avatar' AS avatar,
	NOW() AS created_at,
	NOW() AS updated_at
FROM public.channels_integrations ci
WHERE ci."integrationId" = (SELECT id FROM public.integrations WHERE service = 'DONATIONALERTS' LIMIT 1)
	AND NOT EXISTS (
	SELECT 1
	FROM channels_integrations_donationalerts cid
	WHERE cid.channel_id = ci."channelId"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
