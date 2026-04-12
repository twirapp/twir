-- +goose Up
-- +goose StatementBegin
-- Drop old global banned user agents table (we'll recreate as presets)
DROP TABLE IF EXISTS short_links_global_banned_user_agents CASCADE;

-- Create presets table
CREATE TABLE short_links_banned_ua_presets (
    id TEXT PRIMARY KEY DEFAULT uuidv7(),
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX short_links_banned_ua_presets_user_id_idx ON short_links_banned_ua_presets(user_id);

COMMENT ON TABLE short_links_banned_ua_presets IS 'User-defined presets of banned user agent patterns';
COMMENT ON COLUMN short_links_banned_ua_presets.name IS 'Display name for the preset (e.g., "Chatterino Block")';

-- Create preset patterns table (patterns belonging to a preset)
CREATE TABLE short_links_banned_ua_preset_patterns (
    id TEXT PRIMARY KEY DEFAULT uuidv7(),
    preset_id TEXT NOT NULL,
    pattern TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (preset_id) REFERENCES short_links_banned_ua_presets(id) ON DELETE CASCADE,
    UNIQUE (preset_id, pattern)
);

CREATE INDEX short_links_banned_ua_preset_patterns_preset_id_idx ON short_links_banned_ua_preset_patterns(preset_id);

COMMENT ON TABLE short_links_banned_ua_preset_patterns IS 'Individual regex patterns within a preset';
COMMENT ON COLUMN short_links_banned_ua_preset_patterns.pattern IS 'Go-compatible regex pattern';

-- Create link-presets pivot table (many-to-many)
CREATE TABLE short_links_link_presets (
    id TEXT PRIMARY KEY DEFAULT uuidv7(),
    link_id TEXT NOT NULL,
    preset_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (link_id) REFERENCES shortened_urls(id) ON DELETE CASCADE,
    FOREIGN KEY (preset_id) REFERENCES short_links_banned_ua_presets(id) ON DELETE CASCADE,
    UNIQUE (link_id, preset_id)
);

CREATE INDEX short_links_link_presets_link_id_idx ON short_links_link_presets(link_id);
CREATE INDEX short_links_link_presets_preset_id_idx ON short_links_link_presets(preset_id);

COMMENT ON TABLE short_links_link_presets IS 'Pivot table linking shortened URLs to banned UA presets';

-- Keep the direct link patterns table for patterns applied directly to links
-- (already created in previous migration, just ensure it exists)
-- short_links_link_banned_user_agents

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS short_links_link_presets CASCADE;
DROP TABLE IF EXISTS short_links_banned_ua_preset_patterns CASCADE;
DROP TABLE IF EXISTS short_links_banned_ua_presets CASCADE;

-- Recreate the old global table (for rollback only)
CREATE TABLE short_links_global_banned_user_agents (
    id TEXT PRIMARY KEY DEFAULT uuidv7(),
    user_id TEXT NOT NULL,
    pattern TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, pattern)
);

CREATE INDEX short_links_global_banned_user_agents_user_id_idx ON short_links_global_banned_user_agents(user_id);
-- +goose StatementEnd
