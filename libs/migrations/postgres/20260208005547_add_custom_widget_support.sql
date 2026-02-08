-- +goose Up
-- +goose StatementBegin
-- Add columns for custom widget support
ALTER TABLE channels_dashboard_widgets 
  ADD COLUMN type TEXT NOT NULL DEFAULT 'system',
  ADD COLUMN custom_name TEXT,
  ADD COLUMN custom_url TEXT;

-- Add check constraint for type
ALTER TABLE channels_dashboard_widgets
  ADD CONSTRAINT check_widget_type CHECK (type IN ('system', 'custom'));

-- Add constraint: custom widgets must have custom_name and custom_url
ALTER TABLE channels_dashboard_widgets
  ADD CONSTRAINT check_custom_fields CHECK (
    (type = 'system' AND custom_name IS NULL AND custom_url IS NULL) OR
    (type = 'custom' AND custom_name IS NOT NULL AND custom_url IS NOT NULL)
  );

-- Drop the old custom_dashboard_widgets table if it exists
DROP TABLE IF EXISTS custom_dashboard_widgets;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove constraints
ALTER TABLE channels_dashboard_widgets
  DROP CONSTRAINT IF EXISTS check_custom_fields,
  DROP CONSTRAINT IF EXISTS check_widget_type;

-- Remove columns
ALTER TABLE channels_dashboard_widgets
  DROP COLUMN IF EXISTS custom_url,
  DROP COLUMN IF EXISTS custom_name,
  DROP COLUMN IF EXISTS type;

-- Recreate custom_dashboard_widgets table
CREATE TABLE IF NOT EXISTS custom_dashboard_widgets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  channel_id TEXT NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_custom_dashboard_widgets_channel_id ON custom_dashboard_widgets(channel_id);
-- +goose StatementEnd
