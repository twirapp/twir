package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
)

const (
	emoteForAddArgLink   = "link"
	emoteForAddArgAliase = "aliase"
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
			Hint: "link or name",
		},
		command_arguments.String{
			Name:     emoteForAddArgAliase,
			Optional: true,
			Hint:     "optional name",
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
				Message: fmt.Sprintf("Failed to get 7tv profile: %v", err),
				Err:     err,
			}
		}
		if sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
			return &types.CommandsHandlerResult{
				Result: []string{
					`❌ No active emote set`,
				},
			}, nil
		}

		nameOrLinkArgument := parseCtx.ArgsParser.Get(emoteForAddArgLink).String()

		emote, err := client.GetOneEmoteByNameOrLink(ctx, nameOrLinkArgument)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get 7tv emote: %v", err),
				Err:     err,
			}
		}

		var emoteName string
		aliaseArgument := parseCtx.ArgsParser.Get(emoteForAddArgAliase)
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
				Message: fmt.Sprintf("Failed to add 7tv emote: %v", err),
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				`✅ Emote added`,
			},
		}, nil
	},
}
