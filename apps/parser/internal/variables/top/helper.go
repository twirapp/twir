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

func getTop(
	ctx context.Context,
	parseCtx *types.VariableParseContext,
	topType string,
	page *int,
	limit int,
) []*userStats {
	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*parseCtx.Services.Config,
		parseCtx.Services.Bus,
	)

	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil
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
		return nil
	}

	query, args, err := squirrel.
		Select("users_stats.*").
		From("users_stats").
		Where(
			squirrel.And{
				squirrel.Eq{`"channelId"`: parseCtx.Channel.ID},
				squirrel.NotEq{`"userId"`: channel.BotID},
				squirrel.NotEq{`"userId"`: parseCtx.Channel.ID},
				squirrel.Gt{"messages": 0},
			},
		).
		Where(`NOT EXISTS (select 1 from users_ignored where "id" = "users_stats"."userId")`).
		Where(`NOT EXISTS (select 1 from bots where "id" = "users_stats"."userId")`).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy(fmt.Sprintf(`"%s" DESC`, topType)).
		ToSql()
	query = parseCtx.Services.Sqlx.Rebind(query)

	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil
	}

	records := []model.UsersStats{}

	err = parseCtx.Services.Sqlx.Select(&records, query, args...)
	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil
	}

	ids := lo.Map(
		records, func(record model.UsersStats, _ int) string {
			return record.UserID
		},
	)

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: ids,
		},
	)

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil
	}

	var stats []*userStats
	for _, record := range records {
		twitchUser, ok := lo.Find(
			twitchUsers.Data.Users, func(user helix.User) bool {
				return user.ID == record.UserID
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

	return stats
}
