-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channel_redemptions_history ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE channel_redemptions_history ALTER COLUMN redeemed_at SET DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
