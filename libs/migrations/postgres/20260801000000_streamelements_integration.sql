-- +goose NO TRANSACTION
-- +goose Up
ALTER TYPE integrations_service_enum ADD VALUE IF NOT EXISTS 'STREAMELEMENTS';
INSERT INTO integrations (service) VALUES ('STREAMELEMENTS') ON CONFLICT DO NOTHING;

-- +goose Down
-- PostgreSQL enum values cannot be removed safely while rows may reference them.
SELECT 1;
