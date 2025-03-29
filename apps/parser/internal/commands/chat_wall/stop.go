package chat_wall

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	chatwallservice "github.com/satont/twir/apps/parser/internal/services/chat_wall"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

const stopPhraseArgName = "phrase"

var Stop = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "chat wall stop",
		Description: null.StringFrom("Stop chat wall instance"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "MODERATION",
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name: stopPhraseArgName,
			Hint: "phrase to stop",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		phrase := parseCtx.ArgsParser.Get(banPhraseArgName).String()

		err := parseCtx.Services.ChatWallService.Stop(
			ctx,
			chatwallservice.StopInput{
				ChannelID: parseCtx.Channel.ID,
				Phrase:    phrase,
			},
		)
		if err != nil {
			if errors.Is(err, chatwallservice.ErrChatWallNotFound) {
				return &types.CommandsHandlerResult{
					Result: []string{
						fmt.Sprintf(
							`Chat wall "%s" not found or already stopped`,
							phrase,
						),
					},
				}, nil
			}

			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					`✅ Chat wall "%s" stopped`,
					phrase,
				),
			},
		}

		return result, nil
	},
}
