package top

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type userStats struct {
	DisplayName string
	UserName    string
	Value       int
}

type userStatsPlatform struct {
	UserID            string `db:"userId"`
	Messages          int64  `db:"messages"`
	Watched           int64  `db:"watched"`
	UsedChannelPoints int64  `db:"usedChannelPoints"`
}

func getTop(
	ctx context.Context,
	parseCtx *types.VariableParseContext,
	topType string,
	page *int,
	limit int,
) ([]*userStats, bool) {
	if page == nil {
		newPage := 1
		page = &newPage
	}

	offset := (*page - 1) * limit

	channelID, err := uuid.Parse(parseCtx.Channel.DBChannelID)
	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	channel, err := parseCtx.Services.ChannelsRepo.GetByID(ctx, channelID)
	if err != nil || channel.IsNil() {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	qb := squirrel.
		Select(`"users_stats"."userId"`, `"users_stats"."messages"`, `"users_stats"."watched"`, `"users_stats"."usedChannelPoints"`).
		From("users_stats").
		Where(
			squirrel.And{
				squirrel.Eq{`"users_stats"."channelId"`: parseCtx.Channel.DBChannelID},
				squirrel.Gt{`"users_stats"."messages"`: 0},
			},
		).
		Where(`NOT EXISTS (SELECT 1 FROM users_ignored ui JOIN users u ON u.platform_id = ui.id WHERE u.id = "users_stats"."userId")`).
		Where(`NOT EXISTS (SELECT 1 FROM bots b JOIN users u ON u.platform_id = b.id AND u.platform = 'twitch' WHERE u.id = "users_stats"."userId")`)

	qb = applyTopChannelBotFilters(qb, channel)

	query, args, err := qb.
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy(fmt.Sprintf(`"users_stats"."%s" DESC`, topType)).
		ToSql()
	query = parseCtx.Services.Sqlx.Rebind(query)

	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	records := []userStatsPlatform{}

	err = parseCtx.Services.Sqlx.Select(&records, query, args...)
	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	if len(records) == 0 {
		return nil, false
	}

	ids := lo.Map(records, func(record userStatsPlatform, _ int) string {
		return record.UserID
	})

	users, err := parseCtx.Services.UsersRepo.GetManyByIDS(ctx, usersrepository.GetManyInput{IDs: ids, PerPage: len(ids)})
	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	usersByID := lo.SliceToMap(users, func(user usersmodel.User) (string, usersmodel.User) {
		return user.ID, user
	})

	var stats []*userStats
	for _, record := range records {
		user, ok := usersByID[record.UserID]
		if !ok {
			continue
		}

		res := &userStats{
			DisplayName: user.DisplayName,
			UserName:    user.Login,
		}

		if topType == "messages" {
			res.Value = int(record.Messages)
		}

		if topType == "watched" {
			res.Value = int(record.Watched)
		}

		if topType == "usedChannelPoints" {
			res.Value = int(record.UsedChannelPoints)
		}

		if res.Value == 0 {
			continue
		}

		stats = append(stats, res)
	}

	return stats, false
}

func applyTopChannelBotFilters(qb squirrel.SelectBuilder, channel channelmodel.Channel) squirrel.SelectBuilder {
	if channel.BotID != "" {
		qb = qb.Where(
			squirrel.Expr(
				`"users_stats"."userId" NOT IN (SELECT id FROM users WHERE platform_id = ? AND platform = 'twitch')`,
				channel.BotID,
			),
		)
	}

	if channel.KickBotID != nil {
		qb = qb.Where(squirrel.NotEq{`"users_stats"."userId"`: channel.KickBotID.String()})
	}

	if channel.TwitchUserID != nil {
		qb = qb.Where(squirrel.NotEq{`"users_stats"."userId"`: *channel.TwitchUserID})
	}

	if channel.KickUserID != nil {
		qb = qb.Where(squirrel.NotEq{`"users_stats"."userId"`: channel.KickUserID.String()})
	}

	return qb
}
