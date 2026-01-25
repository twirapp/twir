-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE shortened_urls DROP CONSTRAINT shortened_urls_pkey;

ALTER TABLE shortened_urls
	ADD COLUMN id TEXT;

UPDATE shortened_urls
SET id = uuidv7()
WHERE id IS NULL;

ALTER TABLE shortened_urls
	ALTER COLUMN id SET NOT NULL,
	ALTER COLUMN id SET DEFAULT uuidv7();

ALTER TABLE shortened_urls
	ADD PRIMARY KEY (id);

ALTER TABLE shortened_urls
	ADD COLUMN domain TEXT;

CREATE INDEX shortened_urls_domain_idx ON shortened_urls(domain);

CREATE UNIQUE INDEX shortened_urls_domain_short_id_unique_idx
	ON shortened_urls (COALESCE(domain, ''), short_id);

COMMENT ON COLUMN shortened_urls.domain IS 'Custom domain for this short link (NULL = default twir.app domain)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS shortened_urls_domain_short_id_unique_idx;
DROP INDEX IF EXISTS shortened_urls_domain_idx;

ALTER TABLE shortened_urls DROP CONSTRAINT IF EXISTS shortened_urls_pkey;
ALTER TABLE shortened_urls DROP COLUMN IF EXISTS domain;
ALTER TABLE shortened_urls DROP COLUMN IF EXISTS id;

ALTER TABLE shortened_urls ADD PRIMARY KEY (short_id);
-- +goose StatementEnd
