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

const createChannelAndBindingsQuery = `
WITH requested AS (
	SELECT
		$1::uuid AS twitch_user_id,
		$2::uuid AS kick_user_id,
		$3::boolean AS twitch_bot_enabled,
		$4::boolean AS kick_bot_enabled,
		$5::text AS bot_id,
		$6::uuid AS kick_bot_id
),
created_channel AS (
	INSERT INTO channels (
		twitch_user_id,
		kick_user_id,
		twitch_bot_enabled,
		kick_bot_enabled,
		"isEnabled",
		"botId",
		kick_bot_id
	)
	SELECT
		r.twitch_user_id,
		r.kick_user_id,
		r.twitch_bot_enabled,
		r.kick_bot_enabled,
		r.twitch_bot_enabled OR r.kick_bot_enabled,
		r.bot_id,
		r.kick_bot_id
	FROM requested r
	WHERE (r.twitch_user_id IS NOT NULL OR r.twitch_bot_enabled = false)
		AND (
			r.kick_user_id IS NOT NULL
			OR (r.kick_bot_enabled = false AND r.kick_bot_id IS NULL)
		)
		AND (
			r.twitch_user_id IS NULL
			OR EXISTS (
				SELECT 1
				FROM users u
				WHERE u.id = r.twitch_user_id AND u.platform = 'twitch'
			)
		)
		AND (
			r.kick_user_id IS NULL
			OR EXISTS (
				SELECT 1
				FROM users u
				WHERE u.id = r.kick_user_id AND u.platform = 'kick'
			)
		)
	RETURNING
		id,
		twitch_user_id,
		twitch_bot_enabled,
		kick_user_id,
		kick_bot_enabled,
		"botId",
		"isBotMod",
		"isTwitchBanned",
		kick_bot_id
),
twitch_binding AS (
	INSERT INTO channel_platforms (
		channel_id,
		platform,
		user_id,
		platform_channel_id,
		enabled,
		bot_user_id,
		bot_config
	)
	SELECT
		c.id,
		'twitch',
		c.twitch_user_id,
		u.platform_id,
		c.twitch_bot_enabled,
		NULL,
		jsonb_build_object(
			'bot_id', c."botId",
			'is_bot_mod', c."isBotMod",
			'is_twitch_banned', c."isTwitchBanned"
		)
	FROM created_channel c
	JOIN users u ON u.id = c.twitch_user_id AND u.platform = 'twitch'
	WHERE c.twitch_user_id IS NOT NULL
	RETURNING channel_id
),
kick_binding AS (
	INSERT INTO channel_platforms (
		channel_id,
		platform,
		user_id,
		platform_channel_id,
		enabled,
		bot_user_id,
		bot_config
	)
	SELECT
		c.id,
		'kick',
		c.kick_user_id,
		u.platform_id,
		c.kick_bot_enabled,
		kb.kick_user_id,
		jsonb_strip_nulls(jsonb_build_object('kick_bot_id', c.kick_bot_id))
	FROM created_channel c
	JOIN users u ON u.id = c.kick_user_id AND u.platform = 'kick'
	LEFT JOIN kick_bots kb ON kb.id = c.kick_bot_id
	WHERE c.kick_user_id IS NOT NULL
	RETURNING channel_id
)
SELECT id FROM created_channel`

const updateChannelAndBindingsQuery = `
WITH input AS (
	SELECT
		$1::uuid AS channel_id,
		$2::boolean AS is_enabled,
		$3::boolean AS is_bot_mod,
		$4::uuid AS twitch_user_id,
		$5::uuid AS kick_user_id,
		$6::boolean AS twitch_bot_enabled,
		$7::boolean AS kick_bot_enabled,
		$8::uuid AS kick_bot_id
),
updated_channel AS (
	UPDATE channels c
	SET
		"isEnabled" = COALESCE(i.is_enabled, c."isEnabled"),
		"isBotMod" = COALESCE(i.is_bot_mod, c."isBotMod"),
		twitch_user_id = COALESCE(i.twitch_user_id, c.twitch_user_id),
		kick_user_id = COALESCE(i.kick_user_id, c.kick_user_id),
		twitch_bot_enabled = COALESCE(i.twitch_bot_enabled, c.twitch_bot_enabled),
		kick_bot_enabled = COALESCE(i.kick_bot_enabled, c.kick_bot_enabled),
		kick_bot_id = COALESCE(i.kick_bot_id, c.kick_bot_id)
	FROM input i
	WHERE c.id = i.channel_id
		AND (
			COALESCE(i.twitch_user_id, c.twitch_user_id) IS NULL
			OR EXISTS (
				SELECT 1
				FROM users u
				WHERE u.id = COALESCE(i.twitch_user_id, c.twitch_user_id)
					AND u.platform = 'twitch'
			)
		)
		AND (
			COALESCE(i.kick_user_id, c.kick_user_id) IS NULL
			OR EXISTS (
				SELECT 1
				FROM users u
				WHERE u.id = COALESCE(i.kick_user_id, c.kick_user_id)
					AND u.platform = 'kick'
			)
		)
		AND (
			COALESCE(i.twitch_user_id, c.twitch_user_id) IS NOT NULL
			OR (
				COALESCE(i.twitch_bot_enabled, c.twitch_bot_enabled) = false
				AND COALESCE(i.is_bot_mod, c."isBotMod") = false
				AND c."isTwitchBanned" = false
			)
		)
		AND (
			COALESCE(i.kick_user_id, c.kick_user_id) IS NOT NULL
			OR (
				COALESCE(i.kick_bot_enabled, c.kick_bot_enabled) = false
				AND COALESCE(i.kick_bot_id, c.kick_bot_id) IS NULL
			)
		)
	RETURNING
		c.id,
		c.twitch_user_id,
		c.twitch_bot_enabled,
		c.kick_user_id,
		c.kick_bot_enabled,
		c."botId",
		c."isBotMod",
		c."isTwitchBanned",
		c.kick_bot_id
),
twitch_binding AS (
	INSERT INTO channel_platforms (
		channel_id,
		platform,
		user_id,
		platform_channel_id,
		enabled,
		bot_user_id,
		bot_config
	)
	SELECT
		c.id,
		'twitch',
		c.twitch_user_id,
		u.platform_id,
		c.twitch_bot_enabled,
		NULL,
		jsonb_build_object(
			'bot_id', c."botId",
			'is_bot_mod', c."isBotMod",
			'is_twitch_banned', c."isTwitchBanned"
		)
	FROM updated_channel c
	JOIN users u ON u.id = c.twitch_user_id AND u.platform = 'twitch'
	WHERE c.twitch_user_id IS NOT NULL
	ON CONFLICT (channel_id, platform) DO UPDATE
	SET
		user_id = EXCLUDED.user_id,
		platform_channel_id = EXCLUDED.platform_channel_id,
		enabled = EXCLUDED.enabled,
		bot_user_id = EXCLUDED.bot_user_id,
		bot_config = EXCLUDED.bot_config,
		updated_at = NOW()
	RETURNING channel_id
),
kick_binding AS (
	INSERT INTO channel_platforms (
		channel_id,
		platform,
		user_id,
		platform_channel_id,
		enabled,
		bot_user_id,
		bot_config
	)
	SELECT
		c.id,
		'kick',
		c.kick_user_id,
		u.platform_id,
		c.kick_bot_enabled,
		kb.kick_user_id,
		jsonb_strip_nulls(jsonb_build_object('kick_bot_id', c.kick_bot_id))
	FROM updated_channel c
	JOIN users u ON u.id = c.kick_user_id AND u.platform = 'kick'
	LEFT JOIN kick_bots kb ON kb.id = c.kick_bot_id
	WHERE c.kick_user_id IS NOT NULL
	ON CONFLICT (channel_id, platform) DO UPDATE
	SET
		user_id = EXCLUDED.user_id,
		platform_channel_id = EXCLUDED.platform_channel_id,
		enabled = EXCLUDED.enabled,
		bot_user_id = EXCLUDED.bot_user_id,
		bot_config = EXCLUDED.bot_config,
		updated_at = NOW()
	RETURNING channel_id
)
SELECT id FROM updated_channel`

func (c *Pgx) GetByApiKey(ctx context.Context, apiKey string) (model.Channel, error) {
	return c.getOne(ctx, selectQuery+` WHERE c.api_key = $1`, apiKey)
}

func (c *Pgx) Create(ctx context.Context, input channels.CreateInput) (model.Channel, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	row := conn.QueryRow(
		ctx,
		createChannelAndBindingsQuery,
		valueOrNil(input.TwitchUserID),
		valueOrNil(input.KickUserID),
		input.TwitchBotEnabled,
		input.KickBotEnabled,
		input.BotID,
		valueOrNil(input.KickBotID),
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

func (c *Pgx) GetAllByBindingPlatform(
	ctx context.Context,
	p platform.Platform,
) ([]model.Channel, error) {
	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, getAllByBindingPlatformQuery, p)
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
	row := c.getter.DefaultTrOrDB(ctx, c.pool).QueryRow(
		ctx,
		updateChannelAndBindingsQuery,
		channelID,
		valueOrNil(input.IsEnabled),
		valueOrNil(input.IsBotMod),
		valueOrNil(input.TwitchUserID),
		valueOrNil(input.KickUserID),
		valueOrNil(input.TwitchBotEnabled),
		valueOrNil(input.KickBotEnabled),
		valueOrNil(input.KickBotID),
	)

	var channelId uuid.UUID
	if err := row.Scan(&channelId); err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, channelId)
}

func valueOrNil[T any](value *T) any {
	if value == nil {
		return nil
	}

	return *value
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

func (c *Pgx) getOne(ctx context.Context, query string, args ...any) (model.Channel, error) {
	rows, err := c.getter.DefaultTrOrDB(ctx, c.pool).Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	channel, err := collectExactlyOneChannel(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
		return model.Nil, err
	}

	return channel, nil
}
