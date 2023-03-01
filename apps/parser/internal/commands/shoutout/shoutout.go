package shoutout

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
)

var ShoutOut = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "so",
		Description: null.StringFrom("Shoutout some streamer"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
		result := &types.CommandsHandlerResult{}

		if ctx.Text == nil || *ctx.Text == "" {
			result.Result = append(result.Result, "you have to type streamer for shoutout.")
			return result
		}

		token, err := tokensGrpc.RequestUserToken(context.Background(), &tokens.GetUserTokenRequest{
			UserId: ctx.ChannelId,
		})

		if err != nil {
			result.Result = append(result.Result, "internal error")
			return result
		}

		_, ok := lo.Find(token.Scopes, func(item string) bool {
			return item == "moderator:manage:shoutouts"
		})
		if !ok {
			result.Result = append(result.Result, "we have no permissions for shoutout. Streamer must re-authorize to bot dashboard.")
			return result
		}

		twitchClient, err := twitch.NewUserClient(ctx.ChannelId, cfg, tokensGrpc)
		if err != nil {
			return nil
		}

		usersReq, err := twitchClient.GetUsers(&helix.UsersParams{
			Logins: []string{*ctx.Text},
		})
		if err != nil || len(usersReq.Data.Users) == 0 {
			return nil
		}

		user := usersReq.Data.Users[0]

		go twitchClient.SendShoutout(&helix.SendShoutoutParams{
			FromBroadcasterID: ctx.ChannelId,
			ToBroadcasterID:   user.ID,
			ModeratorID:       ctx.ChannelId,
		})

		streamsReq, err := twitchClient.GetStreams(&helix.StreamsParams{
			UserIDs: []string{user.ID},
		})
		if err != nil {
			return nil
		}

		if len(streamsReq.Data.Streams) != 0 {
			stream := streamsReq.Data.Streams[0]

			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Check out amazing %s, streaming %s - %s for %v viewers",
					stream.UserName,
					stream.GameName,
					stream.Title,
					stream.ViewerCount,
				),
			)
			return result
		} else {
			channelReq, err := twitchClient.GetChannelInformation(&helix.GetChannelInformationParams{
				BroadcasterIDs: []string{user.ID},
			})
			if err != nil || len(channelReq.Data.Channels) == 0 {
				return nil
			}
			channel := channelReq.Data.Channels[0]
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Check out amazing %s, was streaming %s - %s",
					channel.BroadcasterName,
					channel.GameName,
					channel.Title,
				),
			)
			return result
		}
	},
}
