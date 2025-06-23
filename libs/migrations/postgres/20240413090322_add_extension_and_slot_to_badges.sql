-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DELETE FROM "badges";

ALTER TABLE "badges" ADD COLUMN "ffz_slot" INTEGER NOT NULL;
ALTER TABLE "badges" ADD COLUMN "file_name" varchar(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
