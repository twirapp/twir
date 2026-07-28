-- +goose Up
-- +goose NO TRANSACTION
ALTER TYPE channel_overlay_now_playing_preset ADD VALUE IF NOT EXISTS 'PULSE_STRIP';
ALTER TYPE channel_overlay_now_playing_preset ADD VALUE IF NOT EXISTS 'AURA_STACK';
ALTER TYPE channel_overlay_now_playing_preset ADD VALUE IF NOT EXISTS 'VINYL_HAZE';
ALTER TYPE channel_overlay_now_playing_preset ADD VALUE IF NOT EXISTS 'SIGNAL_DECK';

-- +goose Down
-- +goose NO TRANSACTION
-- PostgreSQL does not support removing values from an enum type; left as a no-op.
SELECT 1;
