package pgx

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels/model"
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

var _ channels.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

const selectCols = `"id", "twitch_user_id", "kick_user_id", "isEnabled", "isTwitchBanned", "isBotMod", "botId"`

func (c *Pgx) Create(ctx context.Context, input channels.CreateInput) (model.Channel, error) {
	query := `
INSERT INTO channels (twitch_user_id, kick_user_id, "botId")
VALUES ($1, $2, $3)
RETURNING ` + selectCols

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.TwitchUserID, input.KickUserID, input.BotID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetCount(ctx context.Context, input channels.GetCountInput) (int, error) {
	query := `SELECT COUNT(*) FROM channels`
	if input.OnlyEnabled {
		query += ` WHERE "isEnabled" = true`
	}

	var count int
	err := c.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) GetByID(ctx context.Context, channelID uuid.UUID) (model.Channel, error) {
	query := `SELECT ` + selectCols + ` FROM channels WHERE "id" = $1`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByTwitchUserID(ctx context.Context, twitchUserID uuid.UUID) (model.Channel, error) {
	query := `SELECT ` + selectCols + ` FROM channels WHERE twitch_user_id = $1`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, twitchUserID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByKickUserID(ctx context.Context, kickUserID uuid.UUID) (model.Channel, error) {
	query := `SELECT ` + selectCols + ` FROM channels WHERE kick_user_id = $1`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, kickUserID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, channelID uuid.UUID, input channels.UpdateInput) (model.Channel, error) {
	updateBuilder := sq.Update("channels").Where(`"id" = ?`, channelID)

	if input.IsEnabled != nil {
		updateBuilder = updateBuilder.Set(`"isEnabled"`, *input.IsEnabled)
	}

	if input.IsBotMod != nil {
		updateBuilder = updateBuilder.Set(`"isBotMod"`, *input.IsBotMod)
	}

	updateBuilder = updateBuilder.Suffix(`RETURNING ` + selectCols)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context, input channels.GetManyInput) ([]model.Channel, error) {
	selectBuilder := sq.
		Select(
			`"id"`,
			"twitch_user_id",
			"kick_user_id",
			`"isEnabled"`,
			`"isTwitchBanned"`,
			`"isBotMod"`,
			`"botId"`,
		).
		From("channels")

	if input.Enabled != nil {
		selectBuilder = selectBuilder.Where(`"isEnabled" = ?`, *input.Enabled)
	}

	if input.HasKickUserID != nil && *input.HasKickUserID {
		selectBuilder = selectBuilder.Where("kick_user_id IS NOT NULL")
	}

	if input.HasTwitchUserID != nil && *input.HasTwitchUserID {
		selectBuilder = selectBuilder.Where("twitch_user_id IS NOT NULL")
	}

	if input.PerPage > 0 {
		selectBuilder = selectBuilder.Limit(uint64(input.PerPage))
	}

	if input.Page > 0 {
		selectBuilder = selectBuilder.Offset(uint64(input.Page * input.PerPage))
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		return nil, err
	}

	return result, nil
}
