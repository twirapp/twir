-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers ADD COLUMN z_index INT NOT NULL DEFAULT 0;

-- Backfill existing layers so their persisted order matches the current visual
-- stacking (which was implicitly created_at ASC).
WITH ranked AS (
	SELECT
		id,
		ROW_NUMBER() OVER (PARTITION BY overlay_id ORDER BY created_at ASC, id ASC) - 1 AS new_z_index
	FROM channels_overlays_layers
)
UPDATE channels_overlays_layers AS layers
SET z_index = ranked.new_z_index
FROM ranked
WHERE layers.id = ranked.id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_overlays_layers DROP COLUMN z_index;
-- +goose StatementEnd
