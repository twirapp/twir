-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_EVENTS';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_EVENTS';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DELETE
FROM pg_enum
WHERE enumlabel = 'channels_roles_permissions_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'VIEW_EVENTS');

DELETE
FROM pg_enum
WHERE enumlabel = 'channels_roles_permissions_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'MANAGE_EVENTS');
-- +goose StatementEnd
