-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_commands_responses ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_commands_responses DROP COLUMN IF EXISTS platforms;
-- +goose StatementEnd
