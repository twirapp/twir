-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_visible_name" boolean NOT NULL DEFAULT true;
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_grow_time" integer NOT NULL DEFAULT 300000;
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_grow_max_scale" integer NOT NULL DEFAULT 10;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
