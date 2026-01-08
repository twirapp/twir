-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_overlays_layers" ADD COLUMN "locked" boolean NOT NULL DEFAULT false;
ALTER TABLE "channels_overlays_layers" ADD COLUMN "visible" boolean NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "channels_overlays_layers" DROP COLUMN "locked";
ALTER TABLE "channels_overlays_layers" DROP COLUMN "visible";
-- +goose StatementEnd
