-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_overlays_layers" ADD COLUMN periodically_refetch_data boolean NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
