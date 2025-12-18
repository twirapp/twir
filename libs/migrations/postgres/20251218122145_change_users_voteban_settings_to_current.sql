-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

UPDATE channels_games_voteban SET vote_duration = 10 WHERE vote_duration < 10;
UPDATE channels_games_voteban SET needed_votes = 2 WHERE needed_votes < 2;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
