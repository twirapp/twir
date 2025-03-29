package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/redis_keys"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
)

func (c *MessageHandler) handleChatWall(ctx context.Context, msg handleMessage) error {
	if msg.Message == nil ||
		msg.ChatterUserId == msg.DbChannel.ID ||
		msg.ChatterUserId == msg.DbChannel.BotID {
		return nil
	}

	badges := createUserBadges(msg.Badges)
	for _, b := range badges {
		if slices.Contains(excludedModerationBadges, b) {
			return nil
		}
	}

	walls, err := c.chatWallCacher.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
	}

	wallSettings, err := c.chatWallSettingsCacher.Get(ctx, msg.BroadcasterUserId)
	if err != nil && !errors.Is(err, chatwallrepository.ErrSettingsNotFound) {
		return err
	}

	var shouldInvalidate bool

	for _, wall := range walls {
		if !wall.Enabled || wall.Phrase == "" {
			continue
		}

		if wall.DurationSeconds != 0 && wall.CreatedAt.Before(time.Now().Add(-time.Duration(wall.DurationSeconds)*time.Second)) {
			if _, err := c.chatWallRepository.Update(
				ctx,
				wall.ID,
				chatwallrepository.UpdateInput{
					Enabled: lo.ToPtr(false),
				},
			); err != nil {
				return err
			}
			shouldInvalidate = true

			continue
		}

		msgLowerCased := strings.ToLower(msg.Message.Text)
		phraseLowerCased := strings.ToLower(wall.Phrase)

		if !strings.Contains(msgLowerCased, phraseLowerCased) {
			continue
		}

		if wallSettings.ChannelID != "" {
			if !wallSettings.MuteSubscribers &&
				(slices.Contains(badges, "SUBSCRIBER") || slices.Contains(badges, "FOUNDER")) {
				continue
			}
			if !wallSettings.MuteVips && slices.Contains(badges, "VIP") {
				continue
			}
		}

		alreadyHandled, err := c.redis.SIsMember(
			ctx,
			fmt.Sprintf(redis_keys.NukeRedisPrefix, msg.BroadcasterUserId),
			msg.ID,
		).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}
		if alreadyHandled {
			continue
		}

		switch wall.Action {
		case chatwallmodel.ChatWallActionDelete:
			if err := c.twitchActions.DeleteMessage(
				ctx,
				twitchactions.DeleteMessageOpts{
					BroadcasterID: msg.BroadcasterUserId,
					ModeratorID:   msg.DbChannel.BotID,
					MessageID:     msg.ID,
				},
			); err != nil {
				return err
			}
		case chatwallmodel.ChatWallActionBan, chatwallmodel.ChatWallActionTimeout:
			var duration int
			if wall.Action == chatwallmodel.ChatWallActionTimeout && wall.TimeoutDurationSeconds != nil {
				duration = *wall.TimeoutDurationSeconds
			}

			if err := c.twitchActions.Ban(
				ctx,
				twitchactions.BanOpts{
					Duration:      duration,
					Reason:        "banned by twir for chat wall phrase: " + wall.Phrase,
					BroadcasterID: msg.BroadcasterUserId,
					UserID:        msg.ChatterUserId,
					ModeratorID:   msg.DbChannel.BotID,
				},
			); err != nil {
				return err
			}
		}

		if err := c.chatWallRepository.CreateLog(
			ctx,
			chatwallrepository.CreateLogInput{
				WallID: wall.ID,
				UserID: msg.ChatterUserId,
				Text:   msg.Message.Text,
			},
		); err != nil {
			return err
		}

		_, err = c.redis.Pipelined(
			ctx, func(p redis.Pipeliner) error {
				if err := p.SAdd(
					ctx,
					fmt.Sprintf(redis_keys.NukeRedisPrefix, msg.BroadcasterUserId),
					msg.ID,
				).Err(); err != nil {
					return err
				}
				if err := p.Expire(
					ctx,
					fmt.Sprintf(redis_keys.NukeRedisPrefix, msg.BroadcasterUserId),
					20*time.Minute,
				).Err(); err != nil {
					return err
				}

				return nil
			},
		)
		if err != nil {
			return fmt.Errorf("cannot add handled messages to redis: %s", err)
		}

		break
	}

	if shouldInvalidate {
		c.chatWallCacher.Invalidate(ctx, msg.BroadcasterUserId)
	}

	return nil
}
