-- +goose Up
-- +goose StatementBegin
-- Rename existing table to indicate it's for global (per-user) bans
ALTER TABLE short_links_banned_user_agents RENAME TO short_links_global_banned_user_agents;

-- Update index names
ALTER INDEX short_links_banned_user_agents_user_id_idx 
    RENAME TO short_links_global_banned_user_agents_user_id_idx;

-- Update comments
COMMENT ON TABLE short_links_global_banned_user_agents IS 'Global regex patterns to ban specific user agents from all short links of a user';

-- Create table for per-link specific banned user agents
CREATE TABLE short_links_link_banned_user_agents (
    id TEXT PRIMARY KEY DEFAULT uuidv7(),
    link_id TEXT NOT NULL,
    pattern TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (link_id) REFERENCES shortened_urls(id) ON DELETE CASCADE,
    UNIQUE (link_id, pattern)
);

CREATE INDEX short_links_link_banned_user_agents_link_id_idx ON short_links_link_banned_user_agents(link_id);

COMMENT ON TABLE short_links_link_banned_user_agents IS 'Regex patterns to ban specific user agents from a specific short link';
COMMENT ON COLUMN short_links_link_banned_user_agents.pattern IS 'Go-compatible regex pattern to match against the User-Agent header';
COMMENT ON COLUMN short_links_link_banned_user_agents.description IS 'Optional human-readable description of what this pattern blocks';

-- Add flag to shortened_urls to ignore global bans
ALTER TABLE shortened_urls ADD COLUMN ignore_global_bans BOOLEAN NOT NULL DEFAULT false;

COMMENT ON COLUMN shortened_urls.ignore_global_bans IS 'If true, global banned user agent patterns are not applied to this link';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shortened_urls DROP COLUMN IF EXISTS ignore_global_bans;

DROP TABLE IF EXISTS short_links_link_banned_user_agents;

ALTER INDEX short_links_global_banned_user_agents_user_id_idx 
    RENAME TO short_links_banned_user_agents_user_id_idx;

ALTER TABLE short_links_global_banned_user_agents RENAME TO short_links_banned_user_agents;

COMMENT ON TABLE short_links_banned_user_agents IS 'Regex patterns to ban specific user agents from redirecting short links';
-- +goose StatementEnd
