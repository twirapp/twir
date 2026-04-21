package top

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

type userStats struct {
	DisplayName string
	UserName    string
	Value       int
}

type userStatsPlatform struct {
	model.UsersStats
	PlatformID string `db:"platform_id"`
}

func getTop(
	ctx context.Context,
	parseCtx *types.VariableParseContext,
	topType string,
	page *int,
	limit int,
) ([]*userStats, bool) {
	if parseCtx.Platform != "twitch" {
		return nil, true
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*parseCtx.Services.Config,
		parseCtx.Services.Bus,
	)

	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	if page == nil {
		newPage := 1
		page = &newPage
	}

	offset := (*page - 1) * limit

	channel := &model.Channels{}
	err = parseCtx.Services.Gorm.
		WithContext(ctx).
		Where(`"id" = ?`, parseCtx.Channel.ID).
		Find(channel).Error
	if err != nil || channel.ID == "" {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	qb := squirrel.
		Select(`"users_stats".*`, `"users"."platform_id"`).
		From("users_stats").
		Join(`users ON users.id = "users_stats"."userId"`).
		Where(
			squirrel.And{
				squirrel.Eq{`"users_stats"."channelId"`: parseCtx.Channel.ID},
				squirrel.Gt{`"users_stats"."messages"`: 0},
			},
		).
		Where(`NOT EXISTS (SELECT 1 FROM users_ignored ui JOIN users u ON u.platform_id = ui.id WHERE u.id = "users_stats"."userId")`).
		Where(`NOT EXISTS (SELECT 1 FROM bots b JOIN users u ON u.platform_id = b.id AND u.platform = 'twitch' WHERE u.id = "users_stats"."userId")`)

	if channel.BotID != "" {
		qb = qb.Where(
			squirrel.Expr(
				`"users_stats"."userId" NOT IN (SELECT id FROM users WHERE platform_id = ? AND platform = 'twitch')`,
				channel.BotID,
			),
		)
	}
	if channel.TwitchUserID != nil {
		qb = qb.Where(squirrel.NotEq{`"users_stats"."userId"`: *channel.TwitchUserID})
	}

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

	ids := lo.Map(
		records, func(record userStatsPlatform, _ int) string {
			return record.PlatformID
		},
	)

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: ids,
		},
	)

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil, false
	}

	var stats []*userStats
	for _, record := range records {
		twitchUser, ok := lo.Find(
			twitchUsers.Data.Users, func(user helix.User) bool {
				return user.ID == record.PlatformID
			},
		)

		if !ok {
			continue
		}

		res := &userStats{
			DisplayName: twitchUser.DisplayName,
			UserName:    twitchUser.Login,
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
