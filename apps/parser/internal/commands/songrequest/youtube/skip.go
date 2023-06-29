package sr_youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	youtube "github.com/satont/twir/libs/types/types/api/modules"
	"gorm.io/gorm"
	"math"
	"time"
)

var SkipCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "voteskip",
		Description: null.StringFrom("Vote for skip command"),
		Module:      "SONGS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		moduleSettings := &model.ChannelModulesSettings{}
		parsedSettings := &youtube.YouTubeSettings{}
		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "type" = ?`, parseCtx.Channel.ID, "youtube_song_requests").
			First(moduleSettings).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		if moduleSettings.ID == "" {
			result.Result = append(result.Result, "Song requests not enabled")
			return result
		}

		err = json.Unmarshal(moduleSettings.Settings, parsedSettings)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		if !*parsedSettings.Enabled {
			result.Result = append(result.Result, parsedSettings.Translations.NotEnabled)
			return result
		}

		currentSong := &model.RequestedSong{}
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, parseCtx.Channel.ID).
			Order(`"createdAt" asc`).
			Limit(1).
			Find(&currentSong).
			Error

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		if currentSong.ID == "" {
			result.Result = append(result.Result, "Current song not found")
			return result
		}

		var onlineUsersCount int64
		err = parseCtx.Services.Gorm.WithContext(ctx).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			Model(&model.UsersOnline{}).
			Count(&onlineUsersCount).
			Error

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		redisKey := fmt.Sprintf("songrequests-voteskip-%s", currentSong.ID)
		votesCount, err := parseCtx.Services.Redis.SCard(ctx, redisKey).Result()
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		currentVote, err := parseCtx.Services.Redis.SIsMember(ctx, redisKey, parseCtx.Sender.ID).Result()
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		neededVotes := int64(math.Round(parsedSettings.NeededVotesVorSkip * float64(onlineUsersCount) / 100))

		if currentVote {
			result.Result = append(result.Result, fmt.Sprintf("%v/%v", votesCount, neededVotes))
			return result
		}

		parseCtx.Services.Redis.SAdd(ctx, redisKey, parseCtx.Sender.ID)
		parseCtx.Services.Redis.Expire(ctx, redisKey, 1*time.Hour)

		if votesCount+1 >= neededVotes {
			_, err = parseCtx.Services.GrpcClients.WebSockets.YoutubeRemoveSongToQueue(
				ctx,
				&websockets.YoutubeRemoveSongFromQueueRequest{
					ChannelId: parseCtx.Channel.ID,
					EntityId:  currentSong.ID,
				},
			)

			if err != nil {
				parseCtx.Services.Logger.Sugar().Error(err)
				result.Result = append(result.Result, "Internal error")
				return result
			}

			currentSong.DeletedAt = lo.ToPtr(time.Now().UTC())
			parseCtx.Services.Gorm.WithContext(ctx).Updates(currentSong)
			parseCtx.Services.Redis.Del(ctx, redisKey)

			result.Result = append(result.Result, fmt.Sprintf("Song %s skipped", currentSong.Title))
			return result
		}

		result.Result = append(result.Result, fmt.Sprintf("%v/%v", votesCount+1, neededVotes))
		return result
	},
}
