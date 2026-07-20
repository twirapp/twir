package pgx

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
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

const channelOwnershipJoinCondition = `(u.platform = 'twitch' AND uc.twitch_user_id = u.id) OR (u.platform = 'kick' AND uc.kick_user_id = u.id)`

func (c *Pgx) scanUserWithChannel(row pgx.Row) (model.UserWithChannel, error) {
	userWithChannel := model.UserWithChannel{}
	var channelID sql.Null[uuid.UUID]
	var channelTwitchUserID, channelKickUserID sql.Null[uuid.UUID]
	var channelBotID sql.Null[string]
	var channelIsBotMod, channelIsEnabled, channelIsTwitchBanned sql.Null[bool]

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
		&channelTwitchUserID,
		&channelKickUserID,
		&channelIsBotMod,
		&channelIsEnabled,
		&channelIsTwitchBanned,
		&channelBotID,
	)
	if err != nil {
		return model.Nil, err
	}

	if channelID.Valid {
		ch := &channelmodel.Channel{
			ID:             channelID.V,
			IsBotMod:       channelIsBotMod.V,
			IsEnabled:      channelIsEnabled.V,
			IsTwitchBanned: channelIsTwitchBanned.V,
			BotID:          channelBotID.V,
		}
		if channelTwitchUserID.Valid {
			ch.TwitchUserID = &channelTwitchUserID.V
		}
		if channelKickUserID.Valid {
			ch.KickUserID = &channelKickUserID.V
		}
		userWithChannel.Channel = ch
	}

	return userWithChannel, err
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.UserWithChannel, error) {
	query := `
SELECT u.id, u.platform, u.platform_id, u.login, u.display_name, u.avatar, u."isBotAdmin", u."tokenId", u."apiKey", u.hide_on_landing_page, u.is_banned,
       uc.id, uc.twitch_user_id, uc.kick_user_id, uc."isBotMod", uc."isEnabled", uc."isTwitchBanned", uc."botId"
FROM users u
LEFT JOIN channels uc ON ` + channelOwnershipJoinCondition + `
WHERE u.id = $1::uuid
GROUP BY u.id, uc.id
LIMIT 1
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	user, err := c.scanUserWithChannel(rows)
	if err != nil {
		return model.Nil, err
	}

	return user, nil
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

	if input.ChannelEnabled != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{`uc."isEnabled"`: input.ChannelEnabled})
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
		Select(
			"u.id",
			"u.platform",
			"u.platform_id",
			"u.login",
			"u.display_name",
			"u.avatar",
			`u."isBotAdmin"`,
			`u."tokenId"`,
			`u."apiKey"`,
			"u.hide_on_landing_page",
			"u.is_banned",
			"uc.id AS channel_id",
			"uc.twitch_user_id AS channel_twitch_user_id",
			"uc.kick_user_id AS channel_kick_user_id",
			`uc."isBotMod" AS channel_is_bot_mod`,
			`uc."isEnabled" AS channel_is_enabled`,
			`uc."isTwitchBanned" AS channel_is_twitch_banned`,
			`uc."botId" AS channel_bot_id`,
		).
		From("users u").
		LeftJoin("channels uc ON " + channelOwnershipJoinCondition).
		OrderBy("u.id asc")

	if len(input.HasBadgesIDS) > 0 {
		selectQuery = selectQuery.LeftJoin("badges_users bu ON u.id = bu.user_id")
	}

	selectQuery = applyFilters(selectQuery, input)

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
		selectQuery = selectQuery.LeftJoin("channels uc ON " + channelOwnershipJoinCondition)
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
