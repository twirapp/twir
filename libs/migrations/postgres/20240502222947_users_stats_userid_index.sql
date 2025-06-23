-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX users_userid_idx ON users("id");
CREATE INDEX channels_channelid_idx ON channels("id");
CREATE INDEX users_stats_userid_idx ON users_stats("userId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
