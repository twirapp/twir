-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE "notifications_messages";

ALTER TABLE "notifications" ADD COLUMN "message" text not null;
ALTER TABLE "notifications" DROP COLUMN "imageSrc";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
