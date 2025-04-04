-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE "channels_roles_permissions_enum" ADD VALUE 'VIEW_MODULES';
ALTER TYPE "channels_roles_permissions_enum" ADD VALUE 'MANAGE_MODULES';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
