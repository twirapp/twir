package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	faceitintegrationentity "github.com/twirapp/twir/libs/entities/faceit_integration"
	faceitintegration "github.com/twirapp/twir/libs/repositories/faceit_integration"
)

var _ faceitintegration.Repository = (*Pgx)(nil)
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
	ID           int64
	Enabled      bool
	ChannelID    string
	AccessToken  string
	RefreshToken string
	UserName     string
	Avatar       string
	Game         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	FaceitUserID string

	isNil bool
}

func (c scanModel) toEntity() faceitintegrationentity.Entity {
	return faceitintegrationentity.Entity{
		ID:           c.ID,
		Enabled:      c.Enabled,
		ChannelID:    c.ChannelID,
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		UserName:     c.UserName,
		Avatar:       c.Avatar,
		Game:         c.Game,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		FaceitUserID: c.FaceitUserID,
	}
}

func (p *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	faceitintegrationentity.Entity,
	error,
) {
	query := `
		SELECT
			id,
			channel_id,
			access_token,
			refresh_token,
			username,
			avatar,
			game,
			enabled,
			created_at,
			updated_at,
			faceit_user_id
		FROM channels_integrations_faceit
		WHERE channel_id = $1
		LIMIT 1
	`

	rows, err := p.pool.Query(ctx, query, channelID)
	if err != nil {
		return faceitintegrationentity.Nil, fmt.Errorf(
			"GetByChannelID: failed to execute query: %w",
			err,
		)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return faceitintegrationentity.Nil, nil
		}
		return faceitintegrationentity.Nil, fmt.Errorf("GetByChannelID: failed to collect row: %w", err)
	}

	return result.toEntity(), nil
}

func (p *Pgx) Create(ctx context.Context, opts faceitintegration.CreateOpts) error {
	query := `
		INSERT INTO channels_integrations_faceit (
			channel_id,
			access_token,
			refresh_token,
			username,
			avatar,
			game,
			enabled,
			created_at,
			updated_at,
			faceit_user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), $8)
	`

	_, err := p.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.AccessToken,
		opts.RefreshToken,
		opts.UserName,
		opts.Avatar,
		opts.Game,
		opts.Enabled,
		opts.FaceitUserID,
	)
	if err != nil {
		return fmt.Errorf("Create: failed to insert faceit integration: %w", err)
	}

	return nil
}

func (p *Pgx) Update(ctx context.Context, opts faceitintegration.UpdateOpts) error {
	builder := sq.Update("channels_integrations_faceit").
		Where(squirrel.Eq{"channel_id": opts.ChannelID}).
		Set("updated_at", squirrel.Expr("NOW()"))

	if opts.AccessToken != nil {
		builder = builder.Set("access_token", *opts.AccessToken)
	}
	if opts.RefreshToken != nil {
		builder = builder.Set("refresh_token", *opts.RefreshToken)
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
	if opts.Game != nil {
		builder = builder.Set("game", *opts.Game)
	}
	if opts.FaceitUserID != nil {
		builder = builder.Set("faceit_user_id", *opts.FaceitUserID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("Update: failed to build query: %w", err)
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Update: failed to update faceit integration: %w", err)
	}

	return nil
}

func (p *Pgx) Delete(ctx context.Context, channelID string) error {
	query := `
		DELETE FROM channels_integrations_faceit
		WHERE channel_id = $1
	`

	_, err := p.pool.Exec(ctx, query, channelID)
	if err != nil {
		return fmt.Errorf("Delete: failed to delete faceit integration: %w", err)
	}

	return nil
}
