-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_overlays_dudes_user_settings_dude_sprite ADD VALUE 'random';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
