-- +goose Up
-- +goose StatementBegin

-- This migration changes DEFAULT values from uuid_generate_v4() to uuidv7()
-- Existing IDs are NOT changed, only the default for new records

-- tokens table
ALTER TABLE tokens
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_customvars table
ALTER TABLE channels_customvars
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_greetings table
ALTER TABLE channels_greetings
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_integrations table
ALTER TABLE channels_integrations
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_keywords table
ALTER TABLE channels_keywords
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_timers table
ALTER TABLE channels_timers
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_timers_responses table
ALTER TABLE channels_timers_responses
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_events table
ALTER TABLE channels_events
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_events_operations table
ALTER TABLE channels_events_operations
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_events_operations_filters table
ALTER TABLE channels_events_operations_filters
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_streams table
ALTER TABLE channels_streams
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_commands table
ALTER TABLE channels_commands
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_roles table
ALTER TABLE channels_roles
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_roles_users table
ALTER TABLE channels_roles_users
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_files table
ALTER TABLE channels_files
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_alerts table
ALTER TABLE channels_alerts
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_overlays table
ALTER TABLE channels_overlays
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_overlays_layers table
ALTER TABLE channels_overlays_layers
    ALTER COLUMN id SET DEFAULT uuidv7();

-- discord_sended_notifications table
-- This table has TEXT type but uses uuid_generate_v4(), need to convert to UUID
ALTER TABLE discord_sended_notifications
    ALTER COLUMN id TYPE UUID USING id::uuid,
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_moderation_settings table
ALTER TABLE channels_moderation_settings
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channel_duels table
ALTER TABLE channel_duels
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_integrations_seventv table
ALTER TABLE channels_integrations_seventv
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_overlays_dudes_user_settings table
ALTER TABLE channels_overlays_dudes_user_settings
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_public_settings table
ALTER TABLE channels_public_settings
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_public_settings_links table
-- This table has TEXT type but uses uuid_generate_v4(), need to convert to UUID
ALTER TABLE channels_public_settings_links
    ALTER COLUMN id TYPE UUID USING id::uuid,
    ALTER COLUMN id SET DEFAULT uuidv7();

-- badges table
ALTER TABLE badges
    ALTER COLUMN id SET DEFAULT uuidv7();

-- badges_users table
ALTER TABLE badges_users
    ALTER COLUMN id SET DEFAULT uuidv7();

-- eventsub_topics table
ALTER TABLE eventsub_topics
    ALTER COLUMN id SET DEFAULT uuidv7();

-- eventsub_subscriptions table
ALTER TABLE eventsub_subscriptions
    ALTER COLUMN id SET DEFAULT uuidv7();

-- toxic_messages table
ALTER TABLE toxic_messages
    ALTER COLUMN id SET DEFAULT uuidv7();

-- channels_commands_prefix table
ALTER TABLE channels_commands_prefix
    ALTER COLUMN id SET DEFAULT uuidv7();

-- Old tables that might still exist in some deployments
-- Check if they exist before altering

DO $$
BEGIN
    -- channel_events_list (old version)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'channel_events_list') THEN
        ALTER TABLE channel_events_list ALTER COLUMN id SET DEFAULT uuidv7();
    END IF;

    -- channels_events_list_donations
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'channels_events_list_donations') THEN
        ALTER TABLE channels_events_list_donations ALTER COLUMN id SET DEFAULT uuidv7();
    END IF;

    -- channels_events_list_follows
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'channels_events_list_follows') THEN
        ALTER TABLE channels_events_list_follows ALTER COLUMN id SET DEFAULT uuidv7();
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Rollback: change DEFAULT back to uuid_generate_v4()
-- Existing IDs remain unchanged

ALTER TABLE tokens
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_customvars
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_greetings
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_integrations
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_keywords
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_moderators
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_timers
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_timers_responses
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_events
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_events_operations
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_events_operations_filters
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_streams
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_commands
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_roles
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_roles_users
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_files
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_alerts
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_overlays
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_overlays_layers
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE discord_sended_notifications
    ALTER COLUMN id TYPE TEXT USING id::text,
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_moderation_settings
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channel_duels
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_integrations_seventv
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_overlays_dudes_user_settings
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_public_settings
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_public_settings_links
    ALTER COLUMN id TYPE TEXT USING id::text,
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE badges
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE badges_users
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE eventsub_topics
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE eventsub_subscriptions
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE toxic_messages
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE chat_messages
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE channels_commands_prefix
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'channel_events_list') THEN
        ALTER TABLE channel_events_list ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'channels_events_list_donations') THEN
        ALTER TABLE channels_events_list_donations ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'channels_events_list_follows') THEN
        ALTER TABLE channels_events_list_follows ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    END IF;
END $$;

-- +goose StatementEnd
