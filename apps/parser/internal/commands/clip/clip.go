package clip

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kvizyx/twitchy/helix"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
		channelID, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotFindChannelDb,
				),
				Err: err,
			}
		}
		channel, err := parseCtx.Services.ChannelsRepo.GetByID(ctx, channelID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotFindChannelDb,
				),
				Err: err,
			}
		}
		twitchBinding, twitchBotConfig, ok, err := channel.TwitchBinding()
		if err != nil || !ok {
			if err == nil {
				err = errors.New("channel has no Twitch binding")
			}
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotFindChannelDb,
				),
				Err: err,
			}
		}

		twitchClient, err := twitch.NewChannelBotClientWithContext(
			ctx,
			twitchBotConfig.BotID,
			twitchBinding.PlatformChannelID,
			*parseCtx.Services.Config,
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

		resp, err := twitchClient.Clips.CreateClip(
			ctx,
			helix.CreateClipRequest{
				BroadcasterID: twitchBinding.PlatformChannelID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Clip.CannotCreateClip,
				),
				Err: err,
			}
		}
		if len(resp.Data) == 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Clip.EmptyClipUrl,
				),
				Err: errors.New("empty clip edit url"),
			}
		}

		clipId := resp.Data[0].ID

		var url string

		for i := 0; i < 20; i++ {
			clip, err := twitchClient.Clips.GetClips(
				ctx,
				helix.GetClipsRequest{
					IDs: []string{clipId},
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Clip.CannotGetClip,
					),
					Err: err,
				}
			}

			if len(clip.Data) > 0 {
				url = clip.Data[0].URL
				break
			}

			time.Sleep(1 * time.Second)
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Clip.ClipCreated.SetVars(
						locales.KeysCommandsClipClipCreatedVars{
							Url: url,
						},
					),
				),
			},
		}, nil
	},
}
