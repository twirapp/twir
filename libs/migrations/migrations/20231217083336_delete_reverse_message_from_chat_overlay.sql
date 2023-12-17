-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
UPDATE "channels_modules_settings"
SET settings = settings::jsonb - 'reverseMessages'
WHERE "type" = 'chat_overlay';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
