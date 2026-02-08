-- +goose Up
-- +goose StatementBegin
ALTER TABLE discord_sended_notifications
ADD COLUMN last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT now();

COMMENT ON COLUMN discord_sended_notifications.last_updated_at IS 'Last time the message was updated with current stream info';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE discord_sended_notifications
DROP COLUMN IF EXISTS last_updated_at;
-- +goose StatementEnd
