package pgx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	c.kick_bot_id,
	c.api_key,
	json_build_object(
		'id', tu.id,
		'platform', tu.platform,
		'platform_id', tu.platform_id,
		'token_id', tu."tokenId",
		'is_bot_admin', tu."isBotAdmin",
		'api_key', tu."apiKey",
		'is_banned', tu.is_banned,
		'hide_on_landing_page', tu.hide_on_landing_page,
		'created_at', tu.created_at,
		'login', tu.login,
		'display_name', tu.display_name,
		'avatar', tu.avatar
	) as twitch_user,
	json_build_object(
		'id', ku.id,
		'platform', ku.platform,
		'platform_id', ku.platform_id,
		'token_id', ku."tokenId",
		'is_bot_admin', ku."isBotAdmin",
		'api_key', ku."apiKey",
		'is_banned', ku.is_banned,
		'hide_on_landing_page', ku.hide_on_landing_page,
		'created_at', ku.created_at,
		'login', ku.login,
		'display_name', ku.display_name,
		'avatar', ku.avatar
	)	as kick_user
FROM channels c
LEFT JOIN users tu ON tu.id = c.twitch_user_id AND tu.platform = 'twitch'
LEFT JOIN users ku ON ku.id = c.kick_user_id AND ku.platform = 'kick'`

func (c *Pgx) GetByApiKey(ctx context.Context, apiKey string) (model.Channel, error) {
	rows, err := c.pool.Query(ctx, selectQuery+`WHERE api_key = $1`, apiKey)
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

func (c *Pgx) Create(ctx context.Context, input channels.CreateInput) (model.Channel, error) {
	query := `INSERT INTO channels (twitch_user_id, kick_user_id, twitch_bot_enabled, kick_bot_enabled, "isEnabled", "botId", kick_bot_id)
	VALUES ($1, $2, $3, $4, $3 OR $4, $5, $6)
	RETURNING id`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	row := conn.QueryRow(
		ctx,
		query,
		input.TwitchUserID,
		input.KickUserID,
		input.TwitchBotEnabled,
		input.KickBotEnabled,
		input.BotID,
		input.KickBotID,
	)

	var channelId uuid.UUID
	if err := row.Scan(&channelId); err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, channelId)
}

func (c *Pgx) GetCount(ctx context.Context, input channels.GetCountInput) (int, error) {
	query := `SELECT COUNT(*) FROM channels`
	if input.OnlyTwitchEnabled {
		query += ` WHERE twitch_bot_enabled = true`
	} else if input.OnlyEnabled {
		query += ` WHERE "isEnabled" = true`
	}

	var count int
	err := c.getter.DefaultTrOrDB(ctx, c.pool).QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) GetByID(ctx context.Context, channelID uuid.UUID) (model.Channel, error) {
	query := selectQuery + ` WHERE c."id" = $1`

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, channelID)
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

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, twitchUserID)
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

func (c *Pgx) GetByTwitchPlatformID(ctx context.Context, twitchPlatformID string) (model.Channel, error) {
	query := selectQuery + ` WHERE tu.platform_id = $1`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, twitchPlatformID)
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

func (c *Pgx) GetByKickPlatformID(ctx context.Context, kickPlatformID string) (model.Channel, error) {
	query := selectQuery + ` WHERE ku.platform_id = $1`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, kickPlatformID)
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

	updateBuilder = updateBuilder.Suffix(`RETURNING id`)

	innerQuery, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	row := c.getter.DefaultTrOrDB(ctx, c.pool).QueryRow(ctx, innerQuery, args...)

	var channelId uuid.UUID
	if err := row.Scan(&channelId); err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, channelId)
}

func (c *Pgx) GetMany(ctx context.Context, input channels.GetManyInput) ([]model.Channel, error) {
	query := selectQuery

	var where []string
	var args []any

	if input.Enabled != nil {
		where = append(where, `c."isEnabled" = $`+strconv.Itoa(len(args)+1))
		args = append(args, *input.Enabled)
	}

	if input.TwitchBotEnabled != nil {
		where = append(where, `c.twitch_bot_enabled = $`+strconv.Itoa(len(args)+1))
		args = append(args, *input.TwitchBotEnabled)
	}

	if input.AnyBotEnabled != nil && *input.AnyBotEnabled {
		where = append(where, `(c.twitch_bot_enabled OR c.kick_bot_enabled)`)
	}

	if len(where) > 0 {
		query += "\nWHERE " + strings.Join(where, "\nAND ")
	}

	if input.PerPage == 0 {
		input.PerPage = 10
	}

	if input.PerPage > 0 {
		args = append(args, input.PerPage)
		query += fmt.Sprintf("\nLIMIT $%d", len(args))
	}

	if input.Page > 0 {
		args = append(args, input.Page*input.PerPage)
		query += fmt.Sprintf("\nOFFSET $%d", len(args))
	}

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetBySlug(ctx context.Context, opts channels.GetBySlugInput) (model.Channel, error) {
	query := selectQuery
	args := []any{opts.Slug}

	if opts.Platform == nil {
		query = query + ` WHERE tu."login" = $1 OR ku."login" = $1`
	} else {
		switch *opts.Platform {
		case platform.PlatformKick:
			query = query + ` WHERE ku."login" = $1`
			query = query + ` AND ku."platform" = $2`
		case platform.PlatformTwitch:
			query = query + ` WHERE tu."login" = $1`
			query = query + ` AND tu."platform" = $2`
		}

		args = append(args, *opts.Platform)
	}

	query = query + " LIMIT 1"

	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, args...)
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
