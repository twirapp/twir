package seventv

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
	model "github.com/satont/twir/libs/gomodels"
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
				Message: fmt.Sprintf("Failed to get 7tv profile: %v", err),
				Err:     err,
			}
		}

		if profile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Emote set is not set"),
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
					fmt.Sprintf(`Emote "%s" not found`, arg),
				},
			}, nil
		}

		adderProfile, err := client.GetProfileById(ctx, *foundEmote.AddedById)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get 7tv profile of adder: %v", err),
				Err:     err,
			}
		}

		if adderProfile == nil {
			parseCtx.Services.Logger.Sugar().Error("Failed to get adder 7tv profile")
			return nil, &types.CommandHandlerError{
				Message: "Failed to get adder 7tv profile",
				Err:     err,
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
				fmt.Sprintf(
					"%s: %s · Added by @%s %v ago · Author %s",
					foundEmote.Emote.DefaultName,
					emoteLink,
					adderProfile.Users.User.MainConnection.PlatformDisplayName,
					addedAgo,
					author,
				),
			},
		}, nil
	},
}
