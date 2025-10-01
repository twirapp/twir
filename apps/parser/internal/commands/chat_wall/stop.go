package chat_wall

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	chatwallservice "github.com/twirapp/twir/apps/parser/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
					Result: []string{i18n.GetCtx(
						ctx,
						locales.Translations.Commands.ChatWall.Errors.ChatWallNotFound.
							SetVars(locales.KeysCommandsChatWallErrorsChatWallNotFoundVars{ErrorPhrase: phrase}),
					)},
				}, nil
			}

			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{i18n.GetCtx(
				ctx,
				locales.Translations.Commands.ChatWall.Stop.ChatWalStop.
					SetVars(locales.KeysCommandsChatWallStopChatWalStopVars{ChatWallPhrase: phrase}))},
		}

		return result, nil
	},
}
