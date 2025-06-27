-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE chat_messages;
DROP TABLE channels_emotes_usages;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
