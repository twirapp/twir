-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "badges" (
	"id" uuid default uuid_generate_v4() primary key,
	"name" varchar(255) not null,
	"enabled" boolean not null default true,
	"created_at" timestamp not null default now()
);

CREATE TABLE "badges_users" (
	"id" uuid default uuid_generate_v4() primary key,
	"badge_id" uuid not null references badges("id"),
	"user_id" text not null references users("id"),
	"created_at" timestamp not null default now(),
	constraint "badges_users_unique_badge_user" unique (badge_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
