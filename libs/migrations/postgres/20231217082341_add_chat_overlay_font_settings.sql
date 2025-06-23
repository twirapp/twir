-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
UPDATE "channels_modules_settings"
SET settings = settings || '{"fontFamily": "roboto", "fontWeight": 400, "fontStyle": "normal"}'
WHERE "type" = 'chat_overlay';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
