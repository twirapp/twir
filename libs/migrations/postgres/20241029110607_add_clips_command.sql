-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_commands_module_enum ADD VALUE 'CLIPS';

UPDATE "channels_commands" SET "name" = 'clip_twir_add_own_module' WHERE "name" = 'clip';
UPDATE "channels_commands" SET "aliases" = array_replace("aliases", 'clip', 'clip_twir_add_own_module');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
