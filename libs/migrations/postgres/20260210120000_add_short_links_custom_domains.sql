-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE short_links_custom_domains (
	id TEXT PRIMARY KEY DEFAULT uuidv7(),
	user_id TEXT NOT NULL UNIQUE,
	domain TEXT NOT NULL UNIQUE,
	verified BOOLEAN NOT NULL DEFAULT false,
	verification_token TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON TABLE short_links_custom_domains IS 'Custom domains for users short links';
COMMENT ON COLUMN short_links_custom_domains.verified IS 'Whether the domain has been verified via CNAME check';
COMMENT ON COLUMN short_links_custom_domains.verification_token IS 'Expected CNAME target: short-{token}.twir.app';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS short_links_custom_domains;
-- +goose StatementEnd
