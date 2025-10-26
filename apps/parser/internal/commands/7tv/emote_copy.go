package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	seventvapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	"golang.org/x/sync/errgroup"
)

const (
	emoteForCopyArgName = "emoteName"
	emoteForCopyChannel = "channel"
	emoteForCopyAlias   = "alias"
)

var EmoteCopy = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv copy",
		Description: null.StringFrom("Copy 7tv emote from other channel"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "7tv",
		Visible:     true,
		IsReply:     true,
		Aliases: []string{
			"7tv yoink",
		},
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: emoteForCopyArgName,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.EmoteForCopyArgName)
			},
		},
		command_arguments.String{
			Name:     emoteForCopyChannel,
			Optional: false,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Errors.Generic.ShouldMentionWithAt)
			},
		},
		command_arguments.String{
			Name:     emoteForCopyAlias,
			Optional: true,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.EmoteForCopyAlias)
			},
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if len(parseCtx.Mentions) == 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.ShouldMentionWithAt,
				),
			}
		}

		client := seventvintegration.NewClient(parseCtx.Services.Config.SevenTvToken)

		var (
			wg                 errgroup.Group
			broadcasterProfile *seventvapi.GetProfileByTwitchIdResponse
			targetProfile      *seventvapi.GetProfileByTwitchIdResponse
		)

		wg.Go(
			func() error {
				broadcasterSeventvProfile, err := client.GetProfileByTwitchId(ctx, parseCtx.Channel.ID)
				if err != nil {
					return fmt.Errorf(
						i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Seventv.Errors.ProfileFailedToGet.
								SetVars(locales.KeysCommandsSeventvErrorsProfileFailedToGetVars{Reason: err.Error()}),
						),
					)
				}
				if broadcasterSeventvProfile.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
					return fmt.Errorf(
						i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Seventv.Errors.EmotesetNotActive,
						),
					)
				}

				broadcasterProfile = broadcasterSeventvProfile
				return nil
			},
		)

		wg.Go(
			func() error {
				targetSeventvProfile, err := client.GetProfileByTwitchId(ctx, parseCtx.Mentions[0].UserId)
				if err != nil {
					return fmt.Errorf(
						i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Seventv.Errors.ProfileFailedToGet.
								SetVars(locales.KeysCommandsSeventvErrorsProfileFailedToGetVars{Reason: err.Error()}),
						),
					)
				}
				if targetSeventvProfile.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
					return fmt.Errorf(
						i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Seventv.Errors.ProfileFailedToGet.
								SetVars(locales.KeysCommandsSeventvErrorsProfileFailedToGetVars{Reason: err.Error()}),
						),
					)
				}

				targetProfile = targetSeventvProfile
				return nil
			},
		)

		if err := wg.Wait(); err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		if broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmotesetBroadcasterNotActive,
				),
			}
		}

		emoteForSearch := parseCtx.ArgsParser.Get(emoteForCopyArgName).String()
		var targetEmote *seventvapi.TwirSeventvUserStyleActiveEmoteSetEmotesEmoteSetEmoteSearchResultItemsEmoteSetEmote
		for _, e := range targetProfile.Users.UserByConnection.Style.ActiveEmoteSet.Emotes.Items {
			if e.Alias == emoteForSearch || e.Emote.DefaultName == emoteForSearch {
				targetEmote = &e
				break
			}
		}

		if targetEmote == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmoteNotFoundInChannel.
						SetVars(locales.KeysCommandsSeventvErrorsEmoteNotFoundInChannelVars{EmoteSearch: emoteForSearch}),
				),
			}
		}

		emoteName := targetEmote.Alias
		if alias := parseCtx.ArgsParser.Get(emoteForCopyAlias); alias != nil {
			emoteName = alias.String()
		}

		for _, e := range broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet.Emotes.Items {
			if e.Emote.Id == targetEmote.Emote.Id ||
				e.Alias == emoteForSearch ||
				e.Emote.DefaultName == emoteForSearch ||
				e.Alias == emoteName {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Seventv.Errors.EmoteAlreadyExistInChannel.
							SetVars(locales.KeysCommandsSeventvErrorsEmoteAlreadyExistInChannelVars{EmoteName: emoteName}),
					),
				}
			}
		}

		err := client.AddEmote(
			ctx,
			*broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSetId,
			targetEmote.Emote.Id,
			emoteName,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmoteFailedToAdd.
						SetVars(locales.KeysCommandsSeventvErrorsEmoteFailedToAddVars{Reason: err.Error()}),
				),
				Err: err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Add.EmoteAdd,
				),
			},
		}, nil
	},
}
