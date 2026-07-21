package pgx

import (
	"context"
	"encoding/json"
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
	c.api_key,
	COALESCE(
		(
			SELECT jsonb_agg(
				jsonb_build_object(
					'ID', cp.id,
					'ChannelID', cp.channel_id,
					'Platform', cp.platform,
					'UserID', cp.user_id,
					'PlatformChannelID', cp.platform_channel_id,
					'Enabled', cp.enabled,
					'BotUserID', cp.bot_user_id,
					'BotConfig', cp.bot_config,
					'CreatedAt', cp.created_at,
					'UpdatedAt', cp.updated_at
				)
				ORDER BY cp.platform
			)
			FROM channel_platforms cp
			WHERE cp.channel_id = c.id
		),
		'[]'::jsonb
	) AS bindings
FROM channels c`

func (c *Pgx) GetByApiKey(ctx context.Context, apiKey string) (model.Channel, error) {
	return c.getOne(ctx, selectQuery+` WHERE c.api_key = $1`, apiKey)
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
	return c.getOne(ctx, selectQuery+` WHERE c."id" = $1`, channelID)
}

func (c *Pgx) GetByBindingUserID(
	ctx context.Context,
	p platform.Platform,
	userID uuid.UUID,
) (model.Channel, error) {
	return c.getOne(
		ctx,
		selectQuery+`
			JOIN channel_platforms cp ON cp.channel_id = c.id
			WHERE cp.platform = $1 AND cp.user_id = $2`,
		p,
		userID,
	)
}

func (c *Pgx) GetByPlatformChannelID(
	ctx context.Context,
	p platform.Platform,
	platformChannelID string,
) (model.Channel, error) {
	return c.getOne(
		ctx,
		selectQuery+`
			JOIN channel_platforms cp ON cp.channel_id = c.id
			WHERE cp.platform = $1 AND cp.platform_channel_id = $2`,
		p,
		platformChannelID,
	)
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
	defer rows.Close()

	result := make([]model.Channel, 0)
	for rows.Next() {
		channel, err := scanChannel(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, channel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetBySlug(ctx context.Context, opts channels.GetBySlugInput) (model.Channel, error) {
	query := selectQuery + `
		JOIN channel_platforms cp ON cp.channel_id = c.id
		JOIN users u ON u.id = cp.user_id
		WHERE u.login = $1`
	args := []any{opts.Slug}

	if opts.Platform != nil {
		query += ` AND cp.platform = $2`
		args = append(args, *opts.Platform)
	}

	return c.getOne(ctx, query+" LIMIT 1", args...)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanChannel(row rowScanner) (model.Channel, error) {
	var channel model.Channel
	var bindingsJSON []byte
	if err := row.Scan(&channel.ID, &channel.ApiKey, &bindingsJSON); err != nil {
		return model.Nil, err
	}
	if err := json.Unmarshal(bindingsJSON, &channel.Bindings); err != nil {
		return model.Nil, fmt.Errorf("unmarshal channel platform bindings: %w", err)
	}

	return channel, nil
}

func (c *Pgx) getOne(ctx context.Context, query string, args ...any) (model.Channel, error) {
	channel, err := scanChannel(c.getter.DefaultTrOrDB(ctx, c.pool).QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
		return model.Nil, err
	}

	return channel, nil
}
