package seventv

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
	"go.uber.org/zap"
)

const emoteForDeleteArgName = "name"

var EmoteDelete = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv delete",
		Description: null.StringFrom("Delete 7tv emote in current set"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeBroadcaster.String()},
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

		if len(profile.EmoteSet.Emotes) == 0 {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get 7tv emotes"),
				Err:     fmt.Errorf("emotes are not set"),
			}
		}

		arg := strings.ToLower(parseCtx.ArgsParser.Get(emoteForDeleteArgName).String())

		var foundEmote *seventvintegration.Emote
		for _, emote := range profile.EmoteSet.Emotes {
			loweredName := strings.ToLower(emote.Name)
			if loweredName == arg {
				foundEmote = &emote
				break
			}
		}

		if foundEmote == nil {
			return &types.CommandsHandlerResult{
				Result: []string{
					fmt.Sprintf(`Emote "%s" not found`, arg),
				},
			}, nil
		}

		err = seventvintegration.RemoveEmoteByID(
			ctx,
			parseCtx.Services.Config.SevenTvToken,
			profile.EmoteSet.Id,
			foundEmote.Id,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to rename 7tv emote"),
				Err:     err,
			}
		}

		if err := parseCtx.Services.SevenTvCache.Invalidate(ctx, parseCtx.Channel.ID); err != nil {
			parseCtx.Services.Logger.Error("cannot invalidate channel seventv cache", zap.Error(err))
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(`✅ Emote "%s" removed`, foundEmote.Name),
			},
		}, nil
	},
}
