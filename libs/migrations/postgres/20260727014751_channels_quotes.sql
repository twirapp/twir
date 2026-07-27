-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_quotes (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	"channelId" TEXT NOT NULL REFERENCES channels ON UPDATE CASCADE ON DELETE RESTRICT,
	number INTEGER NOT NULL,
	text TEXT NOT NULL,
	"creatorId" TEXT,
	"creatorName" TEXT,
	"gameId" TEXT,
	"gameName" TEXT,
	"createdAt" TIMESTAMP NOT NULL DEFAULT now(),
	"updatedAt" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX channels_quotes_channel_id_number_key
	ON channels_quotes ("channelId", number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels_quotes;
-- +goose StatementEnd
