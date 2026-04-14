package pgx

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
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

func (c *Pgx) Create(ctx context.Context, input channels.CreateInput) (model.Channel, error) {
	query := `
INSERT INTO channels (user_id, platform, "botId")
VALUES ($1, $2, $3)
RETURNING "id"::text AS "id", "platform", "user_id", "isEnabled", "isTwitchBanned", "isBotMod", "botId"
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.UserID, input.Platform, input.BotID)
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
	query := `
SELECT COUNT(*)
FROM channels
`

	if input.OnlyEnabled {
		query += `
WHERE "isEnabled" = true
`
	}

	var count int
	err := c.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) GetByID(ctx context.Context, channelID string) (model.Channel, error) {
	query := `
SELECT "id"::text AS "id", "platform", "user_id", "isEnabled", "isTwitchBanned", "isBotMod", "botId"
FROM channels
WHERE "id" = $1::uuid
`

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
	}

	return result, nil
}

func (c *Pgx) GetByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, platformVal platform.Platform) (model.Channel, error) {
	query := `
SELECT "id"::text AS "id", "platform", "user_id", "isEnabled", "isTwitchBanned", "isBotMod", "botId"
FROM channels
WHERE "user_id" = $1 AND "platform" = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, userID, platformVal)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
	}

	return result, nil
}

func (c *Pgx) GetByPlatformUserID(ctx context.Context, plat platform.Platform, platformUserID string) (model.Channel, error) {
	query := `
SELECT c."id"::text AS "id", c."platform", c."user_id", c."isEnabled", c."isTwitchBanned", c."isBotMod", c."botId"
FROM user_platform_accounts upa
JOIN channels c ON c.user_id = upa.user_id AND c.platform = $1
WHERE upa.platform = $1 AND upa.platform_user_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, plat, platformUserID)
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

func (c *Pgx) Update(ctx context.Context, channelID string, input channels.UpdateInput) (model.Channel, error) {
	updateBuilder := sq.Update("channels").Where(`"id" = ?::uuid`, channelID)

	if input.IsEnabled != nil {
		updateBuilder = updateBuilder.Set(`"isEnabled"`, *input.IsEnabled)
	}

	if input.IsBotMod != nil {
		updateBuilder = updateBuilder.Set(`"isBotMod"`, *input.IsBotMod)
	}

	updateBuilder = updateBuilder.Suffix(
		`RETURNING "id"::text AS "id", "platform", "user_id", "isEnabled", "isTwitchBanned", "isBotMod", "botId"`,
	)

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
			`"id"::text AS "id"`,
			"platform",
			"user_id",
			`"isEnabled"`,
			`"isTwitchBanned"`,
			`"isBotMod"`,
			`"botId"`,
		).
		From("channels")

	if input.Enabled != nil {
		selectBuilder = selectBuilder.Where(`"isEnabled" = ?`, *input.Enabled)
	}

	// not need to use defaults because i wanna select all channels
	if input.PerPage > 0 {
		selectBuilder = selectBuilder.Limit(uint64(input.PerPage))
	}

	// not need to use defaults because i wanna select all channels
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
