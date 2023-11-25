-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE "channels_modules_settings_type_enum" ADD VALUE 'be_right_back_overlay';
ALTER TYPE "channels_commands_module_enum" ADD VALUE 'OVERLAYS';

UPDATE "channels_commands" SET "name" = 'brb_twir_add_own_module' WHERE "name" = 'brb';
UPDATE "channels_commands" SET "aliases" = array_replace("aliases", 'brb', 'brb_twir_add_own_module');

UPDATE "channels_commands" SET "name" = 'brbstop_twir_add_own_module' WHERE "name" = 'brbstop';
-- same for aliases array

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
