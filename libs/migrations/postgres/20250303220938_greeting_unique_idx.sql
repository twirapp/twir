-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE UNIQUE INDEX channels_greetings_userid_channelid_unique_idx ON channels_greetings ("userId", "channelId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
