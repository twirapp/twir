package pgx

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (c *Pgx) scanUserWithChannel(row pgx.Row) (model.UserWithChannel, error) {
	userWithChannel := model.UserWithChannel{}
	var channelID, channelBotID sql.Null[string]
	var channelIsBotMod, channelIsEnabled, channelIsTwitchBanned sql.Null[bool]

	err := row.Scan(
		&userWithChannel.User.ID,
		&userWithChannel.User.IsTester,
		&userWithChannel.User.IsBotAdmin,
		&userWithChannel.User.TokenID,
		&userWithChannel.User.ApiKey,
		&userWithChannel.User.HideOnLandingPage,
		&userWithChannel.User.IsBanned,
		&channelID,
		&channelIsBotMod,
		&channelIsEnabled,
		&channelIsTwitchBanned,
		&channelBotID,
	)
	if err != nil {
		return model.Nil, err
	}

	if channelID.Valid {
		userWithChannel.Channel = &channelmodel.Channel{
			ID:             channelID.V,
			IsBotMod:       channelIsBotMod.V,
			IsEnabled:      channelIsEnabled.V,
			IsTwitchBanned: channelIsTwitchBanned.V,
			BotID:          channelBotID.V,
		}
	}

	return userWithChannel, err
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.UserWithChannel, error) {
	query := `
SELECT u.id, u."isTester", u."isBotAdmin", u."tokenId", u."apiKey", u.hide_on_landing_page, u.is_banned,
			 uc.id, uc."isBotMod", uc."isEnabled", uc."isTwitchBanned", uc."botId"
FROM users u
LEFT JOIN channels uc ON u.id = uc.id
WHERE u.id = $1
GROUP BY u.id
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

func (c *Pgx) GetManyByIDS(
	ctx context.Context,
	input users_with_channel.GetManyInput,
) ([]model.UserWithChannel, error) {
	selectQuery := sq.
		Select(
			"u.id",
			`u."isTester"`,
			`u."isBotAdmin"`,
			`u."tokenId"`,
			`u."apiKey"`,
			"u.hide_on_landing_page",
			"u.is_banned",
			"uc.id",
			`uc."isBotMod"`,
			`uc."isEnabled"`,
			`uc."isTwitchBanned"`,
			`uc."botId"`,
		).
		From("users u").
		LeftJoin("channels uc ON u.id = uc.id").
		LeftJoin("badges_users bu ON u.id = bu.user_id").
		GroupBy("u.id", "uc.id").
		OrderBy("u.id asc")

	if len(input.IDs) > 0 {
		selectQuery = selectQuery.Where(squirrel.Eq{"u.id": input.IDs})
	}

	if len(input.HasBadgesIDS) > 0 {
		selectQuery = selectQuery.Where(squirrel.Eq{"bu.badge_id": input.HasBadgesIDS})
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

	if input.ChannelEnabled != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{`uc."isEnabled"`: input.ChannelEnabled})
	}

	if input.ChannelIsBotAdmin != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{`u."isBotAdmin"`: input.ChannelIsBotAdmin})
	}

	if input.IsBanned != nil {
		selectQuery = selectQuery.Where(squirrel.Eq{`u.is_banned`: input.IsBanned})
	}

	query, args, err := selectQuery.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	users := make([]model.UserWithChannel, 0)

	for rows.Next() {
		user, err := c.scanUserWithChannel(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (c *Pgx) GetManyCount(ctx context.Context, input users_with_channel.GetManyInput) (
	int,
	error,
) {
	selectQuery := sq.
		Select("COUNT(*)").
		From("users u")

	if len(input.IDs) > 0 {
		selectQuery = selectQuery.Where(squirrel.Eq{"u.id": input.IDs})
	}

	if input.ChannelIsBotAdmin != nil || input.ChannelEnabled != nil {
		selectQuery = selectQuery.LeftJoin("channels uc ON u.id = uc.id")
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

	if len(input.HasBadgesIDS) > 0 {
		selectQuery = selectQuery.
			LeftJoin("channels uc ON u.id = uc.id").
			Where(squirrel.Eq{"bu.badge_id": input.HasBadgesIDS})
	}

	query, args, err := selectQuery.ToSql()
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
