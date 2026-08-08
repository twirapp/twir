-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers ADD COLUMN name TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers DROP COLUMN name;
-- +goose StatementEnd
