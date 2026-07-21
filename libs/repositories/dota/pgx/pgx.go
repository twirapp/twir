package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

type Opts struct {
	PgxPool   *pgxpool.Pool
	TrManager trm.Manager
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:      opts.PgxPool,
		getter:    trmpgx.DefaultCtxGetter,
		trManager: opts.TrManager,
	}
}

func NewFx(pool *pgxpool.Pool, trManager trm.Manager) *Pgx {
	return New(Opts{PgxPool: pool, TrManager: trManager})
}

var _ dota.Repository = (*Pgx)(nil)

type Pgx struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	trManager trm.Manager
}

type matchResultExecutor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type matchStateTransitionExecutor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type predictionActionClaimExecutor interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

type predictionActionMutationExecutor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

const selectColumns = `
id, channel_id, enabled, steam_account_id, gsi_token, mmr, mmr_delta,
session_wins, session_losses, prediction_settings, chat_events, commands_settings,
created_at, updated_at
`

func (p *Pgx) scanOne(row pgx.Row) (model.ChannelDotaSettings, error) {
	settings := model.ChannelDotaSettings{}

	err := row.Scan(
		&settings.ID,
		&settings.ChannelID,
		&settings.Enabled,
		&settings.SteamAccountID,
		&settings.GsiToken,
		&settings.Mmr,
		&settings.MmrDelta,
		&settings.SessionWins,
		&settings.SessionLosses,
		&settings.PredictionSettings,
		&settings.ChatEvents,
		&settings.CommandsSettings,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, dota.ErrNotFound
		}
		return model.Nil, fmt.Errorf("dota settings scan: %w", err)
	}

	return settings, nil
}

func (p *Pgx) GetByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (model.ChannelDotaSettings, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_dota_settings
WHERE channel_id = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	return p.scanOne(conn.QueryRow(ctx, query, channelID))
}

func (p *Pgx) GetByGsiToken(
	ctx context.Context,
	token string,
) (model.ChannelDotaSettings, error) {
	query := `
SELECT ` + selectColumns + `
FROM channels_dota_settings
WHERE gsi_token = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	return p.scanOne(conn.QueryRow(ctx, query, token))
}

func (p *Pgx) Create(
	ctx context.Context,
	input dota.CreateInput,
) (model.ChannelDotaSettings, error) {
	commandsSettings := dota.CommandSettingsOrDefault(input.CommandsSettings)

	query := `
INSERT INTO channels_dota_settings (
	channel_id, enabled, steam_account_id, mmr, mmr_delta,
	prediction_settings, chat_events, commands_settings
)
VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8::jsonb, '{"mmr":true,"wl":true,"lg":true,"gm":true,"np":true,"wp":true}'::jsonb));
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.ChannelID,
		input.Enabled,
		input.SteamAccountID,
		input.Mmr,
		input.MmrDelta,
		input.PredictionSettings,
		input.ChatEvents,
		commandsSettings,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("dota settings create: %w", err)
	}

	return p.GetByChannelID(ctx, input.ChannelID)
}

func (p *Pgx) Update(
	ctx context.Context,
	channelID uuid.UUID,
	input dota.UpdateInput,
) (model.ChannelDotaSettings, error) {
	query := `
UPDATE channels_dota_settings
SET enabled = $2,
	steam_account_id = $3,
	mmr = $4,
	mmr_delta = $5,
	prediction_settings = $6,
	chat_events = $7,
	commands_settings = $8,
	updated_at = now()
WHERE channel_id = $1
RETURNING channel_id
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(
		ctx,
		query,
		channelID,
		input.Enabled,
		input.SteamAccountID,
		input.Mmr,
		input.MmrDelta,
		input.PredictionSettings,
		input.ChatEvents,
		input.CommandsSettings,
	)

	var updatedChannelID uuid.UUID
	if err := row.Scan(&updatedChannelID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, dota.ErrNotFound
		}
		return model.Nil, fmt.Errorf("dota settings update: %w", err)
	}

	return p.GetByChannelID(ctx, updatedChannelID)
}

func (p *Pgx) UpdateMatchResult(
	ctx context.Context,
	channelID uuid.UUID,
	won bool,
	mmrDelta int,
) (model.ChannelDotaSettings, error) {
	query := `
UPDATE channels_dota_settings
SET mmr = mmr + $2,
	session_wins = session_wins + $3,
	session_losses = session_losses + $4,
	updated_at = now()
WHERE channel_id = $1
RETURNING channel_id
`

	wins, losses := 0, 0
	if won {
		wins = 1
	} else {
		losses = 1
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID, mmrDelta, wins, losses)

	var updatedChannelID uuid.UUID
	if err := row.Scan(&updatedChannelID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, dota.ErrNotFound
		}
		return model.Nil, fmt.Errorf("dota settings update match result: %w", err)
	}

	return p.GetByChannelID(ctx, updatedChannelID)
}

func (p *Pgx) ApplyMatchResultOnce(
	ctx context.Context,
	input dota.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
	if err := dota.ValidateMatchResultInput(input); err != nil {
		return model.Nil, fmt.Errorf("validate dota match result: %w", err)
	}

	var settings model.ChannelDotaSettings
	err := p.trManager.Do(ctx, func(txCtx context.Context) error {
		conn := p.getter.DefaultTrOrDB(txCtx, p.pool)

		var err error
		settings, err = p.applyMatchResultOnce(txCtx, conn, input)
		return err
	})
	if err != nil {
		return model.Nil, fmt.Errorf("dota settings apply match result once: %w", err)
	}

	return settings, nil
}

func (p *Pgx) GetMatchState(ctx context.Context, channelID uuid.UUID) (model.MatchState, error) {
	if channelID == uuid.Nil {
		return model.MatchState{}, errors.New("channel ID is required")
	}

	query := `
SELECT channel_id, revision, provider_timestamp, snapshot, updated_at
FROM dota_channel_match_states
WHERE channel_id = $1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	state, err := p.scanMatchState(conn.QueryRow(ctx, query, channelID), channelID)
	if err != nil {
		return model.MatchState{}, fmt.Errorf("dota match state get: %w", err)
	}

	return state, nil
}

func (p *Pgx) scanMatchState(row pgx.Row, channelID uuid.UUID) (model.MatchState, error) {
	state := model.MatchState{}
	var snapshot []byte
	if err := row.Scan(
		&state.ChannelID,
		&state.Revision,
		&state.ProviderTimestamp,
		&snapshot,
		&state.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.MatchState{
				ChannelID: channelID,
				Snapshot:  json.RawMessage(`{}`),
			}, nil
		}

		return model.MatchState{}, fmt.Errorf("dota match state scan: %w", err)
	}

	state.Snapshot = append(json.RawMessage(nil), snapshot...)
	return state, nil
}

func (p *Pgx) ApplyMatchStateTransition(
	ctx context.Context,
	input dota.ApplyMatchStateTransitionInput,
) (bool, error) {
	if err := dota.ValidateApplyMatchStateTransitionInput(input); err != nil {
		return false, fmt.Errorf("validate dota match state transition: %w", err)
	}

	var applied bool
	err := p.trManager.Do(ctx, func(txCtx context.Context) error {
		conn := p.getter.DefaultTrOrDB(txCtx, p.pool)
		var err error
		applied, err = p.applyMatchStateTransition(txCtx, conn, input)
		return err
	})
	if err != nil {
		return false, fmt.Errorf("dota match state transition: %w", err)
	}

	return applied, nil
}

func (p *Pgx) applyMatchStateTransition(
	ctx context.Context,
	conn matchStateTransitionExecutor,
	input dota.ApplyMatchStateTransitionInput,
) (bool, error) {
	insertIdleStateQuery := `
INSERT INTO dota_channel_match_states (channel_id, snapshot)
VALUES ($1, '{}'::jsonb)
ON CONFLICT (channel_id) DO NOTHING;
`
	if _, err := conn.Exec(ctx, insertIdleStateQuery, input.ChannelID); err != nil {
		return false, fmt.Errorf("dota match state insert idle: %w", err)
	}

	selectRevisionQuery := `
SELECT revision
FROM dota_channel_match_states
WHERE channel_id = $1
FOR UPDATE;
`
	var revision int64
	if err := conn.QueryRow(ctx, selectRevisionQuery, input.ChannelID).Scan(&revision); err != nil {
		return false, fmt.Errorf("dota match state lock: %w", err)
	}
	if revision != input.ExpectedRevision {
		return false, nil
	}

	updateStateQuery := `
UPDATE dota_channel_match_states
SET revision = $2,
	provider_timestamp = $3,
	snapshot = $4,
	updated_at = now()
WHERE channel_id = $1;
`
	result, err := conn.Exec(
		ctx,
		updateStateQuery,
		input.ChannelID,
		revision+1,
		input.ProviderTimestamp,
		string(input.Snapshot),
	)
	if err != nil {
		return false, fmt.Errorf("dota match state update: %w", err)
	}
	if result.RowsAffected() != 1 {
		return false, errors.New("dota match state update affected no rows")
	}

	insertActionQuery := `
INSERT INTO dota_prediction_outbox (channel_id, match_id, action, sequence, payload)
VALUES ($1, $2, $3, $4, $5);
`
	for _, action := range input.Actions {
		if _, err := conn.Exec(
			ctx,
			insertActionQuery,
			action.ChannelID,
			action.MatchID,
			string(action.Action),
			action.Sequence,
			string(action.Payload),
		); err != nil {
			return false, fmt.Errorf("dota prediction action insert: %w", err)
		}
	}

	return true, nil
}

func (p *Pgx) ClaimPredictionActions(
	ctx context.Context,
	input dota.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	if err := dota.ValidateClaimPredictionActionsInput(input); err != nil {
		return nil, fmt.Errorf("validate claim prediction actions: %w", err)
	}

	var actions []model.ClaimedOutboxAction
	err := p.trManager.Do(ctx, func(txCtx context.Context) error {
		conn := p.getter.DefaultTrOrDB(txCtx, p.pool)
		var err error
		actions, err = p.claimPredictionActions(txCtx, conn, input)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("claim prediction actions: %w", err)
	}

	return actions, nil
}

func (p *Pgx) claimPredictionActions(
	ctx context.Context,
	conn predictionActionClaimExecutor,
	input dota.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	lockToken := uuid.New()
	query := `
WITH earliest_actions AS MATERIALIZED (
	SELECT DISTINCT ON (channel_id, match_id) id
	FROM dota_prediction_outbox
	WHERE completed_at IS NULL
	ORDER BY channel_id, match_id, sequence, created_at
), claimable_actions AS (
	SELECT outbox.id
	FROM dota_prediction_outbox AS outbox
	JOIN earliest_actions AS earliest ON earliest.id = outbox.id
	WHERE outbox.available_at <= now()
		AND (outbox.locked_at IS NULL OR outbox.locked_at < now() - ($1 * INTERVAL '1 microsecond'))
	ORDER BY outbox.available_at, outbox.sequence, outbox.created_at
	LIMIT $2
	FOR UPDATE OF outbox SKIP LOCKED
)
UPDATE dota_prediction_outbox AS outbox
SET locked_at = now(),
	lock_token = $3,
	attempts = outbox.attempts + 1
FROM claimable_actions
WHERE outbox.id = claimable_actions.id
RETURNING outbox.id, outbox.channel_id, outbox.match_id, outbox.action,
	outbox.sequence, outbox.payload, outbox.attempts, outbox.lock_token;
`

	rows, err := conn.Query(ctx, query, input.Lease.Microseconds(), input.Limit, lockToken)
	if err != nil {
		return nil, fmt.Errorf("dota prediction actions claim: %w", err)
	}
	defer rows.Close()

	var actions []model.ClaimedOutboxAction
	for rows.Next() {
		var action model.ClaimedOutboxAction
		var actionName string
		var payload []byte
		if err := rows.Scan(
			&action.ID,
			&action.ChannelID,
			&action.MatchID,
			&actionName,
			&action.Sequence,
			&payload,
			&action.Attempts,
			&action.LockToken,
		); err != nil {
			return nil, fmt.Errorf("dota prediction actions claim scan: %w", err)
		}
		action.Action = model.OutboxAction(actionName)
		action.Payload = append(json.RawMessage(nil), payload...)
		actions = append(actions, action)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dota prediction actions claim rows: %w", err)
	}

	return actions, nil
}

func (p *Pgx) CompletePredictionAction(ctx context.Context, actionID uuid.UUID, lockToken uuid.UUID) error {
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	if err := p.completePredictionAction(ctx, conn, actionID, lockToken); err != nil {
		return fmt.Errorf("complete prediction action: %w", err)
	}

	return nil
}

func (p *Pgx) completePredictionAction(
	ctx context.Context,
	conn predictionActionMutationExecutor,
	actionID uuid.UUID,
	lockToken uuid.UUID,
) error {
	query := `
UPDATE dota_prediction_outbox
SET completed_at = now(),
	locked_at = NULL,
	lock_token = NULL
WHERE id = $1
	AND lock_token = $2
	AND completed_at IS NULL;
`
	result, err := conn.Exec(ctx, query, actionID, lockToken)
	if err != nil {
		return fmt.Errorf("complete prediction action update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("complete prediction action ownership: %w", dota.ErrPredictionActionOwnershipLost)
	}

	return nil
}

func (p *Pgx) RetryPredictionAction(
	ctx context.Context,
	actionID uuid.UUID,
	lockToken uuid.UUID,
	availableAt time.Time,
) error {
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	if err := p.retryPredictionAction(ctx, conn, actionID, lockToken, availableAt); err != nil {
		return fmt.Errorf("retry prediction action: %w", err)
	}

	return nil
}

func (p *Pgx) retryPredictionAction(
	ctx context.Context,
	conn predictionActionMutationExecutor,
	actionID uuid.UUID,
	lockToken uuid.UUID,
	availableAt time.Time,
) error {
	query := `
UPDATE dota_prediction_outbox
SET available_at = $3,
	locked_at = NULL,
	lock_token = NULL
WHERE id = $1
	AND lock_token = $2
	AND completed_at IS NULL;
`
	result, err := conn.Exec(ctx, query, actionID, lockToken, availableAt)
	if err != nil {
		return fmt.Errorf("retry prediction action update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("retry prediction action ownership: %w", dota.ErrPredictionActionOwnershipLost)
	}

	return nil
}

func (p *Pgx) applyMatchResultOnce(
	ctx context.Context,
	conn matchResultExecutor,
	input dota.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
	insertSettlementQuery := `
INSERT INTO dota_match_settlements (channel_id, match_id, won, mmr_delta)
VALUES ($1, $2, $3, $4)
ON CONFLICT (channel_id, match_id) DO NOTHING;
`

	insertResult, err := conn.Exec(
		ctx,
		insertSettlementQuery,
		input.ChannelID,
		input.MatchID,
		input.Won,
		input.MmrDelta,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("dota match settlement create: %w", err)
	}

	if insertResult.RowsAffected() == 1 {
		wins, losses := 0, 0
		if input.Won {
			wins = 1
		} else {
			losses = 1
		}

		updateSettingsQuery := `
UPDATE channels_dota_settings
SET mmr = mmr + $2,
	session_wins = session_wins + $3,
	session_losses = session_losses + $4,
	updated_at = now()
WHERE channel_id = $1;
`

		updateResult, err := conn.Exec(
			ctx,
			updateSettingsQuery,
			input.ChannelID,
			input.MmrDelta,
			wins,
			losses,
		)
		if err != nil {
			return model.Nil, fmt.Errorf("dota settings apply match result: %w", err)
		}
		if updateResult.RowsAffected() == 0 {
			return model.Nil, dota.ErrNotFound
		}
	}

	getSettingsQuery := `
SELECT ` + selectColumns + `
FROM channels_dota_settings
WHERE channel_id = $1
LIMIT 1;
`

	settings, err := p.scanOne(conn.QueryRow(ctx, getSettingsQuery, input.ChannelID))
	if err != nil {
		return model.Nil, fmt.Errorf("dota settings get after match settlement: %w", err)
	}

	return settings, nil
}

func (p *Pgx) ResetSession(
	ctx context.Context,
	channelID uuid.UUID,
) (model.ChannelDotaSettings, error) {
	query := `
UPDATE channels_dota_settings
SET session_wins = 0,
	session_losses = 0,
	updated_at = now()
WHERE channel_id = $1
RETURNING channel_id
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID)

	var updatedChannelID uuid.UUID
	if err := row.Scan(&updatedChannelID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, dota.ErrNotFound
		}
		return model.Nil, fmt.Errorf("dota settings reset session: %w", err)
	}

	return p.GetByChannelID(ctx, updatedChannelID)
}

func (p *Pgx) RegenerateGsiToken(
	ctx context.Context,
	channelID uuid.UUID,
) (model.ChannelDotaSettings, error) {
	query := `
UPDATE channels_dota_settings
SET gsi_token = replace(uuidv7()::text, '-', ''),
	updated_at = now()
WHERE channel_id = $1
RETURNING channel_id
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID)

	var updatedChannelID uuid.UUID
	if err := row.Scan(&updatedChannelID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, dota.ErrNotFound
		}
		return model.Nil, fmt.Errorf("dota settings regenerate gsi token: %w", err)
	}

	return p.GetByChannelID(ctx, updatedChannelID)
}
