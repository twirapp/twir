package top

import (
	"fmt"
	model "tsuwari/parser/internal/models"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/nicklaw5/helix"
	"github.com/samber/lo"
)

type UserStats struct {
	DisplayName string
	UserName    string
	Value       int
}

func GetTop(ctx *variables_cache.VariablesCacheService, channelId string, topType string, page *int) *[]*UserStats {
	if page == nil {
		newPage := 1
		page = &newPage
	}

	offset := (*page - 1) * 10

	var records []model.UsersStats

	err := ctx.Services.Db.
		Where(`"channelId" = ?`, channelId).
		Limit(10).
		Offset(offset).
		Order(fmt.Sprintf("%s DESC", topType)).
		Find(&records).Error

	if err != nil {
		return nil
	}

	ids := lo.Map(records, func(record model.UsersStats, _ int) string {
		return record.UserID
	})

	twitchUsers, err := ctx.Services.Twitch.Client.GetUsers(&helix.UsersParams{
		IDs: ids,
	})

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil
	}

	stats := lo.Map(records, func(record model.UsersStats, _ int) *UserStats {
		twitchUser, ok := lo.Find(twitchUsers.Data.Users, func(user helix.User) bool {
			return user.ID == record.UserID
		})

		if !ok {
			return nil
		}

		res := &UserStats{
			DisplayName: twitchUser.DisplayName,
			UserName:    twitchUser.Login,
		}

		if topType == "messages" {
			res.Value = int(record.Messages)
		}

		return res
	})

	return &stats
}
