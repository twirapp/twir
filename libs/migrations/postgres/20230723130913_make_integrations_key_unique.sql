-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE integrations
	ADD CONSTRAINT integrations_unique_service UNIQUE (service);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
