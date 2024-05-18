package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/kr/pretty"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *MessageHandler) handleGamesVoteban(ctx context.Context, msg handleMessage) error {
	c.votebanMutex.Lock()
	defer c.votebanMutex.Unlock()

	redisKey := fmt.Sprintf("channels:%s:games:voteban", msg.BroadcasterUserId)

	voteExists, err := c.redis.Exists(ctx, redisKey).Result()
	if err != nil {
		return err
	}

	if voteExists == 0 {
		return nil
	}

	userVoteExists, err := c.redis.Exists(
		ctx,
		fmt.Sprintf("%s:totalVotes:%s", redisKey, msg.ChatterUserId),
	).Result()
	if err != nil {
		return err
	}

	if userVoteExists == 1 {
		return nil
	}

	var voteEntity model.ChannelGamesVoteBanRedisStruct
	err = c.redis.HGetAll(ctx, redisKey).Scan(&voteEntity)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	gameEntity := model.ChannelGamesVoteBan{}
	err = c.gorm.
		WithContext(ctx).
		Where(
			`"channel_id" = ?`,
			msg.BroadcasterUserId,
		).
		First(&gameEntity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if !gameEntity.Enabled {
		return nil
	}

	splittedChatMessage := strings.Fields(msg.Message.Text)

	for _, word := range splittedChatMessage {
		if slices.Contains(gameEntity.ChatVotesWordsPositive, word) {
			voteEntity.TotalVotes++
			voteEntity.PositiveVotes++
			break
		} else if slices.Contains(gameEntity.ChatVotesWordsNegative, word) {
			voteEntity.TotalVotes++
			voteEntity.NegativeVotes++
			break
		}
	}

	pretty.Println(voteEntity)

	if voteEntity.TotalVotes >= gameEntity.NeededVotes {
		// TODO: make logic for timeout/not timeout

		pretty.Println(
			"Voteban success",
			voteEntity,
		)

		if err := c.redis.Del(ctx, redisKey).Err(); err != nil {
			return err
		}
	} else {
		if err := c.redis.HSet(
			ctx,
			redisKey,
			voteEntity,
		).Err(); err != nil {
			return err
		}
	}

	return nil
}
