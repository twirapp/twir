-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'channels_commands' AND column_name = 'platform'
    ) THEN
        ALTER TABLE channels_commands DROP COLUMN platform;
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'channels_commands' AND column_name = 'platforms'
    ) THEN
        ALTER TABLE channels_commands ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'channels_timers' AND column_name = 'platform'
    ) THEN
        ALTER TABLE channels_timers DROP COLUMN platform;
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'channels_timers' AND column_name = 'platforms'
    ) THEN
        ALTER TABLE channels_timers ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'channels_keywords' AND column_name = 'platform'
    ) THEN
        ALTER TABLE channels_keywords DROP COLUMN platform;
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'channels_keywords' AND column_name = 'platforms'
    ) THEN
        ALTER TABLE channels_keywords ADD COLUMN platforms platform[] NOT NULL DEFAULT '{}';
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'no-op';
-- +goose StatementEnd
