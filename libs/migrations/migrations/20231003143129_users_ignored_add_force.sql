-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "users_ignored" ADD COLUMN "force" boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "users_ignored" DROP COLUMN "force";
-- +goose StatementEnd
