-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_dashboard_widgets" (
	"id" UUID PRIMARY KEY default uuidv7(),
	"channel_id" TEXT NOT NULL,
	"widget_id" TEXT NOT NULL,
	"x" INT NOT NULL,
	"y" INT NOT NULL,
	"w" INT NOT NULL,
	"h" INT NOT NULL,
	"min_w" INT NOT NULL,
	"min_h" INT NOT NULL,
	"visible" BOOLEAN NOT NULL DEFAULT true,
	"created_at" TIMESTAMP default now() NOT NULL,
	"updated_at" TIMESTAMP default now() NOT NULL,
	FOREIGN KEY ("channel_id") REFERENCES "users" ("id") ON DELETE CASCADE,
	UNIQUE("channel_id", "widget_id")
);

CREATE INDEX "channels_dashboard_widgets_channel_id_idx" ON "channels_dashboard_widgets" ("channel_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS "channels_dashboard_widgets";

-- +goose StatementEnd
