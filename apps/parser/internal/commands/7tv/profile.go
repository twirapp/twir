package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	seventvvariables "github.com/twirapp/twir/apps/parser/internal/variables/7tv"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

const profileArg = "channelName"

var Profile = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv profile",
		Description: null.StringFrom("Information about 7tv profile"),
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
			Name:     profileArg,
			Optional: true,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Seventv.Hints.CopySetChannelName)
			},
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		userIdForCheck := parseCtx.Channel.ID
		if parseCtx.ArgsParser.Get(profileArg) != nil && len(parseCtx.Mentions) >= 1 {
			userIdForCheck = parseCtx.Mentions[0].UserId
		}

		if _, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, userIdForCheck); err != nil {
			return &types.CommandsHandlerResult{
				Result: []string{
					i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Seventv.Errors.ProfileNotFound,
					),
				},
			}, nil
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx, locales.Translations.Commands.Seventv.ProfileInfo.Response.
						SetVars(
							locales.KeysCommandsSeventvProfileInfoResponseVars{
								ProfileLink:      fmt.Sprintf("$(%s)", seventvvariables.ProfileLink.Name),
								PaintName:        fmt.Sprintf("$(%s)", seventvvariables.Paint.Name),
								UnlockedPaints:   fmt.Sprintf("$(%s)", seventvvariables.UnlockedPaints.Name),
								Roles:            fmt.Sprintf("$(%s)", seventvvariables.Roles.Name),
								EditorCount:      fmt.Sprintf("$(%s)", seventvvariables.EditorForCount.Name),
								EmoteSetName:     fmt.Sprintf("$(%s)", seventvvariables.EmoteSetName.Name),
								EmoteSetCount:    fmt.Sprintf("$(%s)", seventvvariables.EmoteSetCount.Name),
								EmoteSetCapacity: fmt.Sprintf("$(%s)", seventvvariables.EmoteSetCapacity.Name),
								ProfileCreatedAt: fmt.Sprintf("$(%s)", seventvvariables.ProfileCreatedAt.Name),
							},
						),
				),
			},
		}

		return result, nil
	},
}
