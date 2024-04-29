package shoutout

import (
	"context"
	"errors"
	"fmt"

	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"

	"github.com/samber/lo"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
)

const (
	userArgName = "@nickname"
)

var ShoutOut = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "so",
		Description: null.StringFrom("Shoutout some streamer"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: userArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		token, err := parseCtx.Services.GrpcClients.Tokens.RequestUserToken(
			ctx, &tokens.GetUserTokenRequest{
				UserId: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			result.Result = append(result.Result, "internal error")
			return result, nil
		}

		_, ok := lo.Find(
			token.Scopes, func(item string) bool {
				return item == "moderator:manage:shoutouts"
			},
		)
		if !ok {
			result.Result = append(
				result.Result,
				"we have no permissions for shoutout. Streamer must re-authorize to bot dashboard.",
			)
			return result, nil
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		if len(parseCtx.Mentions) == 0 {
			result.Result = []string{"user not found."}
			return result, nil
		}

		user := parseCtx.Mentions[0]

		go twitchClient.SendShoutout(
			&helix.SendShoutoutParams{
				FromBroadcasterID: parseCtx.Channel.ID,
				ToBroadcasterID:   user.UserId,
				ModeratorID:       parseCtx.Channel.ID,
			},
		)

		streamsReq, err := twitchClient.GetStreams(
			&helix.StreamsParams{
				UserIDs: []string{user.UserId},
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get stream",
				Err:     err,
			}
		}
		if streamsReq.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: "cannot get stream",
				Err:     errors.New(streamsReq.ErrorMessage),
			}
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
			return result, nil
		} else {
			channelReq, err := twitchClient.GetChannelInformation(
				&helix.GetChannelInformationParams{
					BroadcasterIDs: []string{user.UserId},
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot get channel",
					Err:     err,
				}
			}
			if channelReq.ErrorMessage != "" {
				return nil, &types.CommandHandlerError{
					Message: "cannot get channel",
					Err:     errors.New(channelReq.ErrorMessage),
				}
			}

			if len(channelReq.Data.Channels) == 0 {
				result.Result = append(result.Result, "cannot find user with this name.")
				return result, nil
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
			return result, nil
		}
	},
}
