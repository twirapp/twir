package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
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

const getAllByBindingPlatformQuery = selectQuery + `
WHERE EXISTS (
	SELECT 1
	FROM channel_platforms cp_filter
	WHERE cp_filter.channel_id = c.id
		AND cp_filter.platform = $1
)
ORDER BY c.id`

const createChannelQuery = `
INSERT INTO channels DEFAULT VALUES
RETURNING id`

const updateChannelQuery = `
UPDATE channels
SET "isEnabled" = COALESCE($2, "isEnabled")
WHERE id = $1
RETURNING id`

func (c *Pgx) GetByApiKey(ctx context.Context, apiKey string) (channelentity.Channel, error) {
	return c.getOne(ctx, selectQuery+` WHERE c.api_key = $1`, apiKey)
}

func (c *Pgx) Create(ctx context.Context) (channelentity.Channel, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	row := conn.QueryRow(ctx, createChannelQuery)

	var channelId uuid.UUID
	if err := row.Scan(&channelId); err != nil {
		return channelentity.Nil, err
	}

	return c.GetByID(ctx, channelId)
}

func (c *Pgx) GetByID(ctx context.Context, channelID uuid.UUID) (channelentity.Channel, error) {
	return c.getOne(ctx, selectQuery+` WHERE c."id" = $1`, channelID)
}

func (c *Pgx) GetAllByBindingPlatform(
	ctx context.Context,
	p platform.Platform,
) ([]channelentity.Channel, error) {
	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, getAllByBindingPlatformQuery, p)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]channelentity.Channel, 0)
	for rows.Next() {
		channel, err := scanChannel(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, mapChannelToEntity(channel))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetByBindingUserID(
	ctx context.Context,
	p platform.Platform,
	userID uuid.UUID,
) (channelentity.Channel, error) {
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
) (channelentity.Channel, error) {
	return c.getOne(
		ctx,
		selectQuery+`
			JOIN channel_platforms cp ON cp.channel_id = c.id
			WHERE cp.platform = $1 AND cp.platform_channel_id = $2`,
		p,
		platformChannelID,
	)
}

func (c *Pgx) Update(ctx context.Context, channelID uuid.UUID, input channels.UpdateInput) (channelentity.Channel, error) {
	row := c.getter.DefaultTrOrDB(ctx, c.pool).QueryRow(
		ctx,
		updateChannelQuery,
		channelID,
		valueOrNil(input.IsEnabled),
	)

	var channelId uuid.UUID
	if err := row.Scan(&channelId); err != nil {
		return channelentity.Nil, err
	}

	return c.GetByID(ctx, channelId)
}

func valueOrNil[T any](value *T) any {
	if value == nil {
		return nil
	}

	return *value
}

func (c *Pgx) GetMany(ctx context.Context, input channels.GetManyInput) ([]channelentity.Channel, error) {
	query := selectQuery

	var where []string
	var args []any

	if input.Enabled != nil {
		where = append(where, `c."isEnabled" = $`+strconv.Itoa(len(args)+1))
		args = append(args, *input.Enabled)
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

	result := make([]channelentity.Channel, 0)
	for rows.Next() {
		channel, err := scanChannel(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, mapChannelToEntity(channel))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetBySlug(ctx context.Context, opts channels.GetBySlugInput) (channelentity.Channel, error) {
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

func scanChannel(row pgx.CollectableRow) (model.Channel, error) {
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

func collectExactlyOneChannel(rows pgx.Rows) (model.Channel, error) {
	return pgx.CollectExactlyOneRow(rows, scanChannel)
}

func (c *Pgx) getOne(ctx context.Context, query string, args ...any) (channelentity.Channel, error) {
	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, args...)
	if err != nil {
		return channelentity.Nil, err
	}

	channel, err := collectExactlyOneChannel(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return channelentity.Nil, channels.ErrNotFound
		}
		return channelentity.Nil, err
	}

	return mapChannelToEntity(channel), nil
}
