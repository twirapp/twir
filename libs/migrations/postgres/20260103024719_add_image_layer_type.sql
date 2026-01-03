-- +goose Up
-- +goose StatementBegin
ALTER TYPE channels_overlays_layers_type ADD VALUE IF NOT EXISTS 'IMAGE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Cannot remove enum values in PostgreSQL, would need to recreate the type
-- +goose StatementEnd
