-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Create enum type for giveaway types
CREATE TYPE giveaway_type AS ENUM ('KEYWORD', 'ONLINE_CHATTERS');

-- Add new columns to channels_giveaways table
ALTER TABLE channels_giveaways
	ADD COLUMN type giveaway_type NOT NULL DEFAULT 'KEYWORD',
	ADD COLUMN min_watched_time BIGINT,
	ADD COLUMN min_messages INTEGER,
	ADD COLUMN min_used_channel_points BIGINT,
	ADD COLUMN min_follow_duration BIGINT,
	ADD COLUMN require_subscription BOOLEAN NOT NULL DEFAULT false;

-- Make keyword nullable since ONLINE_CHATTERS type doesn't need it
ALTER TABLE channels_giveaways
	ALTER COLUMN keyword DROP NOT NULL;

-- Add constraint to ensure keyword is present for KEYWORD type
ALTER TABLE channels_giveaways
	ADD CONSTRAINT keyword_required_for_keyword_type
	CHECK (type != 'KEYWORD' OR keyword IS NOT NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- Remove constraint
ALTER TABLE channels_giveaways
	DROP CONSTRAINT IF EXISTS keyword_required_for_keyword_type;

-- Restore keyword NOT NULL
ALTER TABLE channels_giveaways
	ALTER COLUMN keyword SET NOT NULL;

-- Remove new columns
ALTER TABLE channels_giveaways
	DROP COLUMN IF EXISTS require_subscription,
	DROP COLUMN IF EXISTS min_follow_duration,
	DROP COLUMN IF EXISTS min_used_channel_points,
	DROP COLUMN IF EXISTS min_messages,
	DROP COLUMN IF EXISTS min_watched_time,
	DROP COLUMN IF EXISTS type;

-- Drop enum type
DROP TYPE IF EXISTS giveaway_type;

-- +goose StatementEnd
