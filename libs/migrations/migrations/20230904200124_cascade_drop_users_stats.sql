-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "users_stats" DROP CONSTRAINT "FK_3d6cc217af2451426c44a30e678";
ALTER TABLE "users_stats" ADD CONSTRAINT "users_stats_user_id" FOREIGN KEY ("userId") REFERENCES "users" ("id") ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
