package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	"github.com/twirapp/twir/libs/repositories/scheduled_vips"
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

var (
	_  scheduled_vips.Repository = (*Pgx)(nil)
	sq                           = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type scanModel struct {
	ID         uuid.UUID
	UserID     string
	ChannelID  string
	CreatedAt  time.Time
	RemoveType *string
	RemoveAt   *time.Time

	isNil bool
}

func (c scanModel) toEntity() scheduledvipsentity.ScheduledVip {
	e := scheduledvipsentity.ScheduledVip{
		ID:        c.ID,
		UserID:    c.UserID,
		ChannelID: c.ChannelID,
		CreatedAt: c.CreatedAt,
		RemoveAt:  c.RemoveAt,
	}

	if c.RemoveType != nil {
		removeType := scheduledvipsentity.RemoveType(*c.RemoveType)
		e.RemoveType = &removeType
	}

	return e
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input scheduled_vips.UpdateInput) error {
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

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (scheduledvipsentity.ScheduledVip, error) {
	query := `
SELECT id, channel_id, user_id, created_at, remove_at, remove_type
FROM channels_scheduled_vips
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id.String())
	if err != nil {
		return scheduledvipsentity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scheduledvipsentity.Nil, nil
		}
		return scheduledvipsentity.Nil, err
	}

	return result.toEntity(), nil
}

func (c *Pgx) GetByUserAndChannelID(
	ctx context.Context,
	userID, channelID string,
) (scheduledvipsentity.ScheduledVip, error) {
	query := `
SELECT id, channel_id, user_id, created_at, remove_at, remove_type
FROM channels_scheduled_vips
WHERE channel_id = $1 AND user_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID, userID)
	if err != nil {
		return scheduledvipsentity.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scheduledvipsentity.Nil, nil
		}
		return scheduledvipsentity.Nil, err
	}

	return result.toEntity(), nil
}

func (c *Pgx) GetMany(ctx context.Context, input scheduled_vips.GetManyInput) (
	[]scheduledvipsentity.ScheduledVip,
	error,
) {
	builder := sq.Select(
		"id",
		"channel_id",
		"user_id",
		"created_at",
		"remove_at",
		"remove_type",
	).From("channels_scheduled_vips")

	if input.Expired != nil {
		if *input.Expired {
			builder = builder.Where("remove_at < NOW()")
		} else {
			builder = builder.Where("remove_at > NOW()")
		}
	}

	if input.RemoveType != nil {
		builder = builder.Where(squirrel.Eq{"remove_type": *input.RemoveType})
	}

	if input.ChannelID != nil {
		builder = builder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
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

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	mappedResult := make([]scheduledvipsentity.ScheduledVip, len(result))
	for i, r := range result {
		mappedResult[i] = r.toEntity()
	}

	return mappedResult, nil
}

func (c *Pgx) Create(ctx context.Context, input scheduled_vips.CreateInput) error {
	query := `
INSERT INTO channels_scheduled_vips (channel_id, user_id, remove_at, remove_type)
VALUES ($1, $2, $3, $4)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, input.ChannelID, input.UserID, input.RemoveAt, input.RemoveType)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_scheduled_vips
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id.String())
	return err
}
