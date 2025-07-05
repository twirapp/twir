-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE notifications
	ALTER COLUMN message DROP NOT NULL;
ALTER TABLE notifications
	ADD COLUMN editor_js_json TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
