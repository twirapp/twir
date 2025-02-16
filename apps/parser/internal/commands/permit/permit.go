package permit

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	"go.uber.org/zap"

	model "github.com/satont/twir/libs/gomodels"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	userArgName  = "@nickname"
	countArgName = "count"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "permit",
		Description: null.StringFrom("Permits user."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: userArgName,
		},
		command_arguments.Int{
			Name:     countArgName,
			Optional: true,
			Min:      lo.ToPtr(1),
			Max:      lo.ToPtr(100),
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		count := 1

		countArg := parseCtx.ArgsParser.Get(countArgName)
		if countArg != nil {
			count = countArg.Int()
		}

		if len(parseCtx.Mentions) == 0 {
			result.Result = []string{"user not found."}
			return result, nil
		}

		user := parseCtx.Mentions[0]

		txErr := parseCtx.Services.Gorm.WithContext(ctx).Transaction(
			func(tx *gorm.DB) error {
				for i := 0; i < count; i++ {
					permit := model.ChannelsPermits{
						ID:        uuid.NewV4().String(),
						ChannelID: parseCtx.Channel.ID,
						UserID:    user.UserId,
					}
					err := tx.Create(&permit).Error
					if err != nil {
						zap.S().Error(err)
						return err
					}
				}
				return nil
			},
		)

		if txErr != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create permit",
				Err:     txErr,
			}
		}

		result.Result = []string{
			fmt.Sprintf("✅ added %v permits to %s", count, user.UserName),
		}
		return result, nil
	},
}
