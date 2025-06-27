-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE shortened_urls ADD COLUMN views int NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
