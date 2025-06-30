-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE channel_redemptions_history;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
