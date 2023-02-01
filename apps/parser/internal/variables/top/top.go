package top

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

type UserStats struct {
	DisplayName string
	UserName    string
	Value       int
}

func GetTop(
	ctx *variables_cache.VariablesCacheService,
	channelId string,
	topType string,
	page *int,
) []*UserStats {
	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	db := do.MustInvoke[gorm.DB](di.Provider)

	if page == nil {
		newPage := 1
		page = &newPage
	}

	offset := (*page - 1) * 10

	var records []model.UsersStats

	err = db.
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

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
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

		if topType == "watched" {
			res.Value = int(record.Watched)
		}

		return res
	})

	return stats
}
