-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE channels_streams
	ALTER COLUMN "gameId" TYPE varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
