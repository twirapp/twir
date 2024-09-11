-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX users_api_key_idx ON users("apiKey");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
