-- +goose Up
-- +goose StatementBegin

-- Add stack_id column for grouping widgets into tabs
ALTER TABLE "channels_dashboard_widgets" ADD COLUMN "stack_id" TEXT;

-- Add stack_order column for ordering tabs within a stack
ALTER TABLE "channels_dashboard_widgets" ADD COLUMN "stack_order" INT NOT NULL DEFAULT 0;

-- Create index for efficient stack lookups
CREATE INDEX "channels_dashboard_widgets_stack_id_idx" ON "channels_dashboard_widgets" ("channel_id", "stack_id") WHERE "stack_id" IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS "channels_dashboard_widgets_stack_id_idx";
ALTER TABLE "channels_dashboard_widgets" DROP COLUMN IF EXISTS "stack_order";
ALTER TABLE "channels_dashboard_widgets" DROP COLUMN IF EXISTS "stack_id";

-- +goose StatementEnd
