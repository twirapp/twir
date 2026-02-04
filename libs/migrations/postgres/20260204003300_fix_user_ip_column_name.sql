-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE pastebins RENAME COLUMN ip TO user_ip;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
