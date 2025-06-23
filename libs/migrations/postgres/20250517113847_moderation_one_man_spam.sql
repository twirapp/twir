-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_moderation_settings_type_enum ADD VALUE 'one_man_spam';

ALTER TABLE channels_moderation_settings ADD COLUMN one_man_spam_minimum_stored_messages integer NOT NULL DEFAULT 0;
ALTER TABLE channels_moderation_settings ADD COLUMN one_man_spam_message_memory_seconds integer NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
