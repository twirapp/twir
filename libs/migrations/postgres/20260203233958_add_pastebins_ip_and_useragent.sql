-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE pastebins
ADD COLUMN ip inet NULL,
ADD COLUMN user_agent TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
