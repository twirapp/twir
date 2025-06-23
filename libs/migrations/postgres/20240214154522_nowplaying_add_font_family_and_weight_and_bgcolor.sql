-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_overlays_now_playing" ADD COLUMN "font_family" VARCHAR(255) NOT NULL DEFAULT 'inter';
ALTER TABLE "channels_overlays_now_playing" ADD COLUMN "font_weight" int NOT NULL DEFAULT 400;
ALTER TABLE "channels_overlays_now_playing" ADD COLUMN "background_color" VARCHAR(255) NOT NULL DEFAULT 'rgba(0, 0, 0, 0)';

ALTER TYPE "channel_overlay_now_playing_preset" ADD VALUE 'SIMPLE_LINE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
