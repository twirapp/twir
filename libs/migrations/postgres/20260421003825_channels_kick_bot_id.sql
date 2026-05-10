-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels
    ADD COLUMN kick_bot_id UUID NULL REFERENCES kick_bots(id) ON DELETE SET NULL ON UPDATE CASCADE;

CREATE INDEX channels_kick_bot_id_idx ON channels(kick_bot_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channels_kick_bot_id_idx;
ALTER TABLE channels DROP COLUMN IF EXISTS kick_bot_id;
-- +goose StatementEnd
