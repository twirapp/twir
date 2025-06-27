-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_OVERLAYS';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_OVERLAYS';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DELETE
FROM pg_enum
WHERE enumlabel = 'channels_roles_permissions_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'VIEW_OVERLAYS');

DELETE
FROM pg_enum
WHERE enumlabel = 'channels_roles_permissions_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'MANAGE_OVERLAYS');
-- +goose StatementEnd
