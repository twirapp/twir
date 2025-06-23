-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_alerts"
	ADD COLUMN "audio_volume" int2 not null default 100;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "channels_alerts"
	DROP COLUMN "audio_volume"
-- +goose StatementEnd
