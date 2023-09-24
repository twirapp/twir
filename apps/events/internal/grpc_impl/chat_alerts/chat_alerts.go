package chat_alerts

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/satont/twir/apps/events/internal"
	model "github.com/satont/twir/libs/gomodels"
)

type ChatAlerts struct {
	services *internal.Services
	settings *model.ChatAlertsSettings
}

func New(channelId string, services *internal.Services) (*ChatAlerts, error) {
	entity := model.ChannelModulesSettings{}

	if err := services.DB.Where(
		`"channelId" = ? AND "userId" IS NULL AND type = 'chat_alerts'`,
		channelId,
	).First(&entity).Error; err != nil {
		return nil, err
	}

	parsedSettings := model.ChatAlertsSettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, err
	}

	return &ChatAlerts{
		services: services,
		settings: &parsedSettings,
	}, nil
}

func (c *ChatAlerts) buildRedisCooldownKey(channelId, eventName string) string {
	return fmt.Sprintf(
		"channels:%s:chat_alerts_events:cooldowns:%s",
		channelId,
		eventName,
	)
}

func (c *ChatAlerts) IsOnCooldown(ctx context.Context, channelId, eventName string) (
	bool,
	error,
) {
	exists, err := c.services.Redis.Exists(
		ctx,
		c.buildRedisCooldownKey(channelId, eventName),
	).Result()
	return exists == 1, err
}

func (c *ChatAlerts) SetCooldown(
	ctx context.Context,
	channelId, eventName string,
	seconds int,
) error {
	if seconds == 0 {
		return nil
	}

	return c.services.Redis.Set(
		ctx,
		c.buildRedisCooldownKey(channelId, eventName),
		"",
		time.Duration(seconds)*time.Second,
	).Err()
}
