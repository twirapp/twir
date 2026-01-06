-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL UNIQUE,
    max_commands INTEGER NOT NULL DEFAULT 500,
    max_timers INTEGER NOT NULL DEFAULT 50,
    max_variables INTEGER NOT NULL DEFAULT 50,
    max_alerts INTEGER NOT NULL DEFAULT 50,
    max_events INTEGER NOT NULL DEFAULT 50,
    max_chat_alerts_messages INTEGER NOT NULL DEFAULT 20,
    max_custom_overlays INTEGER NOT NULL DEFAULT 10,
    max_eightball_answers INTEGER NOT NULL DEFAULT 25,
    max_commands_responses INTEGER NOT NULL DEFAULT 3,
    max_moderation_rules INTEGER NOT NULL DEFAULT 50,
    max_keywords INTEGER NOT NULL DEFAULT 50,
    max_greetings INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert default free plan
INSERT INTO plans (
    name,
    max_commands,
    max_timers,
    max_variables,
    max_alerts,
    max_events,
    max_chat_alerts_messages,
    max_custom_overlays,
    max_eightball_answers,
    max_commands_responses,
    max_moderation_rules,
    max_keywords,
    max_greetings
)
VALUES (
    'Free',
    500,
    50,
    50,
    50,
    50,
    20,
    10,
    25,
    3,
    50,
    50,
    100
);

-- Add plan_id column to channels table
ALTER TABLE channels
ADD COLUMN plan_id UUID REFERENCES plans(id) ON DELETE SET NULL;

-- Set all existing channels to free plan
UPDATE channels
SET plan_id = (SELECT id FROM plans WHERE name = 'free');

-- Create function to set default plan
CREATE OR REPLACE FUNCTION set_default_plan_for_channel()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.plan_id IS NULL THEN
        NEW.plan_id := (SELECT id FROM plans WHERE name = 'free' LIMIT 1);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically set plan for new channels
CREATE TRIGGER trigger_set_default_plan
    BEFORE INSERT ON channels
    FOR EACH ROW
    EXECUTE FUNCTION set_default_plan_for_channel();

-- Add index for better performance
CREATE INDEX idx_channels_plan_id ON channels(plan_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_set_default_plan ON channels;
DROP FUNCTION IF EXISTS set_default_plan_for_channel();
DROP INDEX IF EXISTS idx_channels_plan_id;
ALTER TABLE channels DROP COLUMN plan_id;
DROP TABLE plans;
-- +goose StatementEnd
