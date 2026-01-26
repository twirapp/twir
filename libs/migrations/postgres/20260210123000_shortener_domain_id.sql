-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE shortened_urls
	ADD COLUMN domain_id UUID;

UPDATE shortened_urls
SET domain_id = domains.id
FROM short_links_custom_domains AS domains
WHERE shortened_urls.domain IS NOT NULL
	AND shortened_urls.domain = domains.domain;

DROP INDEX IF EXISTS shortened_urls_domain_short_id_unique_idx;
DROP INDEX IF EXISTS shortened_urls_domain_idx;

ALTER TABLE shortened_urls
	DROP COLUMN IF EXISTS domain;

ALTER TABLE shortened_urls
	ADD CONSTRAINT shortened_urls_domain_id_fkey
		FOREIGN KEY (domain_id) REFERENCES short_links_custom_domains(id) ON DELETE SET NULL;

CREATE INDEX shortened_urls_domain_id_idx ON shortened_urls(domain_id);

CREATE UNIQUE INDEX shortened_urls_domain_id_short_id_unique_idx
	ON shortened_urls (COALESCE(domain_id, ''), short_id);

COMMENT ON COLUMN shortened_urls.domain_id IS 'Custom domain id for this short link (NULL = default twir.app domain)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE shortened_urls
	ADD COLUMN domain TEXT;

UPDATE shortened_urls
SET domain = domains.domain
FROM short_links_custom_domains AS domains
WHERE shortened_urls.domain_id = domains.id;

DROP INDEX IF EXISTS shortened_urls_domain_id_short_id_unique_idx;
DROP INDEX IF EXISTS shortened_urls_domain_id_idx;

ALTER TABLE shortened_urls
	DROP CONSTRAINT IF EXISTS shortened_urls_domain_id_fkey;

ALTER TABLE shortened_urls
	DROP COLUMN IF EXISTS domain_id;

CREATE INDEX shortened_urls_domain_idx ON shortened_urls(domain);

CREATE UNIQUE INDEX shortened_urls_domain_short_id_unique_idx
	ON shortened_urls (COALESCE(domain, ''), short_id);

COMMENT ON COLUMN shortened_urls.domain IS 'Custom domain for this short link (NULL = default twir.app domain)';
-- +goose StatementEnd
