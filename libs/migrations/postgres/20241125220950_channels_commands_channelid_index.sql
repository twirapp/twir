-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX channels_commands_channelid_index ON channels_commands("channelId");
CREATE INDEX channels_commands_responses_commandid_index ON channels_commands_responses("commandId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
