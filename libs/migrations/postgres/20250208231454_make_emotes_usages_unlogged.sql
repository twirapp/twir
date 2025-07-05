-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_emotes_usages SET UNLOGGED;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
