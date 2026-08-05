package sr_youtube

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"

	"github.com/twirapp/twir/libs/bus-core/api"
	song_request_mode "github.com/twirapp/twir/libs/entities/song_request_mode"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
)

var SkipCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "voteskip",
		Description: null.StringFrom("Vote for skip command"),
		Module:      "SONGS",
		IsReply:     true,
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
			result.Result = append(result.Result, "Voteskip not supported in Spotify mode")
			return result, nil
		}

		currentSong := &model.RequestedSong{}
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ?::uuid AND "deletedAt" IS NULL`, parseCtx.Channel.DBChannelID).
			Order(`"createdAt" asc`).
			Limit(1).
			Find(&currentSong).
			Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetCurrentSong),
				Err:     err,
			}
		}

		if currentSong.ID == "" {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.NotFound))
			return result, nil
		}

		var onlineUsersCount int64
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ?::uuid`, parseCtx.Channel.DBChannelID).
			Model(&model.UsersOnline{}).
			Count(&onlineUsersCount).
			Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetUsersCount),
				Err:     err,
			}
		}

		redisKey := fmt.Sprintf("songrequests-voteskip-%s", currentSong.ID)
		votesCount, err := parseCtx.Services.Redis.SCard(ctx, redisKey).Result()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetVotesCount),
				Err:     err,
			}
		}

		currentVote, err := parseCtx.Services.Redis.SIsMember(
			ctx,
			redisKey,
			parseCtx.Sender.ID,
		).Result()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetCurrentVote),
				Err:     err,
			}
		}

		neededVotes := int64(math.Round(moduleSettings.NeededVotesForSkip * float64(onlineUsersCount) / 100))

		if currentVote {
			result.Result = append(result.Result, fmt.Sprintf("%v/%v", votesCount, neededVotes))
			return result, nil
		}

		parseCtx.Services.Redis.SAdd(ctx, redisKey, parseCtx.Sender.ID)
		parseCtx.Services.Redis.Expire(ctx, redisKey, 1*time.Hour)

		if votesCount+1 >= neededVotes {
			parseCtx.Services.Bus.Api.SongRequestRemoveFromQueue.Publish(
				ctx,
				api.SongRequestRemoveFromQueue{
					ChannelID: parseCtx.Channel.DBChannelID,
					VideoID:   currentSong.VideoID,
				},
			)

			currentSong.DeletedAt = lo.ToPtr(time.Now().UTC())
			parseCtx.Services.Gorm.WithContext(ctx).Updates(currentSong)
			parseCtx.Services.Redis.Del(ctx, redisKey)

			result.Result = append(result.Result, i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Songrequest.Info.SongSkipped.
					SetVars(locales.KeysCommandsSongrequestInfoSongSkippedVars{SongTitle: currentSong.Title})),
			)
			return result, nil
		}

		result.Result = append(result.Result, fmt.Sprintf("%v/%v", votesCount+1, neededVotes))
		return result, nil
	},
}
