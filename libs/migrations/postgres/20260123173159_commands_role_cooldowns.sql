-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels_commands_role_cooldowns (
	id uuid PRIMARY KEY DEFAULT uuidv7(),
	command_id uuid NOT NULL REFERENCES channels_commands(id) ON DELETE CASCADE,
	role_id uuid NOT NULL REFERENCES channels_roles(id) ON DELETE CASCADE,
	cooldown INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	UNIQUE(command_id, role_id)
);

CREATE INDEX idx_commands_role_cooldowns_command_id ON channels_commands_role_cooldowns(command_id);
CREATE INDEX idx_commands_role_cooldowns_role_id ON channels_commands_role_cooldowns(role_id);

INSERT INTO channels_commands_role_cooldowns (command_id, role_id, cooldown)
SELECT
    c.id AS command_id,
    cr.id AS role_id,
    c.cooldown
FROM channels_commands c
CROSS JOIN LATERAL unnest(c."cooldown_roles_ids") AS role_id_str
INNER JOIN channels_roles cr ON cr.id = role_id_str::uuid
WHERE c."cooldown" is not null and c."cooldown" > 0 AND c."cooldown_roles_ids" IS NOT NULL
  AND array_length(c."cooldown_roles_ids", 1) > 0
ON CONFLICT (command_id, role_id) DO NOTHING;

ALTER TABLE channels_commands DROP COLUMN IF EXISTS cooldown_roles_ids;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS channels_commands_role_cooldowns;
-- +goose StatementEnd
