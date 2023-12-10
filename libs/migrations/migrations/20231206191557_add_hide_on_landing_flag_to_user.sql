-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "users" ADD COLUMN "hide_on_landing_page" boolean default false NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
