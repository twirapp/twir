-- +goose Up
-- +goose StatementBegin
ALTER TABLE "channels_overlays"
ADD COLUMN "insta_save" boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "channels_overlays" DROP COLUMN "insta_save";
-- +goose StatementEnd
