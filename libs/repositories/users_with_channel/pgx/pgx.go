package pgx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/repositories/users_with_channel"
	"github.com/twirapp/twir/libs/repositories/users_with_channel/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ users_with_channel.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

const userWithChannelColumns = `
	u.id,
	u.platform,
	u.platform_id,
	u.login,
	u.display_name,
	u.avatar,
	u."isBotAdmin",
	u."tokenId",
	u."apiKey",
	u.hide_on_landing_page,
	u.is_banned,
	cb.channel_id,
	CASE WHEN cb.id IS NULL THEN NULL ELSE jsonb_build_object(
		'ID', cb.id,
		'ChannelID', cb.channel_id,
		'Platform', cb.platform,
		'UserID', cb.user_id,
		'PlatformChannelID', cb.platform_channel_id,
		'Enabled', cb.enabled,
		'BotUserID', cb.bot_user_id,
		'BotConfig', cb.bot_config,
		'CreatedAt', cb.created_at,
		'UpdatedAt', cb.updated_at
	) END AS channel_binding`

const ownerBindingLateral = `LATERAL (
	SELECT
		cp.id,
		cp.channel_id,
		cp.platform,
		cp.user_id,
		cp.platform_channel_id,
		cp.enabled,
		cp.bot_user_id,
		cp.bot_config,
		cp.created_at,
		cp.updated_at
	FROM channel_platforms cp
	WHERE cp.user_id = u.id
		AND cp.platform = u.platform
	ORDER BY cp.channel_id
	LIMIT 1
) cb ON TRUE`

const ownerBindingJoin = `LEFT JOIN ` + ownerBindingLateral

const getByIDQuery = `
SELECT ` + userWithChannelColumns + `
FROM users u
` + ownerBindingJoin + `
WHERE u.id = $1::uuid
LIMIT 1`

func (c *Pgx) scanUserWithChannel(row pgx.Row) (model.UserWithChannel, error) {
	userWithChannel := model.UserWithChannel{}
	var channelID pgtype.UUID
	var bindingJSON []byte

	err := row.Scan(
		&userWithChannel.User.ID,
		&userWithChannel.User.Platform,
		&userWithChannel.User.PlatformID,
		&userWithChannel.User.Login,
		&userWithChannel.User.DisplayName,
		&userWithChannel.User.Avatar,
		&userWithChannel.User.IsBotAdmin,
		&userWithChannel.User.TokenID,
		&userWithChannel.User.ApiKey,
		&userWithChannel.User.HideOnLandingPage,
		&userWithChannel.User.IsBanned,
		&channelID,
		&bindingJSON,
	)
	if err != nil {
		return model.Nil, err
	}

	return mapUserWithChannelProjection(userWithChannel.User, channelID, bindingJSON)
}

func mapUserWithChannelProjection(
	user usermodel.User,
	channelID pgtype.UUID,
	bindingJSON []byte,
) (model.UserWithChannel, error) {
	result := model.UserWithChannel{User: user}
	if !channelID.Valid {
		return result, nil
	}

	var binding channelplatformsmodel.ChannelPlatform
	if err := json.Unmarshal(bindingJSON, &binding); err != nil {
		return model.Nil, fmt.Errorf("unmarshal channel platform binding: %w", err)
	}

	result.Channel = &channelmodel.Channel{
		ID:       uuid.UUID(channelID.Bytes),
		Bindings: []channelplatformsmodel.ChannelPlatform{binding},
	}

	return result, nil
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.UserWithChannel, error) {
	return c.scanUserWithChannel(c.pool.QueryRow(ctx, getByIDQuery, id))
}

func platformsToStrings(platforms []platformentity.Platform) []string {
	result := make([]string, 0, len(platforms))
	for _, platform := range platforms {
		result = append(result, platform.String())
	}

	return result
}

func applyFilters(selectQuery squirrel.SelectBuilder, input users_with_channel.GetManyInput) squirrel.SelectBuilder {
	if len(input.IDs) > 0 {
		selectQuery = selectQuery.Where(squirrel.Eq{"u.platform_id": input.IDs})
	}

	if input.SearchQuery != "" {
		searchQuery := "%" + strings.TrimSpace(input.SearchQuery) + "%"
		selectQuery = selectQuery.Where(squirrel.Or{
			squirrel.Expr("u.login ILIKE ?", searchQuery),
			squirrel.Expr("u.display_name ILIKE ?", searchQuery),
			squirrel.Expr("u.platform_id ILIKE ?", searchQuery),
		})
	}

	if len(input.Platforms) > 0 {
		selectQuery = selectQuery.Where(squirrel.Eq{"u.platform": platformsToStrings(input.Platforms)})
	}

	if len(input.HasBadgesIDS) > 0 {
		selectQuery = selectQuery.Where(squirrel.Eq{"bu.badge_id": input.HasBadgesIDS})
	}

	if input.ChannelIsBotAdmin != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{`u."isBotAdmin"`: input.ChannelIsBotAdmin})
	}

	if input.IsBanned != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{`u.is_banned`: input.IsBanned})
	}

	return selectQuery
}

func buildGetManyQuery(input users_with_channel.GetManyInput) (string, []any, error) {
	selectQuery := sq.
		Select(userWithChannelColumns).
		From("users u").
		LeftJoin(ownerBindingLateral).
		OrderBy("u.id asc")

	if len(input.HasBadgesIDS) > 0 {
		selectQuery = selectQuery.LeftJoin("badges_users bu ON u.id = bu.user_id")
	}

	selectQuery = applyFilters(selectQuery, input)
	if input.ChannelEnabled != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{"cb.enabled": *input.ChannelEnabled})
	}

	var page int
	perPage := 20

	if input.Page > 0 {
		page = input.Page
	}

	if input.PerPage > 0 {
		perPage = input.PerPage
	}

	offset := page * perPage

	selectQuery = selectQuery.Limit(uint64(perPage)).Offset(uint64(offset))

	return selectQuery.ToSql()
}

func (c *Pgx) GetManyByIDS(
	ctx context.Context,
	input users_with_channel.GetManyInput,
) ([]model.UserWithChannel, error) {
	query, args, err := buildGetManyQuery(input)
	if err != nil {
		return nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.UserWithChannel, 0)

	for rows.Next() {
		user, err := c.scanUserWithChannel(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func buildGetManyCountQuery(input users_with_channel.GetManyInput) (string, []any, error) {
	countColumn := "COUNT(*)"
	if input.ChannelEnabled != nil || len(input.HasBadgesIDS) > 0 {
		countColumn = "COUNT(DISTINCT u.id)"
	}

	selectQuery := sq.Select(countColumn).From("users u")

	if input.ChannelEnabled != nil {
		selectQuery = selectQuery.
			LeftJoin(ownerBindingLateral).
			Where(squirrel.Eq{"cb.enabled": *input.ChannelEnabled})
	}

	if len(input.HasBadgesIDS) > 0 {
		selectQuery = selectQuery.
			LeftJoin("badges_users bu ON u.id = bu.user_id").
			Where(squirrel.Eq{"bu.badge_id": input.HasBadgesIDS})
	}

	selectQuery = applyFilters(selectQuery, input)

	return selectQuery.ToSql()
}

func (c *Pgx) GetManyCount(ctx context.Context, input users_with_channel.GetManyInput) (
	int,
	error,
) {
	query, args, err := buildGetManyCountQuery(input)
	if err != nil {
		return 0, err
	}

	var count int
	err = c.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
