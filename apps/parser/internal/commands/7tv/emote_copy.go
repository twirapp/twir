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
	seventvapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	"golang.org/x/sync/errgroup"
)

const (
	emoteForCopyArgName = "emoteName"
	emoteForCopyChannel = "channel"
	emoteForCopyAliase  = "aliase"
)

var EmoteCopy = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv copy",
		Description: null.StringFrom("Copy 7tv emote from other channel"),
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
			Name: emoteForCopyArgName,
			Hint: "name of emote to copy",
		},
		command_arguments.String{
			Name:     emoteForCopyChannel,
			Optional: false,
			Hint:     "@channel",
		},
		command_arguments.String{
			Name:     emoteForCopyAliase,
			Optional: true,
			Hint:     "aliase for emote",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if len(parseCtx.Mentions) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "you should tag user with @",
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
					return fmt.Errorf("failed to get 7tv profile: %w", err)
				}
				if broadcasterSeventvProfile.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
					return fmt.Errorf(`❌ No active emote set`)
				}

				broadcasterProfile = broadcasterSeventvProfile
				return nil
			},
		)

		wg.Go(
			func() error {
				targetSeventvProfile, err := client.GetProfileByTwitchId(ctx, parseCtx.Mentions[0].UserId)
				if err != nil {
					return fmt.Errorf("failed to get 7tv profile: %w", err)
				}
				if targetSeventvProfile.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
					return fmt.Errorf(`❌ No active emote set`)
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
				Message: `❌ No active emote set for broadcaster`,
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
				Message: fmt.Sprintf(`❌ Emote "%s" not found in target channel`, emoteForSearch),
			}
		}

		emoteName := targetEmote.Alias
		if alias := parseCtx.ArgsParser.Get(emoteForCopyAliase); alias != nil {
			emoteName = alias.String()
		}

		for _, e := range broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet.Emotes.Items {
			if e.Emote.Id == targetEmote.Emote.Id ||
				e.Alias == emoteForSearch ||
				e.Emote.DefaultName == emoteForSearch ||
				e.Alias == emoteName {
				return nil, &types.CommandHandlerError{
					Message: fmt.Sprintf(`❌ Emote "%s" already exists in this channel`, emoteName),
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
