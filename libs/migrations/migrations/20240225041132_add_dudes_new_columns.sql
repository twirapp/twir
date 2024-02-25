-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_overlays_dudes_user_settings_dude_sprite ADD VALUE 'random';
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_max_on_screen" integer NOT NULL DEFAULT 0;
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_default_sprite" channels_overlays_dudes_user_settings_dude_sprite NOT NULL DEFAULT 'random';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
