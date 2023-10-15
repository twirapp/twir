-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE integrations_service_enum ADD VALUE 'DISCORD';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
