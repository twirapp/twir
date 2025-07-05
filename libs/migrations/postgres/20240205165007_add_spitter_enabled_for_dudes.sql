-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_overlays_dudes" ADD COLUMN "spitter_emote_enabled" boolean NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
