-- +goose Up
-- +goose StatementBegin
CREATE TYPE platform AS ENUM ('twitch', 'kick');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE platform;
-- +goose StatementEnd
