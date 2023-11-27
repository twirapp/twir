-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DELETE FROM "channels_moderation_settings" WHERE "type" = 'language';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
