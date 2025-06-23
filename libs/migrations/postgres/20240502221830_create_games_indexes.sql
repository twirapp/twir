-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE UNIQUE INDEX channels_games_8ball_channel_id_idx ON channels_games_8ball (channel_id);
CREATE UNIQUE INDEX channels_games_duel_channel_id_idx ON channels_games_duel (channel_id);
CREATE UNIQUE INDEX channels_games_russian_roulette_channel_id_idx ON channels_games_russian_roulette (channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
