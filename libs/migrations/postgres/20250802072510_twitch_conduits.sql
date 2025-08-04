-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE twitch_conduits (
	id text NOT NULL PRIMARY KEY,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	shard_count int2 NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
