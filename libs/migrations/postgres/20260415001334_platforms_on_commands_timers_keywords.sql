-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_commands
    ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';

ALTER TABLE channels_timers
    ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';

ALTER TABLE channels_keywords
    ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_commands  DROP COLUMN IF EXISTS platforms;
ALTER TABLE channels_timers    DROP COLUMN IF EXISTS platforms;
ALTER TABLE channels_keywords  DROP COLUMN IF EXISTS platforms;
-- +goose StatementEnd
