-- +goose Up
DROP EXTENSION IF EXISTS ulid CASCADE;

-- Convert channels_categories_aliases
ALTER TABLE channels_categories_aliases
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Convert channels_scheduled_vips
ALTER TABLE channels_scheduled_vips
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Convert channels_chat_wall_settings
ALTER TABLE channels_chat_wall_settings
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Convert channels_chat_wall
ALTER TABLE channels_chat_wall
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Convert channels_chat_wall_log
-- First, drop the foreign key constraint
ALTER TABLE channels_chat_wall_log DROP CONSTRAINT IF EXISTS channels_chat_wall_log_wall_id_fkey;

-- Convert the columns
ALTER TABLE channels_chat_wall_log
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7(),
    ALTER COLUMN wall_id TYPE UUID USING uuidv7();

-- Recreate the foreign key constraint
ALTER TABLE channels_chat_wall_log
    ADD CONSTRAINT channels_chat_wall_log_wall_id_fkey
    FOREIGN KEY (wall_id) REFERENCES channels_chat_wall(id) ON DELETE CASCADE;

-- Convert channels_chat_translation_settings
ALTER TABLE channels_chat_translation_settings
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Convert channels_giveaways
ALTER TABLE channels_giveaways
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Convert channels_giveaways_participants
-- First, drop the foreign key constraint
ALTER TABLE channels_giveaways_participants DROP CONSTRAINT IF EXISTS channels_giveaways_participants_giveaway_id_fkey;

-- Convert the columns
ALTER TABLE channels_giveaways_participants
    ALTER COLUMN id TYPE UUID USING uuidv7(),
    ALTER COLUMN id SET DEFAULT uuidv7(),
    ALTER COLUMN giveaway_id TYPE UUID USING uuidv7();

-- Recreate the foreign key constraint
ALTER TABLE channels_giveaways_participants
    ADD CONSTRAINT channels_giveaways_participants_giveaway_id_fkey
    FOREIGN KEY (giveaway_id) REFERENCES channels_giveaways(id) ON DELETE CASCADE;

-- Convert channels_integrations_donationalerts (public_id only)
ALTER TABLE channels_integrations_donationalerts
    ALTER COLUMN public_id TYPE UUID USING uuidv7(),
    ALTER COLUMN public_id SET DEFAULT uuidv7();

-- Convert channels_integrations_donatepay (if exists)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_name = 'channels_integrations_donatepay'
    ) THEN
        ALTER TABLE channels_integrations_donatepay
            ALTER COLUMN id TYPE UUID USING uuidv7(),
            ALTER COLUMN id SET DEFAULT uuidv7();
    END IF;
END $$;

-- +goose Down
-- This migration is not reversible as we're changing ID types
-- Rolling back would require restoring from backup
SELECT 'Migration ulid_to_uuidv7 is not reversible';
