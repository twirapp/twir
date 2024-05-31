-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE UNIQUE INDEX "eventsub_subscriptions_topic_id_user_id_idx" ON "eventsub_subscriptions" ("topic_id", "user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
