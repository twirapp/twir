-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE channels_giveaways (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	channel_id TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	started_at TIMESTAMPTZ,
	ended_at TIMESTAMPTZ,
	archived_at TIMESTAMPTZ,
	stopped_at TIMESTAMPTZ,
	keyword TEXT NOT NULL,
	created_by_user_id TEXT,

	FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
	FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE
);

CREATE TABLE channels_giveaways_participants (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	giveaway_id ulid NOT NULL,
	is_winner BOOLEAN NOT NULL DEFAULT false,
	display_name TEXT NOT NULL,
	user_id TEXT NOT NULL,

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
	FOREIGN KEY(giveaway_id) REFERENCES channels_giveaways(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX channels_giveaways_participants_unique ON channels_giveaways_participants (giveaway_id, user_id);

ALTER TYPE channels_roles_permissions_enum ADD VALUE 'VIEW_GIVEAWAYS';
ALTER TYPE channels_roles_permissions_enum ADD VALUE 'MANAGE_GIVEAWAYS';


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE channels_giveaways_participants;
DROP TABLE channels_giveaways;
-- +goose StatementEnd
