-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_GAMES';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_GAMES';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
