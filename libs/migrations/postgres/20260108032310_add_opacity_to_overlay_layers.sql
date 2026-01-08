-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers ADD COLUMN opacity REAL NOT NULL DEFAULT 1.0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers DROP COLUMN opacity;
-- +goose StatementEnd
