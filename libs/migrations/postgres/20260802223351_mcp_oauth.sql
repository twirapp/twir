-- +goose Up
-- +goose StatementBegin
CREATE TABLE oauth_clients (
	id                        UUID        PRIMARY KEY DEFAULT uuidv7(),
	client_id                 TEXT        NOT NULL,
	metadata                  JSONB       NOT NULL DEFAULT '{}'::jsonb,
	redirect_uris             TEXT[]      NOT NULL DEFAULT '{}',
	grant_types               TEXT[]      NOT NULL DEFAULT '{}',
	response_types            TEXT[]      NOT NULL DEFAULT '{}',
	token_endpoint_auth_method TEXT       NOT NULL,
	scopes                    TEXT[]      NOT NULL DEFAULT '{}',
	created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT oauth_clients_client_id_key UNIQUE (client_id)
);

CREATE TABLE oauth_authorization_codes (
	id             UUID        PRIMARY KEY DEFAULT uuidv7(),
	code_hash      BYTEA       NOT NULL,
	client_id      TEXT        NOT NULL REFERENCES oauth_clients(client_id) ON DELETE CASCADE ON UPDATE CASCADE,
	channel_id     UUID        NOT NULL REFERENCES channels(id) ON DELETE CASCADE ON UPDATE CASCADE,
	user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
	redirect_uri   TEXT        NOT NULL,
	pkce_challenge TEXT        NOT NULL,
	scopes         TEXT[]      NOT NULL DEFAULT '{}',
	resource       TEXT        NOT NULL,
	expires_at     TIMESTAMPTZ NOT NULL,
	created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT oauth_authorization_codes_code_hash_key UNIQUE (code_hash),
	CONSTRAINT oauth_authorization_codes_code_hash_length CHECK (octet_length(code_hash) = 32)
);

CREATE TABLE oauth_tokens (
	id                 UUID        PRIMARY KEY DEFAULT uuidv7(),
	family_id          UUID        NOT NULL DEFAULT uuidv7(),
	client_id          TEXT        NOT NULL REFERENCES oauth_clients(client_id) ON DELETE CASCADE ON UPDATE CASCADE,
	channel_id         UUID        NOT NULL REFERENCES channels(id) ON DELETE CASCADE ON UPDATE CASCADE,
	user_id            UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
	access_token_hash  BYTEA       NOT NULL,
	refresh_token_hash BYTEA       NOT NULL,
	scopes             TEXT[]      NOT NULL DEFAULT '{}',
	resource           TEXT        NOT NULL,
	access_expires_at  TIMESTAMPTZ NOT NULL,
	refresh_expires_at TIMESTAMPTZ NOT NULL,
	revoked_at         TIMESTAMPTZ,
	replaced_by_id     UUID        REFERENCES oauth_tokens(id) ON DELETE SET NULL ON UPDATE CASCADE,
	created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT oauth_tokens_access_token_hash_key UNIQUE (access_token_hash),
	CONSTRAINT oauth_tokens_refresh_token_hash_key UNIQUE (refresh_token_hash),
	CONSTRAINT oauth_tokens_access_token_hash_length CHECK (octet_length(access_token_hash) = 32),
	CONSTRAINT oauth_tokens_refresh_token_hash_length CHECK (octet_length(refresh_token_hash) = 32)
);

CREATE INDEX oauth_authorization_codes_client_id_idx ON oauth_authorization_codes (client_id);
CREATE INDEX oauth_authorization_codes_channel_id_idx ON oauth_authorization_codes (channel_id);
CREATE INDEX oauth_authorization_codes_user_id_idx ON oauth_authorization_codes (user_id);
CREATE INDEX oauth_tokens_family_id_idx ON oauth_tokens (family_id);
CREATE INDEX oauth_tokens_client_id_idx ON oauth_tokens (client_id);
CREATE INDEX oauth_tokens_channel_id_idx ON oauth_tokens (channel_id);
CREATE INDEX oauth_tokens_user_id_idx ON oauth_tokens (user_id);
CREATE INDEX oauth_tokens_replaced_by_id_idx ON oauth_tokens (replaced_by_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth_tokens;
DROP TABLE IF EXISTS oauth_authorization_codes;
DROP TABLE IF EXISTS oauth_clients;
-- +goose StatementEnd
