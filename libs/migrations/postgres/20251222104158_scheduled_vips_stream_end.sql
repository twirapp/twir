-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE scheduled_vip_remove_type AS ENUM ('time', 'stream_end');

ALTER TABLE "channels_scheduled_vips"
	ADD COLUMN "remove_type" scheduled_vip_remove_type NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
