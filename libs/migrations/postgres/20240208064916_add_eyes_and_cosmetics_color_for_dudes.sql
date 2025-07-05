-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_eyes_color" text NOT NULL DEFAULT '#FFFFFF';
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "dude_cosmetics_color" text NOT NULL DEFAULT '#FFFFFF';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
