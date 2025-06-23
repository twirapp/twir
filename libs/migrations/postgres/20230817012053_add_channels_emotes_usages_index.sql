-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE INDEX "channels_emotes_usages_channelId_userId" ON "channels_emotes_usages" ("channelId", "userId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
