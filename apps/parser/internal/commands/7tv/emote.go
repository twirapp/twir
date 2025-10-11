package seventv

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
)

const emoteFindArgName = "emoteName"

var EmoteFind = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv emote",
		Description: null.StringFrom("Search emote by name in current set"),
		RolesIDS:    pq.StringArray{},
		Module:      "7tv",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     false,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: emoteFindArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		client := seventvintegration.NewClient(parseCtx.Services.Config.SevenTvToken)

		profile, err := client.GetProfileByTwitchId(ctx, parseCtx.Channel.ID)
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

		if profile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmotesetNotActive,
				),
			}
		}

		arg := strings.ToLower(parseCtx.ArgsParser.Get(emoteFindArgName).String())

		var foundEmote *seventvintegrationapi.TwirSeventvUserStyleActiveEmoteSetEmotesEmoteSetEmoteSearchResultItemsEmoteSetEmote
		for _, emote := range profile.Users.UserByConnection.Style.ActiveEmoteSet.Emotes.Items {
			if strings.ToLower(emote.Alias) == arg {
				foundEmote = &emote
				break
			}

			if strings.ToLower(emote.Emote.DefaultName) == arg {
				foundEmote = &emote
				break
			}
		}

		if foundEmote == nil {
			return &types.CommandsHandlerResult{
				Result: []string{
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Seventv.Errors.EmoteNotFound.
							SetVars(locales.KeysCommandsSeventvErrorsEmoteNotFoundVars{EmoteName: arg}),
					),
				},
			}, nil
		}

		adderProfile, err := client.GetProfileById(ctx, *foundEmote.AddedById)
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

		if adderProfile == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.ProfileFailedToGet.
						SetVars(locales.KeysCommandsSeventvErrorsProfileFailedToGetVars{Reason: err.Error()}),
				),
				Err: err,
			}
		}

		emoteLink := "https://7tv.app/emotes/" + foundEmote.Id

		addedAgo := helpers.Duration(
			foundEmote.AddedAt,
			&helpers.DurationOpts{
				UseUtc: true,
				Hide: helpers.DurationOptsHide{
					Seconds: true,
				},
			},
		)
		author := fmt.Sprintf(
			"%s: https://7tv.app/users/%s",
			foundEmote.Emote.Owner.MainConnection.PlatformDisplayName,
			foundEmote.Emote.Owner.Id,
		)

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx, locales.Translations.Commands.Seventv.EmoteInfo.Response.
						SetVars(locales.KeysCommandsSeventvEmoteInfoResponseVars{
							Name:            foundEmote.Emote.DefaultName,
							Link:            emoteLink,
							AddedByUserName: adderProfile.Users.User.MainConnection.PlatformDisplayName,
							AddedByTime:     addedAgo,
							EmoteAuthor:     author,
						}),
				),
			},
		}, nil
	},
}
