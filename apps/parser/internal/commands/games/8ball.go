package games

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		entity := model.ChannelModulesSettings{}
		if err := parseCtx.Services.Gorm.WithContext(ctx).Where(
			`"channelId" = ? and "userId" is null and "type" = '8ball'`,
			parseCtx.Channel.ID,
		).First(&entity).Error; err != nil {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}
		}

		var parsedSettings model.EightBallSettings
		if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}
		}

		if !parsedSettings.Enabled {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}
		}

		if len(parsedSettings.Answers) == 0 {
			return &types.CommandsHandlerResult{
				Result: []string{
					"I cannot answer to your question, " +
						"because 8ball not configured properly (missed answers)",
				},
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{lo.Sample(parsedSettings.Answers)},
		}
	},
}
