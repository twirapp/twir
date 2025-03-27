-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS channels_scheduled_vips (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	channel_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	remove_at TIMESTAMPTZ,

	FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

UPDATE channels_commands SET name = 'vip_twir_add_own_module' WHERE name = 'vips add';
UPDATE channels_commands SET name = 'unvip_twir_add_own_module' WHERE name = 'vips remove';
UPDATE channels_commands SET name = 'vips_list_twir_add_own_module' WHERE name = 'vips list';

ALTER type channels_commands_module_enum ADD VALUE 'VIPS';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE channels_scheduled_vips;
-- +goose StatementEnd
