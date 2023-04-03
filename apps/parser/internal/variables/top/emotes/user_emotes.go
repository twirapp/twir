package emotes

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type topUser struct {
	Name  string
	Count int
}

var UsersVariable = types.Variable{
	Name:        "top.emotes.users",
	Description: lo.ToPtr("Top users by emotes. Prints counts of used emotes"),
	Example:     lo.ToPtr("top.emotes.users|10"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		cfg := do.MustInvoke[config.Config](di.Provider)
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
		twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)

		limit := 10
		if data.Params != nil {
			newLimit, err := strconv.Atoi(*data.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		usages := []model.ChannelEmoteUsage{}
		err = db.Model(&model.ChannelEmoteUsage{}).
			Select(`"userId", COUNT(*) as count`).
			Group(`"userId"`).
			Order("count DESC").
			Limit(10).
			Scan(&usages).
			Error

		if err != nil {
			return nil, err
		}

		twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
			IDs: lo.Map(usages, func(item model.ChannelEmoteUsage, _ int) string {
				return item.UserID
			}),
		})

		if err != nil {
			return nil, err
		}

		mappedTop := []string{}

		for _, usage := range usages {
			user, ok := lo.Find(twitchUsers.Data.Users, func(item helix.User) bool {
				return item.ID == usage.UserID
			})

			if !ok {
				continue
			}

			mappedTop = append(mappedTop, fmt.Sprintf(
				"%s × %v",
				user.Login,
				lo.If(usage.Count != nil, *usage.Count).Else(0),
			))
		}

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
