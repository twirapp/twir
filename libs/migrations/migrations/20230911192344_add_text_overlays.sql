-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "channels_overlays"
(
	"id"         uuid         NOT NULL primary key DEFAULT uuid_generate_v4(),
	"channel_id" text         NOT NULL REFERENCES "channels" ("id") ON DELETE CASCADE,
	"name"       varchar(255) NOT NULL,
	"created_at" timestamp    NOT NULL             DEFAULT now(),
	"updated_at" timestamp    NOT NULL             DEFAULT now(),
	"height" 	 int          NOT NULL             DEFAULT 1080,
	"width" 	 int          NOT NULL             DEFAULT 1920
);

CREATE TYPE "channels_overlays_layers_type" AS ENUM ('HTML');

CREATE TABLE "channels_overlays_layers"
(
	"id"         uuid                          NOT NULL primary key DEFAULT uuid_generate_v4(),
	"overlay_id" uuid                          NOT NULL REFERENCES "channels_overlays" ("id") ON DELETE CASCADE,
	"type"       channels_overlays_layers_type NOT NULL,
	"settings"   jsonb                         NOT NULL             DEFAULT '{}'::jsonb,
	"pos_x"      int                           NOT NULL             DEFAULT 0,
	"pos_y"      int                           NOT NULL             DEFAULT 0,
	"width"      int                           NOT NULL             DEFAULT 0,
	"height"     int                           NOT NULL             DEFAULT 0,
	"created_at" timestamp                     NOT NULL             DEFAULT now(),
	"updated_at" timestamp                     NOT NULL             DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
