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

const selectQuery = `
SELECT
	c."id",
	c."twitch_user_id",
	tu.platform_id AS twitch_platform_id,
	c.twitch_bot_enabled,
	c."kick_user_id",
	ku.platform_id AS kick_platform_id,
	c.kick_bot_enabled,
	c."isEnabled",
	c."isTwitchBanned",
	c."isBotMod",
	c."botId",
	c.kick_bot_id
FROM channels c
LEFT JOIN users tu ON tu.id = c.twitch_user_id
LEFT JOIN users ku ON ku.id = c.kick_user_id`

func (c *Pgx) Create(ctx context.Context, input channels.CreateInput) (model.Channel, error) {
	query := `
WITH inserted AS (
	INSERT INTO channels (twitch_user_id, kick_user_id, twitch_bot_enabled, kick_bot_enabled, "botId", kick_bot_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
)
SELECT
	i."id",
	i."twitch_user_id",
	tu.platform_id AS twitch_platform_id,
	i.twitch_bot_enabled,
	i."kick_user_id",
	ku.platform_id AS kick_platform_id,
	i.kick_bot_enabled,
	i."isEnabled",
	i."isTwitchBanned",
	i."isBotMod",
	i."botId",
	i.kick_bot_id
FROM inserted i
LEFT JOIN users tu ON tu.id = i.twitch_user_id
LEFT JOIN users ku ON ku.id = i.kick_user_id`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.TwitchUserID, input.KickUserID, input.TwitchBotEnabled, input.KickBotEnabled, input.BotID, input.KickBotID)
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
	if input.OnlyTwitchEnabled {
		query += ` WHERE twitch_bot_enabled = true`
	} else if input.OnlyEnabled {
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
	query := selectQuery + ` WHERE c."id" = $1`

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
	query := selectQuery + ` WHERE c.twitch_user_id = $1`

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
	query := selectQuery + ` WHERE c.kick_user_id = $1`

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

	if input.TwitchUserID != nil {
		updateBuilder = updateBuilder.Set("twitch_user_id", *input.TwitchUserID)
	}

	if input.KickUserID != nil {
		updateBuilder = updateBuilder.Set("kick_user_id", *input.KickUserID)
	}

	if input.TwitchBotEnabled != nil {
		updateBuilder = updateBuilder.Set("twitch_bot_enabled", *input.TwitchBotEnabled)
	}

	if input.KickBotEnabled != nil {
		updateBuilder = updateBuilder.Set("kick_bot_enabled", *input.KickBotEnabled)
	}

	if input.KickBotID != nil {
		updateBuilder = updateBuilder.Set("kick_bot_id", *input.KickBotID)
	}

	updateBuilder = updateBuilder.Suffix(`RETURNING *`)

	innerQuery, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	query := `
WITH updated AS (` + innerQuery + `)
SELECT
	u."id",
	u."twitch_user_id",
	tu.platform_id AS twitch_platform_id,
	u.twitch_bot_enabled,
	u."kick_user_id",
	ku.platform_id AS kick_platform_id,
	u.kick_bot_enabled,
	u."isEnabled",
	u."isTwitchBanned",
	u."isBotMod",
	u."botId",
	u.kick_bot_id
FROM updated u
LEFT JOIN users tu ON tu.id = u.twitch_user_id
LEFT JOIN users ku ON ku.id = u.kick_user_id`

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
			`c."id"`,
			"c.twitch_user_id",
			"tu.platform_id AS twitch_platform_id",
			"c.twitch_bot_enabled",
			"c.kick_user_id",
			"ku.platform_id AS kick_platform_id",
			"c.kick_bot_enabled",
			`c."isEnabled"`,
			`c."isTwitchBanned"`,
			`c."isBotMod"`,
			`c."botId"`,
			"c.kick_bot_id",
		).
		From("channels c").
		LeftJoin("users tu ON tu.id = c.twitch_user_id").
		LeftJoin("users ku ON ku.id = c.kick_user_id")

	if input.Enabled != nil {
		selectBuilder = selectBuilder.Where(`c."isEnabled" = ?`, *input.Enabled)
	}

	if input.HasKickUserID != nil && *input.HasKickUserID {
		selectBuilder = selectBuilder.Where("c.kick_user_id IS NOT NULL")
	}

	if input.HasTwitchUserID != nil && *input.HasTwitchUserID {
		selectBuilder = selectBuilder.Where("c.twitch_user_id IS NOT NULL")
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
