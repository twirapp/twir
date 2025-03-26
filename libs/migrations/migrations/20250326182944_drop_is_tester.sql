-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "users" DROP COLUMN "isTester";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
