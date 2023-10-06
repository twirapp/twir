-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE "public"."integrations_service_enum" ADD VALUE 'DISCORD';
INSERT INTO "integrations" ("service") VALUES ('DISCORD');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
