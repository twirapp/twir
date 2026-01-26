-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_timers ADD COLUMN "online_enabled" boolean NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_timers DROP COLUMN "online_enabled";
-- +goose StatementEnd
