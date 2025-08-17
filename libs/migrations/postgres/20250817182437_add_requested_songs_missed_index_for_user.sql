-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX IF NOT EXISTS channels_requested_songs_orderedbyid_idx ON channels_requested_songs("orderedById");
CREATE INDEX IF NOT EXISTS channels_requested_songs_channelid_idx ON channels_requested_songs("channelId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
