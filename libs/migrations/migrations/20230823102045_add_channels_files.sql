-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "channels_files"
(
	id         uuid default uuid_generate_v4() not null primary key,
	channel_id text                            not null
		constraint "channel_files_channel_id" references channels ("id")
			on update cascade
			on delete cascade,
	mime_type  text                            not null,
	file_name  text                            not null,
	size       int8                            not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "channels_files"
-- +goose StatementEnd
