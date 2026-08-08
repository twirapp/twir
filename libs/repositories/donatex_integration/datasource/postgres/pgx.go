package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	donatexintegrationentity "github.com/twirapp/twir/libs/entities/donatex_integration"
	donatexintegration "github.com/twirapp/twir/libs/repositories/donatex_integration"
)

var _ donatexintegration.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(pool)
}

func New(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

type scanModel struct {
	ID            int64
	PublicID      uuid.UUID
	Enabled       bool
	ChannelID     string
	AccessToken   string
	RefreshToken  string
	DonateXUserID string
	UserName      string
	Avatar        string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	isNil bool
}

func (c scanModel) toEntity() donatexintegrationentity.Entity {
	return donatexintegrationentity.Entity{
		ID:            c.ID,
		PublicID:      c.PublicID,
		Enabled:       c.Enabled,
		ChannelID:     c.ChannelID,
		AccessToken:   c.AccessToken,
		RefreshToken:  c.RefreshToken,
		DonateXUserID: c.DonateXUserID,
		UserName:      c.UserName,
		Avatar:        c.Avatar,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func (p *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	donatexintegrationentity.Entity,
	error,
) {
	query := `
		SELECT
			id,
			public_id,
			channel_id,
			access_token,
			refresh_token,
			donatex_user_id,
			username,
			avatar,
			enabled,
			created_at,
			updated_at
		FROM channels_integrations_donatex
		WHERE channel_id = $1
		LIMIT 1
	`

	rows, err := p.pool.Query(ctx, query, channelID)
	if err != nil {
		return donatexintegrationentity.Nil, fmt.Errorf(
			"GetByChannelID: failed to execute query: %w",
			err,
		)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return donatexintegrationentity.Nil, nil
		}
		return donatexintegrationentity.Nil, fmt.Errorf("GetByChannelID: failed to collect row: %w", err)
	}

	return result.toEntity(), nil
}

func (p *Pgx) Create(ctx context.Context, opts donatexintegration.CreateOpts) error {
	query := `
		INSERT INTO channels_integrations_donatex (
			channel_id,
			access_token,
			refresh_token,
			donatex_user_id,
			username,
			avatar,
			enabled,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err := p.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.AccessToken,
		opts.RefreshToken,
		opts.DonateXUserID,
		opts.UserName,
		opts.Avatar,
		opts.Enabled,
	)
	if err != nil {
		return fmt.Errorf("Create: failed to insert donatex integration: %w", err)
	}

	return nil
}

func (p *Pgx) Update(ctx context.Context, opts donatexintegration.UpdateOpts) error {
	builder := sq.Update("channels_integrations_donatex").
		Where(squirrel.Eq{"channel_id": opts.ChannelID}).
		Set("updated_at", squirrel.Expr("NOW()"))

	if opts.AccessToken != nil {
		builder = builder.Set("access_token", *opts.AccessToken)
	}
	if opts.RefreshToken != nil {
		builder = builder.Set("refresh_token", *opts.RefreshToken)
	}
	if opts.DonateXUserID != nil {
		builder = builder.Set("donatex_user_id", *opts.DonateXUserID)
	}
	if opts.Enabled != nil {
		builder = builder.Set("enabled", *opts.Enabled)
	}
	if opts.UserName != nil {
		builder = builder.Set("username", *opts.UserName)
	}
	if opts.Avatar != nil {
		builder = builder.Set("avatar", *opts.Avatar)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("Update: failed to build query: %w", err)
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Update: failed to update donatex integration: %w", err)
	}

	return nil
}

func (p *Pgx) Delete(ctx context.Context, channelID string) error {
	query := `
		DELETE FROM channels_integrations_donatex
		WHERE channel_id = $1
	`

	_, err := p.pool.Exec(ctx, query, channelID)
	if err != nil {
		return fmt.Errorf("Delete: failed to delete donatex integration: %w", err)
	}

	return nil
}
