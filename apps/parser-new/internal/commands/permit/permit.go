package permit

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"strconv"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/nicklaw5/helix/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "permit",
		Description: null.StringFrom("Permits user."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)

		result := &types.CommandsHandlerResult{}

		count := 1
		params := strings.Split(*parseCtx.Text, " ")

		paramsLen := len(params)
		if paramsLen < 1 {
			result.Result = []string{"you have type user name to permit."}
			return result
		}

		if paramsLen == 2 {
			newCount, err := strconv.Atoi(params[1])
			if err == nil {
				count = newCount
			}
		}

		if count > 100 {
			result.Result = []string{"cannot create more then 100 permits."}
			return result
		}

		target, err := twitchClient.GetUsers(&helix.UsersParams{
			Logins: []string{params[0]},
		})

		if err != nil || target.StatusCode != 200 || len(target.Data.Users) == 0 {
			result.Result = []string{"user not found."}
			return result
		}

		parseCtx.Services.Gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for i := 0; i < count; i++ {
				permit := model.ChannelsPermits{
					ID:        uuid.NewV4().String(),
					ChannelID: parseCtx.Channel.ID,
					UserID:    target.Data.Users[0].ID,
				}
				err := tx.Create(&permit).Error
				if err != nil {
					zap.S().Error(err)
					return err
				}
			}
			return nil
		})

		result.Result = []string{
			fmt.Sprintf("✅ added %v permits to %s", count, target.Data.Users[0].DisplayName),
		}
		return result
	},
}
