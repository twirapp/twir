package shoutout

import (
	"context"
	"errors"

	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/tokens"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"

	"github.com/samber/lo"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
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
		IsReply:     true,
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

		token, err := parseCtx.Services.Bus.Tokens.RequestUserToken.Request(
			ctx,
			tokens.GetUserTokenRequest{
				UserId: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotFindChannelDb,
				),
			)
			return result, nil
		}

		_, ok := lo.Find(
			token.Data.Scopes, func(item string) bool {
				return item == "moderator:manage:shoutouts"
			},
		)
		if !ok {
			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Shoutout.Errors.BotHaveNoPermissions,
				),
			)
			return result, nil
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.BroadcasterClient,
				),
				Err: err,
			}
		}

		if len(parseCtx.Mentions) == 0 {
			result.Result = []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.ShouldMentionWithAt,
				),
			}
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
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotGetStream.SetVars(locales.KeysErrorsGenericCannotGetStreamVars{Reason: err.Error()}),
				),
				Err: err,
			}
		}
		if streamsReq.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotGetStream.SetVars(locales.KeysErrorsGenericCannotGetStreamVars{Reason: streamsReq.ErrorMessage}),
				),
				Err: errors.New(streamsReq.ErrorMessage),
			}
		}

		if len(streamsReq.Data.Streams) != 0 {
			stream := streamsReq.Data.Streams[0]

			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Shoutout.ResponseOnline.SetVars(
						locales.KeysCommandsShoutoutResponseOnlineVars{
							UserName:     stream.UserName,
							CategoryName: stream.GameName,
							Title:        stream.Title,
							Viewers:      stream.ViewerCount,
						},
					),
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
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Errors.Generic.CannotFindChannelTwitch.SetVars(locales.KeysErrorsGenericCannotFindChannelTwitchVars{Reason: err.Error()}),
					),
					Err: err,
				}
			}
			if channelReq.ErrorMessage != "" {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Errors.Generic.CannotFindChannelTwitch.SetVars(locales.KeysErrorsGenericCannotFindChannelTwitchVars{Reason: channelReq.ErrorMessage}),
					),
					Err: errors.New(channelReq.ErrorMessage),
				}
			}

			if len(channelReq.Data.Channels) == 0 {
				result.Result = append(
					result.Result,
					i18n.GetCtx(
						ctx,
						locales.Translations.Errors.Generic.CannotFindChannelTwitch.SetVars(locales.KeysErrorsGenericCannotFindChannelTwitchVars{Reason: ""}),
					),
				)
				return result, nil
			}
			channel := channelReq.Data.Channels[0]
			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Shoutout.ResponseOffline.SetVars(
						locales.KeysCommandsShoutoutResponseOfflineVars{
							UserName:     channel.BroadcasterName,
							CategoryName: channel.GameName,
							Title:        channel.Title,
						},
					),
				),
			)
			return result, nil
		}
	},
}
