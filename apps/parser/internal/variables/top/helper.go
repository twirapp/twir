package top

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/channelbinding"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/entities/platform"
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
	UserID            string `db:"user_id"`
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
	if err != nil && parseCtx.Channel.TwitchUserID != uuid.Nil {
		channel, chanErr := parseCtx.Services.ChannelService.GetChannelByBindingUserID(ctx, platform.PlatformTwitch, parseCtx.Channel.TwitchUserID)
		if chanErr == nil && !channel.IsNil() {
			channelID = channel.ID
		}
	}
	if err != nil && channelID == uuid.Nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	channel, err := parseCtx.Services.ChannelService.GetChannelByID(ctx, channelID)
	if err != nil || channel.IsNil() {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	qb := squirrel.
		Select(`users_stats.user_id`, `users_stats.messages`, `users_stats.watched`, `users_stats."usedChannelPoints"`).
		From("users_stats").
		Where(
			squirrel.And{
				squirrel.Eq{`users_stats.channel_id`: parseCtx.Channel.DBChannelID},
				squirrel.Gt{`users_stats.messages`: 0},
			},
		).
		Where(`NOT EXISTS (SELECT 1 FROM users_ignored ui JOIN users u ON u.platform = 'twitch' AND u.platform_id = ui.id WHERE u.id = users_stats.user_id)`).
		Where(`NOT EXISTS (SELECT 1 FROM bots b JOIN users u ON u.platform_id = b.id AND u.platform = 'twitch' WHERE u.id = users_stats.user_id)`)

	qb, err = applyTopChannelBotFilters(qb, channel)
	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	query, args, err := qb.
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy(fmt.Sprintf(`users_stats.%s DESC`, topType)).
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

	ids := lo.FilterMap(records, func(record userStatsPlatform, _ int) (uuid.UUID, bool) {
		parsedID, err := uuid.Parse(record.UserID)
		if err != nil {
			return uuid.Nil, false
		}

		return parsedID, true
	})

	users, err := parseCtx.Services.UsersRepo.GetManyByIDS(ctx, usersrepository.GetManyInput{IDs: ids, PerPage: len(ids)})
	if err != nil {
		parseCtx.Services.Logger.Sugar().Error(err)
		return nil, false
	}

	usersByID := lo.SliceToMap(users, func(user usersmodel.User) (string, usersmodel.User) {
		return user.ID.String(), user
	})
	twitchLoginsByPlatformID := make(map[string]string)
	if parseCtx.Platform == platform.PlatformTwitch && parseCtx.Services.CacheTwitchClient != nil {
		platformIDs := lo.Map(users, func(user usersmodel.User, _ int) string {
			return user.PlatformID
		})
		twitchUsers, err := parseCtx.Services.CacheTwitchClient.GetUsersByIds(ctx, platformIDs)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
		} else {
			for _, twitchUser := range twitchUsers {
				twitchLoginsByPlatformID[twitchUser.ID] = twitchUser.Login
			}
		}
	}

	var stats []*userStats
	for _, record := range records {
		user, ok := usersByID[record.UserID]
		if !ok {
			continue
		}

		userName := user.Login
		if twitchLogin, ok := twitchLoginsByPlatformID[user.PlatformID]; ok {
			userName = twitchLogin
		}
		if userName == "" {
			userName = user.DisplayName
		}
		if userName == "" {
			userName = user.PlatformID
		}

		res := &userStats{
			DisplayName: user.DisplayName,
			UserName:    userName,
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

func applyTopChannelBotFilters(qb squirrel.SelectBuilder, channel channelmodel.Channel) (squirrel.SelectBuilder, error) {
	_, twitchBotConfig, hasTwitchBinding, err := channelbinding.FindTwitch(channel)
	if err != nil {
		return qb, err
	}
	if hasTwitchBinding && twitchBotConfig.BotID != "" {
		qb = qb.Where(
			squirrel.Expr(
				`users_stats.user_id NOT IN (SELECT id FROM users WHERE platform_id = ? AND platform = 'twitch')`,
				twitchBotConfig.BotID,
			),
		)
	}

	for _, binding := range channel.Bindings {
		if binding.BotUserID != nil {
			qb = qb.Where(squirrel.NotEq{`users_stats.user_id`: binding.BotUserID.String()})
		}
		if binding.UserID != uuid.Nil {
			qb = qb.Where(squirrel.NotEq{`users_stats.user_id`: binding.UserID.String()})
		}
	}

	return qb, nil
}
