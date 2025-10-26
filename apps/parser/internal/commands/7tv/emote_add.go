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

const (
	emoteForAddArgLink  = "link"
	emoteForAddArgAlias = "alias"
)

var EmoteAdd = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv add",
		Description: null.StringFrom("Add 7tv emote to current set by link"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "7tv",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     false,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: emoteForAddArgLink,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.EmoteForAddArgLink)
			},
		},
		command_arguments.String{
			Name:     emoteForAddArgAlias,
			Optional: true,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.EmoteForAddArgAlias)
			},
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

		nameOrLinkArgument := parseCtx.ArgsParser.Get(emoteForAddArgLink).String()

		emote, err := client.GetOneEmoteByNameOrLink(ctx, nameOrLinkArgument)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Seventv.Errors.EmotesetNotActive,
				),
				Err: err,
			}
		}

		var emoteName string
		aliaseArgument := parseCtx.ArgsParser.Get(emoteForAddArgAlias)
		if aliaseArgument != nil {
			emoteName = aliaseArgument.String()
		} else {
			emoteName = emote.DefaultName
		}

		err = client.AddEmote(
			ctx,
			*sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSetId,
			emote.Id,
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
				i18n.Get(locales.Translations.Commands.Seventv.Add.EmoteAdd),
			},
		}, nil
	},
}
