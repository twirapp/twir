package postgres

import (
	"context"
	"errors"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/vk_integration"
	vkintegrationrepo "github.com/twirapp/twir/libs/repositories/vk_integration"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ vkintegrationrepo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

// scanModel is used for scanning database rows
type scanModel struct {
	ID          int64     `db:"id"`
	ChannelID   string    `db:"channel_id"`
	AccessToken string    `db:"access_token"`
	UserName    string    `db:"username"`
	Avatar      string    `db:"avatar"`
	Enabled     bool      `db:"enabled"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// toEntity converts scan model to entity
func (s scanModel) toEntity() vk_integration.Entity {
	return vk_integration.Entity{
		ID:          s.ID,
		Enabled:     s.Enabled,
		ChannelID:   s.ChannelID,
		AccessToken: s.AccessToken,
		UserName:    s.UserName,
		Avatar:      s.Avatar,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	vk_integration.Entity,
	error,
) {
	query := `
SELECT id, channel_id, access_token, username, avatar, created_at, updated_at, enabled
FROM channels_integrations_vk
WHERE channel_id = $1
LIMIT 1;
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return vk_integration.Nil, err
	}

	data, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[scanModel],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return vk_integration.Nil, vkintegrationrepo.ErrNotFound
		}
		return vk_integration.Nil, err
	}

	return data.toEntity(), nil
}

func (c *Pgx) Update(
	ctx context.Context,
	opts vkintegrationrepo.UpdateOpts,
) error {
	query := `
UPDATE channels_integrations_vk
SET
	"enabled" = COALESCE($2, enabled),
	"access_token" = COALESCE($3, "access_token"),
	username = COALESCE($4, username),
	avatar = COALESCE($5, avatar),
	updated_at = NOW()
WHERE channel_id = $1
`

	cmd, err := c.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.Enabled,
		opts.AccessToken,
		opts.UserName,
		opts.Avatar,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return vkintegrationrepo.ErrNotFound
	}

	return err
}

func (c *Pgx) Delete(ctx context.Context, channelID string) error {
	query := `
DELETE FROM channels_integrations_vk
WHERE channel_id = $1
`

	_, err := c.pool.Exec(ctx, query, channelID)
	return err
}

func (c *Pgx) Create(
	ctx context.Context,
	opts vkintegrationrepo.CreateOpts,
) error {
	query := `
INSERT INTO channels_integrations_vk (channel_id, enabled, access_token, username, avatar)
VALUES ($1, $2, $3, $4, $5)
`

	_, err := c.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.Enabled,
		opts.AccessToken,
		opts.UserName,
		opts.Avatar,
	)

	return err
}
