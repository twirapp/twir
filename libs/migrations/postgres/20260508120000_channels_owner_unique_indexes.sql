-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM channels
        WHERE twitch_user_id IS NOT NULL
        GROUP BY twitch_user_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot create unique index on channels.twitch_user_id: duplicate owners exist';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM channels
        WHERE kick_user_id IS NOT NULL
        GROUP BY kick_user_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot create unique index on channels.kick_user_id: duplicate owners exist';
    END IF;
END $$;

CREATE UNIQUE INDEX channels_twitch_user_id_unique_idx
    ON channels (twitch_user_id)
    WHERE twitch_user_id IS NOT NULL;

CREATE UNIQUE INDEX channels_kick_user_id_unique_idx
    ON channels (kick_user_id)
    WHERE kick_user_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channels_twitch_user_id_unique_idx;
DROP INDEX IF EXISTS channels_kick_user_id_unique_idx;
-- +goose StatementEnd
