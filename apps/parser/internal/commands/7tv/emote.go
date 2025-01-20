package seventv

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
	model "github.com/satont/twir/libs/gomodels"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
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

		arg := strings.ToLower(parseCtx.ArgsParser.Get(emoteFindArgName).String())

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

		adderProfile, err := parseCtx.Services.SevenTvCacheBySevenTvID.Get(
			ctx,
			foundEmote.ActorId,
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get adder 7tv profile"),
			}
		}

		if adderProfile == nil {
			parseCtx.Services.Logger.Sugar().Error("Failed to get adder 7tv profile")
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get adder 7tv profile"),
				Err:     err,
			}
		}

		addedBy := fmt.Sprintf("%s", adderProfile.Username)
		emoteLink := fmt.Sprintf("https://7tv.app/emotes/%s", foundEmote.Id)

		addedAgo := helpers.Duration(
			time.UnixMilli(foundEmote.Timestamp),
			&helpers.DurationOpts{
				UseUtc: true,
				Hide: helpers.DurationOptsHide{
					Seconds: true,
				},
			},
		)
		author := fmt.Sprintf(
			"%s: https://7tv.app/users/%s",
			foundEmote.Data.Owner.Username,
			foundEmote.Data.Owner.Id,
		)

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"%s: %s · Added by @%s %v ago · Author %s",
					foundEmote.Name,
					emoteLink,
					addedBy,
					addedAgo,
					author,
				),
			},
		}, nil
	},
}
