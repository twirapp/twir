-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "users_stats" ADD COLUMN "is_mod" boolean NOT NULL DEFAULT false;
ALTER TABLE "users_stats" ADD COLUMN "is_vip" boolean NOT NULL DEFAULT false;
ALTER TABLE "users_stats" ADD COLUMN "is_subscriber" boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "users_stats" DROP COLUMN "is_mod";
ALTER TABLE "users_stats" DROP COLUMN "is_vip";
ALTER TABLE "users_stats" DROP COLUMN "is_subscriber";
-- +goose StatementEnd
