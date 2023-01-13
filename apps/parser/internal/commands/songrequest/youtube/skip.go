package sr_youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"time"
)

var SkipCommand = types.DefaultCommand{
	Command: types.Command{
		Name:               "voteskip",
		Description:        lo.ToPtr("Vote for skip command"),
		Permission:         "VIEWER",
		Visible:            false,
		Module:             lo.ToPtr("SONGREQUEST"),
		IsReply:            true,
		KeepResponsesOrder: lo.ToPtr(false),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		logger := do.MustInvoke[zap.Logger](di.Provider)
		db := do.MustInvoke[gorm.DB](di.Provider)
		redisClient := do.MustInvoke[redis.Client](di.Provider)
		websocketGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)

		result := &types.CommandsHandlerResult{}

		moduleSettings := &model.ChannelModulesSettings{}
		parsedSettings := &youtube.YouTubeSettings{}
		err := db.
			Where(`"channelId" = ? AND "type" = ?`, ctx.ChannelId, "youtube_song_requests").
			First(moduleSettings).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		if moduleSettings.ID == "" {
			result.Result = append(result.Result, "Song requests not enabled")
			return result
		}

		err = json.Unmarshal(moduleSettings.Settings, parsedSettings)
		if err != nil {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		if !*parsedSettings.Enabled {
			result.Result = append(result.Result, parsedSettings.Translations.NotEnabled)
			return result
		}

		currentSong := &model.RequestedSong{}
		err = db.
			Where(`"channelId" = ? AND "deletedAt" IS NULL`, ctx.ChannelId).
			Order(`"createdAt" desc`).
			Limit(1).
			Find(&currentSong).
			Error

		if err != nil {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		if currentSong.ID == "" {
			result.Result = append(result.Result, "Current song not found")
			return result
		}

		var onlineUsersCount int64
		err = db.
			Where(`"channelId" = ?`, ctx.ChannelId).
			Model(&model.UsersOnline{}).
			Count(&onlineUsersCount).
			Error

		if err != nil {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		redisKey := fmt.Sprintf("songrequests-voteskip-%s", currentSong.ID)
		votesCount, err := redisClient.SCard(context.Background(), redisKey).Result()
		if err != nil {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		currentVote, err := redisClient.SIsMember(context.Background(), redisKey, ctx.SenderId).Result()
		if err != nil {
			logger.Sugar().Error(err)
			result.Result = append(result.Result, "Internal error")
			return result
		}

		neededVotes := int64(math.Round(parsedSettings.NeededVotesVorSkip * float64(onlineUsersCount) / 100))

		if currentVote {
			result.Result = append(result.Result, fmt.Sprintf("%v/%v", votesCount, neededVotes))
			return result
		}

		redisClient.SAdd(context.Background(), redisKey, ctx.SenderId)
		redisClient.Expire(context.Background(), redisKey, 1*time.Hour)

		if votesCount+1 >= neededVotes {
			_, err = websocketGrpc.YoutubeRemoveSongToQueue(context.Background(), &websockets.YoutubeRemoveSongFromQueueRequest{
				ChannelId: ctx.ChannelId,
				EntityId:  currentSong.ID,
			})

			if err != nil {
				logger.Sugar().Error(err)
				result.Result = append(result.Result, "Internal error")
				return result
			}

			currentSong.DeletedAt = lo.ToPtr(time.Now().UTC())
			db.Updates(currentSong)
			redisClient.Del(context.Background(), redisKey)

			result.Result = append(result.Result, fmt.Sprintf("Song %s skipped", currentSong.Title))
			return result
		}

		result.Result = append(result.Result, fmt.Sprintf("%v/%v", votesCount+1, neededVotes))
		return result
	},
}
