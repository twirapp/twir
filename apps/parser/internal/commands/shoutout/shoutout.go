package shoutout

import (
	"context"

	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/guregu/null"
	"github.com/kvizyx/twitchy/helix"
	"github.com/lib/pq"

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

		var token model.Tokens
		err := parseCtx.Services.Gorm.
			Table("users").
			Select("tokens.scopes").
			Joins(`JOIN tokens ON users."tokenId" = tokens.id`).
			Where("users.id = ?", parseCtx.Channel.TwitchUserID).
			First(&token).Error
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

		if !hasShoutoutScope(token.Scopes) {
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
			parseCtx.Channel.TwitchUserID,
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

		go func() {
			_, _ = twitchClient.Chat.SendShoutout(
				ctx,
				helix.SendShoutoutRequest{
					FromBroadcasterID: parseCtx.Channel.ID,
					ToBroadcasterID:   user.UserID,
					ModeratorID:       parseCtx.Channel.ID,
				},
			)
		}()

		streamsReq, err := twitchClient.Streams.GetStreams(
			ctx,
			helix.GetStreamsRequest{
				UserID: []string{user.UserID},
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
		if len(streamsReq.Data) != 0 {
			stream := streamsReq.Data[0]

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
			channelReq, err := twitchClient.Channels.GetChannelInformation(
				ctx,
				helix.GetChannelInformationRequest{
					BroadcasterIDs: []string{user.UserID},
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
			if len(channelReq.Data) == 0 {
				result.Result = append(
					result.Result,
					i18n.GetCtx(
						ctx,
						locales.Translations.Errors.Generic.CannotFindChannelTwitch.SetVars(locales.KeysErrorsGenericCannotFindChannelTwitchVars{Reason: ""}),
					),
				)
				return result, nil
			}
			channel := channelReq.Data[0]
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

func hasShoutoutScope(scopes []string) bool {
	return lo.Contains(scopes, "moderator:manage:shoutouts")
}
