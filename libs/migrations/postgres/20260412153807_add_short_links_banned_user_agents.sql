-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_links_banned_user_agents (
	id TEXT PRIMARY KEY DEFAULT uuidv7(),
	user_id TEXT NOT NULL,
	pattern TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	UNIQUE (user_id, pattern)
);

CREATE INDEX short_links_banned_user_agents_user_id_idx ON short_links_banned_user_agents(user_id);

COMMENT ON TABLE short_links_banned_user_agents IS 'Regex patterns to ban specific user agents from redirecting short links';
COMMENT ON COLUMN short_links_banned_user_agents.pattern IS 'Go-compatible regex pattern to match against the User-Agent header';
COMMENT ON COLUMN short_links_banned_user_agents.description IS 'Optional human-readable description of what this pattern blocks';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS short_links_banned_user_agents;
-- +goose StatementEnd
