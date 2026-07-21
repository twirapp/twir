-- +goose Up
-- +goose StatementBegin
CREATE TABLE dota_channel_match_states (
    channel_id UUID PRIMARY KEY REFERENCES channels(id) ON DELETE CASCADE,
    revision BIGINT NOT NULL DEFAULT 0 CHECK (revision >= 0),
    provider_timestamp BIGINT NOT NULL DEFAULT 0,
    snapshot JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE dota_prediction_outbox (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    match_id BIGINT NOT NULL CHECK (match_id > 0),
    action TEXT NOT NULL CHECK (action IN ('create', 'resolve', 'cancel')),
    sequence BIGINT NOT NULL CHECK (sequence > 0),
    payload JSONB NOT NULL,
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_at TIMESTAMPTZ,
    lock_token UUID,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (channel_id, match_id, action)
);

CREATE INDEX dota_prediction_outbox_claim_idx
    ON dota_prediction_outbox (available_at, sequence, created_at)
    WHERE completed_at IS NULL;
CREATE INDEX dota_prediction_outbox_match_order_idx
    ON dota_prediction_outbox (channel_id, match_id, sequence)
    WHERE completed_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dota_prediction_outbox;
DROP TABLE IF EXISTS dota_channel_match_states;
-- +goose StatementEnd
