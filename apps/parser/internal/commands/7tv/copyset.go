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
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
)

const (
	copySetChannelName = "channelName"
	copySetNameOfSet   = "nameOfSet"
)

// request target emote set from channel B -> create emote set on channel A -> create X subflows for add emote to createdEmoteSet

var CopySet = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv copyset",
		Description: null.StringFrom("Copy set from other channel"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "7tv",
		Visible: true,
		IsReply: true,
		Aliases: []string{},
		Enabled: false,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: copySetChannelName,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.CopySetChannelName)
			},
		},
		command_arguments.VariadicString{
			Name: copySetNameOfSet,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.CopySetNameOfSet)
			},
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		client := seventvintegration.NewClient(parseCtx.Services.Config.SevenTvToken)

		if len(parseCtx.Mentions) != 1 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.ShouldMentionWithAt,
				),
			}
		}

		profile, err := client.GetProfileByTwitchId(ctx, parseCtx.Mentions[0].UserId)
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

		var setName string
		if profile.Users.UserByConnection.Style.ActiveEmoteSet != nil {
			setName = profile.Users.UserByConnection.Style.ActiveEmoteSet.Name
		}
		if name := parseCtx.ArgsParser.Get(copySetNameOfSet); name != nil {
			setName = name.String()
		}

		var targetSet *seventvintegrationapi.TwirSeventvUserEmoteSetsEmoteSet
		for _, set := range profile.Users.UserByConnection.EmoteSets {
			if set.Name == setName {
				targetSet = &set
				break
			}
		}

		if targetSet == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmoteNotFound.
						SetVars(locales.KeysCommandsSeventvErrorsEmoteNotFoundVars{EmoteName: setName}),
				),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				"",
			},
		}, nil
	},
}
