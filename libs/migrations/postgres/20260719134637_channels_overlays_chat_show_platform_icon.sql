-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_overlays_chat ADD COLUMN show_platform_icon boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_overlays_chat DROP COLUMN show_platform_icon;
-- +goose StatementEnd
