package games

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

const (
	votebanArgName = "user"
)

var Voteban = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "voteban",
		Description: null.StringFrom("Initiate voteban"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: votebanArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		mu := parseCtx.Services.RedSync.NewMutex(
			"parser:voteban:"+parseCtx.Channel.ID,
			redsync.WithExpiry(5*time.Second),
		)
		if err := mu.Lock(); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot lock voteban",
				Err:     err,
			}
		}
		defer mu.Unlock()

		entity := model.ChannelGamesVoteBan{}
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Preload("Channel").
			Where(
				`"channel_id" = ?`,
				parseCtx.Channel.ID,
			).First(&entity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			return nil, &types.CommandHandlerError{
				Message: "cannot find voteban settings",
				Err:     err,
			}
		}

		if !entity.Enabled {
			return nil, nil
		}

		if len(parseCtx.Mentions) == 0 {
			return nil, nil
		}

		targetUser := parseCtx.Mentions[0]

		if entity.Channel == nil ||
			targetUser.UserId == parseCtx.Channel.ID ||
			targetUser.UserId == entity.Channel.BotID {
			return nil, nil
		}

		redisKey := fmt.Sprintf("channels:%s:games:voteban", parseCtx.Channel.ID)
		if entity.VotingMode == model.ChannelGamesVoteBanVotingModeChat {
			voteInProgress, err := parseCtx.Services.Redis.Exists(ctx, redisKey).Result()
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot check if vote in progress",
					Err:     err,
				}
			}

			if voteInProgress == 1 {
				return &types.CommandsHandlerResult{
					Result: []string{"Another voteban in progress"},
				}, nil
			}

			targetUserStatsEntity := model.UsersStats{}
			if err := parseCtx.Services.Gorm.
				WithContext(ctx).
				Where(`"userId" = ? AND "channelId" = ?`, targetUser.UserId, parseCtx.Channel.ID).
				First(&targetUserStatsEntity).Error; err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot find target user",
					Err:     err,
				}
			}

			if err := parseCtx.Services.Redis.HSet(
				ctx,
				redisKey,
				model.ChannelGamesVoteBanRedisStruct{
					TargetUserId:   targetUser.UserId,
					TargetUserName: targetUser.UserName,
					TargetIsMod:    targetUserStatsEntity.IsMod,
					TotalVotes:     1,
					PositiveVotes:  1,
					NegativeVotes:  0,
				},
			).Err(); err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot set vote",
					Err:     err,
				}
			}

			if err := parseCtx.Services.Redis.Expire(
				ctx,
				redisKey,
				time.Second*time.Duration(entity.VoteDuration),
			).Err(); err != nil {
				parseCtx.Services.Redis.Del(ctx, redisKey)

				return nil, &types.CommandHandlerError{
					Message: "cannot set vote expiration",
					Err:     err,
				}
			}

			if err := parseCtx.Services.Redis.Set(
				ctx,
				fmt.Sprintf("%s:votes:%s", redisKey, parseCtx.Sender.ID),
				1,
				time.Second*time.Duration(entity.VoteDuration),
			).Err(); err != nil {
				parseCtx.Services.Redis.Del(ctx, redisKey)
			}

			initMessage := strings.ReplaceAll(
				entity.InitMessage,
				"{targetUser}",
				targetUser.UserName,
			)
			initMessage = strings.ReplaceAll(
				initMessage,
				"{positiveTexts}",
				strings.Join(entity.ChatVotesWordsPositive, " · "),
			)
			initMessage = strings.ReplaceAll(
				initMessage,
				"{negativeTexts}",
				strings.Join(entity.ChatVotesWordsNegative, " · "),
			)

			return &types.CommandsHandlerResult{
				Result: []string{initMessage},
			}, nil
		}

		return &types.CommandsHandlerResult{
			Result: []string{},
		}, nil
	},
}
