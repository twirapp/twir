-- +goose Up
-- +goose StatementBegin
ALTER TYPE channels_overlays_layers_type ADD VALUE IF NOT EXISTS 'TEXT';
ALTER TYPE channels_overlays_layers_type ADD VALUE IF NOT EXISTS 'VIDEO';
ALTER TYPE channels_overlays_layers_type ADD VALUE IF NOT EXISTS 'IFRAME';
ALTER TYPE channels_overlays_layers_type ADD VALUE IF NOT EXISTS 'YOUTUBE';
ALTER TYPE channels_overlays_layers_type ADD VALUE IF NOT EXISTS 'EMOTE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Cannot remove enum values in PostgreSQL, would need to recreate the type
-- +goose StatementEnd
