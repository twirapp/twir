package shoutout

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/apps/parser/internal/types"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"

	"github.com/samber/lo"

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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			result.Result = append(result.Result, "you have to type streamer for shoutout.")
			return result
		}

		token, err := parseCtx.Services.GrpcClients.Tokens.RequestUserToken(ctx, &tokens.GetUserTokenRequest{
			UserId: parseCtx.Channel.ID,
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

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil
		}

		usersReq, err := twitchClient.GetUsers(&helix.UsersParams{
			Logins: []string{*parseCtx.Text},
		})
		if err != nil || len(usersReq.Data.Users) == 0 {
			return nil
		}

		user := usersReq.Data.Users[0]

		go twitchClient.SendShoutout(&helix.SendShoutoutParams{
			FromBroadcasterID: parseCtx.Channel.ID,
			ToBroadcasterID:   user.ID,
			ModeratorID:       parseCtx.Channel.ID,
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
