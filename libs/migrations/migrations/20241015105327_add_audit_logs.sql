-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "audit_logs" (
	"id" UUID PRIMARY KEY default uuid_generate_v4(),
	"table_name" VARCHAR(255) NOT NULL,
	"operation_type" VARCHAR(255) NOT NULL,
	"old_value" TEXT,
	"new_value" TEXT,
	"object_id" TEXT,
	"user_id" TEXT,
	"channel_id" TEXT,
	"created_at" TIMESTAMP default now() NOT NULL,
	FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION,
	FOREIGN KEY ("channel_id") REFERENCES "users" ("id") ON DELETE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
