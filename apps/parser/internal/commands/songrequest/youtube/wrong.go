package sr_youtube

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/guregu/null"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"

	"github.com/twirapp/twir/libs/bus-core/api"
	buscorespotify "github.com/twirapp/twir/libs/bus-core/spotify"
	song_request_mode "github.com/twirapp/twir/libs/entities/song_request_mode"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"

	"github.com/samber/lo"
)

const (
	songSkipArgName = "number"
)

var WrongCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr wrong",
		Description: null.StringFrom("Delete wrong song from queue"),
		Module:      "SONGS",
		IsReply:     true,
		Visible:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     songSkipArgName,
			Optional: true,
			Min:      lo.ToPtr(1),
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		moduleSettings, err := parseCtx.Services.SongRequestsSettingsRepo.GetByChannelID(ctx, parseCtx.Channel.DBChannelID)
		if err != nil {
			if errors.Is(err, songrequestssettingsrepository.ErrNotFound) {
				return result, nil
			}

			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetSettings),
				Err:     err,
			}
		}

		if !moduleSettings.Enabled {
			return result, nil
		}

		if moduleSettings.Mode == song_request_mode.ModeSpotify {
			request, err := parseCtx.Services.Bus.Spotify.CancelSongRequest.Request(
				ctx,
				buscorespotify.CancelSongRequestRequest{
					ChannelID:     parseCtx.Channel.DBChannelID,
					RequesterName: parseCtx.Sender.Name,
				},
			)
			if err != nil {
				message := spotifySongRequestErrorMessage(err)
				if errors.Is(normalizeSpotifySongRequestError(err), errSpotifyTrackNotFound) {
					message = "No active Spotify request to cancel"
				}
				if message == "" {
					parseCtx.Services.Logger.Sugar().Errorw(
						"spotify song request cancel failed",
						"error", err,
						"channel_id", parseCtx.Channel.DBChannelID,
						"user_id", parseCtx.Sender.DbUser.ID,
					)
					message = "Spotify request cancel failed"
				}

				result.Result = append(result.Result, message)
				return result, nil
			}

			result.Result = append(
				result.Result,
				fmt.Sprintf("Removed Spotify request: %s", request.Data.Request.Title),
			)
			return result, nil
		}

		var songs []*model.RequestedSong
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(
				`"channelId" = ?::uuid AND "orderedById" = ? AND "deletedAt" IS NULL`,
				parseCtx.Channel.DBChannelID,
				parseCtx.Sender.DbUser.ID,
			).
			Limit(5).
			Order(`"createdAt" desc`).
			Find(&songs).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetSongsFromQueue),
				Err:     err,
			}
		}

		if len(songs) == 0 {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Info.NoRequestedSongs))
			return result, nil
		}

		number := 1
		songSkipArg := parseCtx.ArgsParser.Get(songSkipArgName)
		if songSkipArg != nil {
			number = songSkipArg.Int()
		}

		if number < 1 || number > len(songs) {
			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Songrequest.Info.OnlyCountSongs.
						SetVars(locales.KeysCommandsSongrequestInfoOnlyCountSongsVars{SongsCount: len(songs)}),
				),
			)
			return result, nil
		}

		choosedSong := songs[number-1]
		choosedSong.DeletedAt = lo.ToPtr(time.Now().UTC())
		err = parseCtx.Services.Gorm.WithContext(ctx).Updates(&choosedSong).Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.UpdateSong),
				Err:     err,
			}
		}

		parseCtx.Services.Bus.Api.SongRequestRemoveFromQueue.Publish(
			ctx,
			api.SongRequestRemoveFromQueue{
				ChannelID: parseCtx.Channel.DBChannelID,
				VideoID:   choosedSong.VideoID,
			},
		)

		result.Result = append(
			result.Result,
			i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Songrequest.Info.Delete.
					SetVars(locales.KeysCommandsSongrequestInfoDeleteVars{SongTitle: choosedSong.Title}),
			),
		)

		return result, nil
	},
}
