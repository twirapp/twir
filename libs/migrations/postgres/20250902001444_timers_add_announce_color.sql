-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE channels_timers_responses ADD COLUMN announce_color INT4 NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
