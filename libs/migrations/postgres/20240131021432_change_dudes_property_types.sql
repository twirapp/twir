-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_overlays_dudes" ALTER COLUMN "dude_scale" TYPE real;
ALTER TABLE "channels_overlays_dudes" ALTER COLUMN "name_box_font_weight" TYPE integer USING ("name_box_font_weight"::integer);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
