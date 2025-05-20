package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
	"gorm.io/gorm/clause"
)

// order metters
// user_id should be the last
var conditionKeys = []string{
	"broadcaster_user_id",
	"broadcaster_user_id",
	"to_broadcaster_user_id",
	"user_id",
}

var knownTopicsEntitiesCache = map[string]model.EventsubTopic{}

func (c *Handler) onNotification(
	_ *eventsub_bindings.ResponseHeaders,
	notification *eventsub_bindings.EventNotification,
) {
	condition, ok := notification.Subscription.Condition.(map[string]any)
	if !ok {
		c.logger.Error(
			"failed to cast condition",
			slog.Any("condition", notification.Subscription.Condition),
		)
		return
	}

	var userId string
	for _, key := range conditionKeys {
		if val, ok := condition[key].(string); ok {
			userId = val
			break
		}
	}

	if userId == "" {
		c.logger.Error("failed to find user_id")
		return
	}

	redisKey := fmt.Sprintf(
		"eventsub:cache:notification:check:%s:%s",
		notification.Subscription.Type,
		userId,
	)
	if exists, err := c.redisClient.Exists(context.Background(), redisKey).Result(); err != nil {
		c.logger.Error("failed to check redis", slog.Any("err", err))
		return
	} else if exists == 1 {
		return
	}

	var topicId uuid.UUID
	if cachedTopic, topicFound := knownTopicsEntitiesCache[notification.Subscription.Type]; topicFound {
		topicId = cachedTopic.ID
	} else {
		topicEntity := model.EventsubTopic{}
		if err := c.gorm.
			Where("topic = ?", notification.Subscription.Type).
			First(&topicEntity).
			Error; err != nil {
			c.logger.Error("failed to find topic", slog.Any("err", err))
			return
		}
		knownTopicsEntitiesCache[notification.Subscription.Type] = topicEntity
		topicId = topicEntity.ID
	}

	if topicId == uuid.Nil {
		c.logger.Error("failed to find topic_id")
		return
	}

	if err := c.gorm.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "topic_id"}, {Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"status": notification.Subscription.Status}),
		},
	).Create(
		&model.EventsubSubscription{
			ID:          uuid.New(),
			UserID:      userId,
			TopicID:     topicId,
			Status:      notification.Subscription.Status,
			Version:     notification.Subscription.Version,
			CallbackUrl: notification.Subscription.Transport.Callback,
		},
	).Error; err != nil {
		c.logger.Error("failed to create/update subscription", slog.Any("err", err))
	}

	c.redisClient.Set(context.Background(), redisKey, "true", 1*time.Minute)
}
