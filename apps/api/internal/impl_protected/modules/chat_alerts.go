package modules

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/modules_chat_alerts"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Modules) convertChatAlertsCountedMessages(
	entity []model.ChatAlertsCountedMessage,
) []*modules_chat_alerts.ChatAlertsCountedMessage {
	result := make([]*modules_chat_alerts.ChatAlertsCountedMessage, len(entity))
	for i, e := range entity {
		result[i] = &modules_chat_alerts.ChatAlertsCountedMessage{
			Count: int32(e.Count),
			Text:  e.Text,
		}
	}

	return result
}

func (c *Modules) convertChatAlertsMessages(
	entity []model.ChatAlertsMessage,
) []*modules_chat_alerts.ChatAlertsMessage {
	result := make([]*modules_chat_alerts.ChatAlertsMessage, len(entity))
	for i, e := range entity {
		result[i] = &modules_chat_alerts.ChatAlertsMessage{
			Text: e.Text,
		}
	}

	return result
}

func (c *Modules) convertChatAlertsSettings(
	entity model.ChatAlertsSettings,
) *modules_chat_alerts.ChatAlertsSettings {
	return &modules_chat_alerts.ChatAlertsSettings{
		Followers: &modules_chat_alerts.ChatAlertsFollowersSettings{
			Enabled:  entity.Followers.Enabled,
			Messages: c.convertChatAlertsMessages(entity.Followers.Messages),
		},
		Raids: &modules_chat_alerts.ChatAlertsRaids{
			Enabled:  entity.Raids.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Raids.Messages),
		},
		Donations: &modules_chat_alerts.ChatAlertsDonations{
			Enabled:  entity.Donations.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Donations.Messages),
		},
		Subscribers: &modules_chat_alerts.ChatAlertsSubscribers{
			Enabled:  entity.Subscribers.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Subscribers.Messages),
		},
		Cheers: &modules_chat_alerts.ChatAlertsCheers{
			Enabled:  entity.Cheers.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Cheers.Messages),
		},
		Redemptions: &modules_chat_alerts.ChatAlertsRedemptions{
			Enabled:  entity.Redemptions.Enabled,
			Messages: c.convertChatAlertsMessages(entity.Redemptions.Messages),
		},
		FirstUserMessage: &modules_chat_alerts.ChatAlertsFirstUserMessage{
			Enabled:  entity.FirstUserMessage.Enabled,
			Messages: c.convertChatAlertsMessages(entity.FirstUserMessage.Messages),
		},
		StreamOnline: &modules_chat_alerts.ChatAlertsStreamOnline{
			Enabled:  entity.StreamOnline.Enabled,
			Messages: c.convertChatAlertsMessages(entity.StreamOnline.Messages),
		},
		StreamOffline: &modules_chat_alerts.ChatAlertsStreamOffline{
			Enabled:  entity.StreamOffline.Enabled,
			Messages: c.convertChatAlertsMessages(entity.StreamOffline.Messages),
		},
		ChatCleared: &modules_chat_alerts.ChatAlertsChatCleared{
			Enabled:  entity.ChatCleared.Enabled,
			Messages: c.convertChatAlertsMessages(entity.ChatCleared.Messages),
		},
		Ban: &modules_chat_alerts.ChatAlertsBan{
			Enabled:           entity.Ban.Enabled,
			Messages:          c.convertChatAlertsCountedMessages(entity.Ban.Messages),
			IgnoreTimeoutFrom: entity.Ban.IgnoreTimeoutFrom,
		},
	}
}

func (c *Modules) ModulesChatAlertsGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*modules_chat_alerts.ChatAlertsSettings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "userId" IS NULL AND type = 'chat_alerts'`,
			dashboardId,
		).First(&entity).Error; err != nil {
		return nil, err
	}

	parsedSettings := model.ChatAlertsSettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, err
	}

	return c.convertChatAlertsSettings(parsedSettings), nil
}

func (c *Modules) chatAlertsRequestedCountedToDb(
	req []*modules_chat_alerts.ChatAlertsCountedMessage,
) []model.ChatAlertsCountedMessage {
	result := make([]model.ChatAlertsCountedMessage, len(req))

	for i, e := range req {
		result[i] = model.ChatAlertsCountedMessage{
			Count: int(e.Count),
			Text:  e.Text,
		}
	}

	return result
}

func (c *Modules) chatAlertsRequestedToDb(req []*modules_chat_alerts.ChatAlertsMessage) []model.ChatAlertsMessage {
	result := make([]model.ChatAlertsMessage, len(req))

	for i, e := range req {
		result[i] = model.ChatAlertsMessage{
			Text: e.Text,
		}
	}

	return result
}

func (c *Modules) ModulesChatAlertsUpdate(
	ctx context.Context,
	req *modules_chat_alerts.ChatAlertsSettings,
) (*modules_chat_alerts.ChatAlertsSettings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "userId" IS NULL AND type = 'chat_alerts'`,
			dashboardId,
		).Find(&entity).Error; err != nil {
		return nil, err
	}

	if entity.ID == "" || entity.ChannelId == "" {
		entity.ID = uuid.New().String()
		entity.ChannelId = dashboardId
	}

	entity.Type = "chat_alerts"

	newSettings := model.ChatAlertsSettings{
		Followers: model.ChatAlertsFollowersSettings{
			Enabled:  req.Followers.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.Followers.Messages),
		},
		Raids: model.ChatAlertsRaids{
			Enabled:  req.Raids.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Raids.Messages),
		},
		Donations: model.ChatAlertsDonations{
			Enabled:  req.Donations.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Donations.Messages),
		},
		Subscribers: model.ChatAlertsSubscribers{
			Enabled:  req.Subscribers.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Subscribers.Messages),
		},
		Cheers: model.ChatAlertsCheers{
			Enabled:  req.Cheers.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Cheers.Messages),
		},
		Redemptions: model.ChatAlertsRedemptions{
			Enabled:  req.Redemptions.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.Redemptions.Messages),
		},
		FirstUserMessage: model.ChatAlertsFirstUserMessage{
			Enabled:  req.FirstUserMessage.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.FirstUserMessage.Messages),
		},
		StreamOnline: model.ChatAlertsStreamOnline{
			Enabled:  req.StreamOnline.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.StreamOnline.Messages),
		},
		StreamOffline: model.ChatAlertsStreamOffline{
			Enabled:  req.StreamOffline.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.StreamOffline.Messages),
		},
		ChatCleared: model.ChatAlertsChatCleared{
			Enabled:  req.ChatCleared.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.ChatCleared.Messages),
		},
		Ban: model.ChatAlertsBan{
			Enabled:           req.Ban.Enabled,
			Messages:          c.chatAlertsRequestedCountedToDb(req.Ban.Messages),
			IgnoreTimeoutFrom: req.Ban.IgnoreTimeoutFrom,
		},
	}

	bytes, err := json.Marshal(newSettings)
	if err != nil {
		return nil, err
	}

	entity.Settings = bytes

	if err := c.Db.
		WithContext(ctx).
		Save(&entity).Error; err != nil {
		return nil, err
	}

	return c.convertChatAlertsSettings(newSettings), nil
}
