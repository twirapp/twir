-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS channels_integrations_faceit
(
	id           SERIAL PRIMARY KEY,
	channel_id   TEXT        NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
	access_token TEXT        NOT NULL,
	username     TEXT        NOT NULL,
	avatar       TEXT        NOT NULL,
	game         TEXT        NOT NULL DEFAULT '',
	enabled      BOOLEAN     NOT NULL DEFAULT TRUE,
	created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (channel_id)
);

-- Migrate data from channels_integrations table
INSERT INTO channels_integrations_faceit (channel_id,
																					access_token,
																					username,
																					avatar,
																					game,
																					created_at,
																					updated_at,
																					enabled)
SELECT ci."channelId",
			 ci."accessToken",
			 ci.data ->> 'username' AS username,
			 ci.data ->> 'avatar'    AS avatar,
			 COALESCE(ci.data ->> 'game', 'cs2') AS game,
			 NOW()                   AS created_at,
			 NOW()                   AS updated_at,
			 ci.enabled
FROM public.channels_integrations ci
WHERE ci."integrationId" = (SELECT id FROM public.integrations WHERE service = 'FACEIT' LIMIT 1)
	AND NOT EXISTS (SELECT 1
									FROM channels_integrations_faceit cif
									WHERE cif.channel_id = ci."channelId")
	AND ci."accessToken" IS NOT NULL
	AND ci."accessToken" != ''
	AND ci.data ->> 'username' IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS channels_integrations_faceit;
-- +goose StatementEnd
