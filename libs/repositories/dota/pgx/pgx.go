package pgx

import (
	"context"
	"errors"
	"fmt"

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
