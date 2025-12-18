-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

UPDATE channels_games_voteban
SET vote_duration = CASE
											WHEN vote_duration = 1 THEN 10
											ELSE vote_duration
										END,
needed_votes  = CASE
 							  	WHEN needed_votes = 1 THEN 2
									ELSE needed_votes
								END
WHERE vote_duration = 1
	 OR needed_votes = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
