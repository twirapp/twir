package games

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

const (
	eightBallArgName = "question"
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
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: eightBallArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		entity := model.ChannelGames8Ball{}
		if err := parseCtx.Services.Gorm.WithContext(ctx).Where(
			`"channel_id" = ?`,
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

		if !entity.Enabled {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		if len(entity.Answers) == 0 {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		return &types.CommandsHandlerResult{
			Result: []string{lo.Sample(entity.Answers)},
		}, nil
	},
}
