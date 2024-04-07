package games

import (
	"context"
	"errors"

	"github.com/goccy/go-json"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

var EightBall = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "8ball",
		Description: null.StringFrom("Magic ball will answer to all your questions!"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		entity := model.ChannelModulesSettings{}
		if err := parseCtx.Services.Gorm.WithContext(ctx).Where(
			`"channelId" = ? and "userId" is null and "type" = '8ball'`,
			parseCtx.Channel.ID,
		).First(&entity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			return nil, &types.CommandHandlerError{
				Message: "cannot find 8ball settings",
				Err:     err,
			}
		}

		var parsedSettings model.EightBallSettings
		if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot parse 8ball settings",
				Err:     err,
			}
		}

		if !parsedSettings.Enabled {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		if len(parsedSettings.Answers) == 0 {
			return nil, &types.CommandHandlerError{
				Message: `I cannot answer to your question, because 8ball not configured properly (missed answers)`,
				Err:     nil,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{lo.Sample(parsedSettings.Answers)},
		}, nil
	},
}
