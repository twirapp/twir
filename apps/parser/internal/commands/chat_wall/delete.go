package chat_wall

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	chatwallservice "github.com/twirapp/twir/apps/parser/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
)

const deletePhraseArgName = "phrase"

var Delete = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "chat wall delete",
		Description: null.StringFrom("Creates chat wall instance for 10 minutes, which will delete all messages with specific phrase"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "MODERATION",
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name: deletePhraseArgName,
			Hint: "phrase to delete",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		phrase := parseCtx.ArgsParser.Get(deletePhraseArgName).String()

		wall, err := parseCtx.Services.ChatWallService.Create(
			ctx,
			chatwallservice.CreateInput{
				ChannelID: parseCtx.Channel.ID,
				Phrase:    phrase,
				Enabled:   true,
				Action:    chatwallmodel.ChatWallActionDelete,
				Duration:  10 * time.Minute,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		if err := parseCtx.Services.ChatWallService.HandlePastMessages(
			ctx,
			wall,
			chatwallservice.HandlePastMessagesInput{
				ChannelID:       parseCtx.Channel.ID,
				Phrase:          phrase,
				Action:          chatwallmodel.ChatWallActionDelete,
				TimeoutDuration: nil,
			},
		); err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					`✅ Chat wall started for 10 minutes, you can stop it with !chat wall stop "%s"`,
					phrase,
				),
			},
		}

		return result, nil
	},
}
