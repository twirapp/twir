package seventv

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
)

const emoteForDeleteArgName = "name"

var EmoteDelete = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv delete",
		Description: null.StringFrom("Delete 7tv emote in current set"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "7tv",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{"7tv remove"},
		Enabled:     false,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: emoteForDeleteArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		client := seventvintegration.NewClient(parseCtx.Services.Config.SevenTvToken)

		sevenTvUser, err := client.GetProfileByTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.ProfileFailedToGet.
						SetVars(locales.KeysCommandsSeventvErrorsProfileFailedToGetVars{Reason: err.Error()}),
				),
				Err: err,
			}
		}
		if sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
			return &types.CommandsHandlerResult{
				Result: []string{
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Seventv.Errors.EmotesetNotActive,
					),
				},
			}, nil
		}

		name := parseCtx.ArgsParser.Get(emoteForDeleteArgName).String()
		var (
			foundEmoteId   string
			foundEmoteName string
		)
		for _, emote := range sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSet.Emotes.Items {
			if emote.Alias == name {
				foundEmoteId = emote.Emote.Id
				foundEmoteName = emote.Alias
				break
			}

			if emote.Emote.DefaultName == name {
				foundEmoteId = emote.Emote.Id
				foundEmoteName = emote.Emote.DefaultName
				break
			}
		}
		if foundEmoteId == "" || foundEmoteName == "" {
			return &types.CommandsHandlerResult{
				Result: []string{
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Seventv.Errors.EmoteNotFoundInEmoteset.
							SetVars(locales.KeysCommandsSeventvErrorsEmoteNotFoundInEmotesetVars{EmoteName: name, EmoteSet: sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSet.Name}),
					),
				},
			}, nil
		}

		err = client.RemoveEmote(
			ctx,
			*sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSetId,
			foundEmoteName,
			foundEmoteId,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmoteFailedToRemove,
				),
				Err: err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Remove.EmoteRemove.
						SetVars(locales.KeysCommandsSeventvRemoveEmoteRemoveVars{EmoteName: name}),
				),
			},
		}, nil
	},
}
