package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	seventvvariables "github.com/twirapp/twir/apps/parser/internal/variables/7tv"
	model "github.com/twirapp/twir/libs/gomodels"
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
			Hint:     "@channelName",
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
					"7tv profile not found",
				},
			}, nil
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s) · Paint: $(%s) ($(%s) unlocked) · Roles: $(%s) · Editor for $(%s) · Set: $(%s) ($(%s)/$(%s)) · Created: $(%s)",
					seventvvariables.ProfileLink.Name,
					seventvvariables.Paint.Name,
					seventvvariables.UnlockedPaints.Name,
					seventvvariables.Roles.Name,
					seventvvariables.EditorForCount.Name,
					seventvvariables.EmoteSetName.Name,
					seventvvariables.EmoteSetCount.Name,
					seventvvariables.EmoteSetCapacity.Name,
					seventvvariables.ProfileCreatedAt.Name,
				),
			},
		}

		return result, nil
	},
}
