-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
	ADD COLUMN links_shortener_custom_domains INTEGER NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plans
	DROP COLUMN IF EXISTS links_shortener_custom_domains;
-- +goose StatementEnd
