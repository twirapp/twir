-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_overlays_now_playing" ADD COLUMN "show_image" bool default true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
