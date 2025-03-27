package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
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

var _ scheduled_vips.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) Update(ctx context.Context, id ulid.ULID, input scheduled_vips.UpdateInput) error {
	updateBuilder := sq.Update("channels_scheduled_vips").
		Where(squirrel.Eq{"id": id.String()})

	if input.RemoveAt != nil {
		updateBuilder = updateBuilder.Set("remove_at", *input.RemoveAt)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) GetByID(ctx context.Context, id ulid.ULID) (model.ScheduledVip, error) {
	query := `
SELECT id, channel_id, user_id, created_at, remove_at
FROM channels_scheduled_vips
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ScheduledVip])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByUserAndChannelID(
	ctx context.Context,
	userID, channelID string,
) (model.ScheduledVip, error) {
	query := `
SELECT id, channel_id, user_id, created_at, remove_at
FROM channels_scheduled_vips
WHERE channel_id = $1 AND user_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ScheduledVip])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context, input scheduled_vips.GetManyInput) (
	[]model.ScheduledVip,
	error,
) {
	builder := sq.Select(
		"id",
		"channel_id",
		"user_id",
		"created_at",
		"remove_at",
	).From("channels_scheduled_vips")

	if input.Expired != nil {
		if *input.Expired {
			builder = builder.Where("remove_at < NOW()")
		} else {
			builder = builder.Where("remove_at > NOW()")
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ScheduledVip])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) (
	[]model.ScheduledVip,
	error,
) {
	query := `
SELECT id, channel_id, user_id, created_at, remove_at
FROM channels_scheduled_vips
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ScheduledVip])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input scheduled_vips.CreateInput) error {
	query := `
INSERT INTO channels_scheduled_vips (channel_id, user_id, remove_at)
VALUES ($1, $2, $3)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, input.ChannelID, input.UserID, input.RemoveAt)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id ulid.ULID) error {
	query := `
DELETE FROM channels_scheduled_vips
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id.String())
	return err
}
