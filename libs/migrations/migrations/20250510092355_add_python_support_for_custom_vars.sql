-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_customvars ADD COLUMN "script_language" VARCHAR(255) NOT NULL DEFAULT 'javascript';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
