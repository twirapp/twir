-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

UPDATE channels_commands SET name = '7tv_profile_twir_migration' WHERE name = '7tv profile';
UPDATE channels_commands SET name = '7tv_emote_twir_migration' WHERE name = '7tv emote';
UPDATE channels_commands SET name = '7tv_rename_twir_migration' WHERE name = '7tv rename';
UPDATE channels_commands SET name = '7tv_delete_twir_migration' WHERE name = '7tv delete';
UPDATE channels_commands SET name = '7tv_add_twir_migration' WHERE name = '7tv add';
ALTER TYPE "channels_commands_module_enum" ADD VALUE '7tv';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
