package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *MessageHandler) ensureUser(ctx context.Context, msg handleMessage) (*model.Users, error) {
	var user *model.Users

	redisCacheKey := fmt.Sprintf(
		"bots:cache:ensureuser:%s:%s",
		msg.BroadcasterUserId,
		msg.ChatterUserId,
	)
	if cachedBytes, _ := c.redis.Get(ctx, redisCacheKey).Bytes(); len(cachedBytes) > 0 {
		if err := json.Unmarshal(cachedBytes, &user); err != nil {
			return nil, err
		}
	}

	if user == nil {
		if err := c.gorm.WithContext(ctx).
			Where(
				"id = ?",
				msg.ChatterUserId,
			).
			Preload(
				"Stats",
				`"userId" = ? AND "channelId" = ?`,
				msg.ChatterUserId,
				msg.BroadcasterUserId,
			).
			First(&user).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			user.ID = msg.ChatterUserId
			user.ApiKey = uuid.NewString()
			user.Stats = &model.UsersStats{
				ID:                uuid.NewString(),
				UserID:            msg.ChatterUserId,
				ChannelID:         msg.BroadcasterUserId,
				Messages:          0,
				Watched:           0,
				UsedChannelPoints: 0,
				IsMod:             false,
				IsVip:             false,
				IsSubscriber:      false,
				Reputation:        0,
				Emotes:            0,
			}

			err := c.gorm.WithContext(ctx).Create(&user).Error
			if err != nil {
				return nil, err
			}
		}
	}

	badges := createUserBadges(msg.Badges)
	isMod := lo.Contains(badges, "MODERATOR")
	isSubscriber := lo.Contains(badges, "SUBSCRIBER")
	isVip := lo.Contains(badges, "VIP")

	if user.Stats == nil {
		user.Stats = &model.UsersStats{
			ID:                uuid.NewString(),
			UserID:            msg.ChatterUserId,
			ChannelID:         msg.BroadcasterUserId,
			Messages:          0,
			Watched:           0,
			UsedChannelPoints: 0,
			Reputation:        0,
			Emotes:            0,
		}
	} else if msg.EnrichedData.ChannelStream != nil {
		user.Stats.Messages += 1
	}

	usedEmotesInMessage := 0
	for _, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
		usedEmotesInMessage += count
	}

	user.Stats.Emotes += usedEmotesInMessage

	user.Stats.IsMod = isMod
	user.Stats.IsVip = isVip
	user.Stats.IsSubscriber = isSubscriber

	if err := c.gorm.WithContext(ctx).Save(&user.Stats).Error; err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(ctx, redisCacheKey, bytes, 24*time.Hour).Err(); err != nil {
		return nil, err
	}

	return user, nil
}
