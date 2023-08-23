-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE channels_events_operations_type_enum ADD VALUE 'TRIGGER_ALERT';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE
FROM pg_enum
WHERE enumlabel = 'channels_events_operations_type_enum'
	AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'TRIGGER_ALERT');
-- +goose StatementEnd
