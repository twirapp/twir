-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE channels_modules_settings_type_enum ADD VALUE 'russian_roulette';
DELETE FROM "channels_commands" WHERE "name" = 'roulette';
DELETE FROM "channels_commands" WHERE "aliases" @> ARRAY['roulette'];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
