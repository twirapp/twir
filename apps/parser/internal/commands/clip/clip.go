package clip

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

var MakeClip = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "clip",
		Description: null.StringFrom("Create clip"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "CLIPS",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		channel := &model.Channels{}
		if err := parseCtx.Services.Gorm.Where(
			"id = ?",
			parseCtx.Channel.ID,
		).First(channel).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot find channel in db, please contact support",
				Err:     err,
			}
		}

		twitchClient, err := twitch.NewBotClientWithContext(
			ctx,
			channel.BotID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		resp, err := twitchClient.CreateClip(
			&helix.CreateClipParams{
				BroadcasterID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create clip",
				Err:     err,
			}
		}
		if resp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("cannot create clip: %s", resp.ErrorMessage),
				Err:     errors.New(resp.ErrorMessage),
			}
		}

		if len(resp.Data.ClipEditURLs) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "cannot create clip",
				Err:     errors.New("empty clip edit url"),
			}
		}

		clipId := resp.Data.ClipEditURLs[0].ID

		var url string

		for i := 0; i < 20; i++ {
			clip, err := twitchClient.GetClips(
				&helix.ClipsParams{
					IDs: []string{clipId},
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot get clip",
					Err:     err,
				}
			}

			if len(clip.Data.Clips) > 0 {
				url = clip.Data.Clips[0].URL
				break
			}

			time.Sleep(1 * time.Second)
		}

		return &types.CommandsHandlerResult{Result: []string{url}}, nil
	},
}
