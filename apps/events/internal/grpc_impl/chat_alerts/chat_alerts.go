package chat_alerts

import (
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
