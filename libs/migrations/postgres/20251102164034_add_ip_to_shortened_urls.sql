-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE shortened_urls
	ADD COLUMN user_agent TEXT;
ALTER TABLE shortened_urls
	ADD COLUMN user_ip cidr;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
