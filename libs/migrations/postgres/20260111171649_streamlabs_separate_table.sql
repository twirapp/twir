-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS channels_integrations_streamlabs (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	access_token TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	username TEXT NOT NULL,
	avatar TEXT NOT NULL,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE(channel_id)
);

INSERT INTO channels_integrations_streamlabs (
	channel_id,
	access_token,
	refresh_token,
	username,
	avatar,
	created_at,
	updated_at,
	enabled
)
SELECT
	ci."channelId",
	ci."accessToken",
	ci."refreshToken",
	ci.data->>'userName' AS username,
	ci.data->>'avatar' AS avatar,
	NOW() AS created_at,
	NOW() AS updated_at,
	ci.enabled
FROM public.channels_integrations ci
WHERE ci."integrationId" = (SELECT id FROM public.integrations WHERE service = 'STREAMLABS' LIMIT 1)
	AND NOT EXISTS (
	SELECT 1
	FROM channels_integrations_streamlabs cis
	WHERE cis.channel_id = ci."channelId"
) AND ci."refreshToken" IS NOT NULL AND ci."accessToken" IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS channels_integrations_streamlabs;
-- +goose StatementEnd
