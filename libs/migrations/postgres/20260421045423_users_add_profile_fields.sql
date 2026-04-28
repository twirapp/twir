-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN login TEXT DEFAULT '',
    ADD COLUMN display_name TEXT DEFAULT '',
    ADD COLUMN avatar TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN login,
    DROP COLUMN display_name,
    DROP COLUMN avatar;
-- +goose StatementEnd
