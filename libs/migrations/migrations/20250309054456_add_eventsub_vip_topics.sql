-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO eventsub_topics (topic, version, condition_type)
VALUES ('channel.vip.add', '1', 'BROADCASTER_USER_ID'),
       ('channel.vip.remove', '1', 'BROADCASTER_USER_ID');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
