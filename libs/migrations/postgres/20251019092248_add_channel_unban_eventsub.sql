-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO eventsub_topics (topic, version)
VALUES ('channel.unban', '1')
ON CONFLICT (topic) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
