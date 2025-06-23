-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_overlays_dudes" (
	"id" uuid NOT NULL PRIMARY KEY,
	"channel_id" text NOT NULL REFERENCES "channels" ("id") ON DELETE CASCADE,
	"dude_color" text NOT NULL,
	"dude_max_life_time" integer NOT NULL,
	"dude_gravity" integer NOT NULL,
	"dude_scale" integer NOT NULL,
	"dude_sounds_enabled" boolean NOT NULL,
	"dude_sounds_volume" real NOT NULL,
	"message_box_border_radius" integer NOT NULL,
	"message_box_box_color" text NOT NULL,
	"message_box_font_family" text NOT NULL,
	"message_box_font_size" integer NOT NULL,
	"message_box_padding" integer NOT NULL,
	"message_box_show_time" integer NOT NULL,
	"message_box_fill" text NOT NULL,
	"name_box_font_family" text NOT NULL,
	"name_box_font_size" integer NOT NULL,
	"name_box_fill" text[] NOT NULL,
	"name_box_line_join" text NOT NULL,
	"name_box_stroke_thickness" integer NOT NULL,
	"name_box_stroke" text NOT NULL,
	"name_box_fill_gradient_stops" float4[] NOT NULL,
	"name_box_fill_gradient_type" integer NOT NULL,
	"name_box_font_style" text NOT NULL,
	"name_box_font_variant" text NOT NULL,
	"name_box_font_weight" text NOT NULL,
	"name_box_drop_shadow" boolean NOT NULL,
	"name_box_drop_shadow_alpha" real NOT NULL,
	"name_box_drop_shadow_angle" real NOT NULL,
	"name_box_drop_shadow_blur" real NOT NULL,
	"name_box_drop_shadow_distance" real NOT NULL,
	"name_box_drop_shadow_color" text NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
