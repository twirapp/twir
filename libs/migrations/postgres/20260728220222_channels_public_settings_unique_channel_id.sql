-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX channels_public_settings_channel_id_unique_idx ON channels_public_settings (channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channels_public_settings_channel_id_unique_idx;
-- +goose StatementEnd
