-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_dota_settings
    ALTER COLUMN commands_settings SET DEFAULT '{"mmr":true,"wl":true,"lg":true,"gm":true,"np":true,"wp":true}'::jsonb;

UPDATE channels_dota_settings
SET commands_settings = '{"mmr":true,"wl":true,"lg":true,"gm":true,"np":true,"wp":true}'::jsonb || commands_settings
WHERE NOT (commands_settings ?& ARRAY['mmr', 'wl', 'lg', 'gm', 'np', 'wp']);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_dota_settings
    ALTER COLUMN commands_settings SET DEFAULT '{}'::jsonb;
-- +goose StatementEnd
