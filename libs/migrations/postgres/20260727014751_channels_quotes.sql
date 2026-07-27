-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_quotes (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id UUID NOT NULL REFERENCES channels ON UPDATE CASCADE ON DELETE RESTRICT,
	number INTEGER NOT NULL,
	text TEXT NOT NULL,
	creator_id TEXT,
	creator_name TEXT,
	game_id TEXT,
	game_name TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX channels_quotes_channel_id_number_key
	ON channels_quotes (channel_id, number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels_quotes;
-- +goose StatementEnd
