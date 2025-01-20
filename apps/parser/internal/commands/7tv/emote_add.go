package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	"go.uber.org/zap"
)

const emoteForAddArgLink = "link"

var EmoteAdd = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv add",
		Description: null.StringFrom("Add 7tv emote to current set by link"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeBroadcaster.String()},
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
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		profile, err := parseCtx.Services.SevenTvCache.Get(ctx, parseCtx.Channel.ID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get 7tv profile"),
				Err:     err,
			}
		}

		if profile.EmoteSet == nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get 7tv emote set"),
				Err:     fmt.Errorf("emote set is not set"),
			}
		}

		arg := parseCtx.ArgsParser.Get(emoteForAddArgLink).String()

		link := seventvintegration.FindEmoteIdInInput(arg)
		if link == "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Invalid arg, link should be used"),
			}
		}

		err = seventvintegration.AddEmote(
			ctx,
			parseCtx.Services.Config.SevenTvToken,
			arg,
			profile.EmoteSet.Id,
		)

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to add 7tv emote: %s", err.Error()),
				Err:     err,
			}
		}

		if err := parseCtx.Services.SevenTvCache.Invalidate(ctx, parseCtx.Channel.ID); err != nil {
			parseCtx.Services.Logger.Error("cannot invalidate channel seventv cache", zap.Error(err))
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(`✅ Emote added`),
			},
		}, nil
	},
}
