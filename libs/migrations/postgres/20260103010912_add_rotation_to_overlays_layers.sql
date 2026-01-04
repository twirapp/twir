-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers ADD COLUMN IF NOT EXISTS rotation INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers DROP COLUMN IF EXISTS rotation;
-- +goose StatementEnd
