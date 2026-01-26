-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_timers ADD COLUMN "offline_enabled" boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_timers DROP COLUMN "offline_enabled";
-- +goose StatementEnd
