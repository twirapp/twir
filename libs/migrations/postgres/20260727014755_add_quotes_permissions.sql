-- +goose Up
-- +goose StatementBegin
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_QUOTES';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_QUOTES';

ALTER TYPE channels_commands_module_enum ADD VALUE 'QUOTES';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
