-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_moderation_settings ADD COLUMN name varchar(255);
ALTER TABLE channels_moderation_settings ADD COLUMN deny_list_regexp_enabled bool NOT NULL DEFAULT FALSE;
ALTER TABLE channels_moderation_settings ADD COLUMN deny_list_word_boundary_enabled bool NOT NULL DEFAULT FALSE;
ALTER TABLE channels_moderation_settings ADD COLUMN deny_list_sensitivity_enabled bool NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
