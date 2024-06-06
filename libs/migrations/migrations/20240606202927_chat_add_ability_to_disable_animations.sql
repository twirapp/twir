-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE channels_overlays_chat_animation AS ENUM (
		'DISABLED'
);

ALTER TABLE channels_overlays_chat ADD COLUMN animation channels_overlays_chat_animation DEFAULT 'DISABLED';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
