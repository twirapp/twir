package modules

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/modules_chat_alerts"
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
			Cooldown: int32(entity.Followers.Cooldown),
		},
		Raids: &modules_chat_alerts.ChatAlertsRaids{
			Enabled:  entity.Raids.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Raids.Messages),
			Cooldown: int32(entity.Raids.Cooldown),
		},
		Donations: &modules_chat_alerts.ChatAlertsDonations{
			Enabled:  entity.Donations.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Donations.Messages),
			Cooldown: int32(entity.Donations.Cooldown),
		},
		Subscribers: &modules_chat_alerts.ChatAlertsSubscribers{
			Enabled:  entity.Subscribers.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Subscribers.Messages),
			Cooldown: int32(entity.Subscribers.Cooldown),
		},
		Cheers: &modules_chat_alerts.ChatAlertsCheers{
			Enabled:  entity.Cheers.Enabled,
			Messages: c.convertChatAlertsCountedMessages(entity.Cheers.Messages),
			Cooldown: int32(entity.Cheers.Cooldown),
		},
		Redemptions: &modules_chat_alerts.ChatAlertsRedemptions{
			Enabled:  entity.Redemptions.Enabled,
			Messages: c.convertChatAlertsMessages(entity.Redemptions.Messages),
			Cooldown: int32(entity.Redemptions.Cooldown),
		},
		FirstUserMessage: &modules_chat_alerts.ChatAlertsFirstUserMessage{
			Enabled:  entity.FirstUserMessage.Enabled,
			Messages: c.convertChatAlertsMessages(entity.FirstUserMessage.Messages),
			Cooldown: int32(entity.FirstUserMessage.Cooldown),
		},
		StreamOnline: &modules_chat_alerts.ChatAlertsStreamOnline{
			Enabled:  entity.StreamOnline.Enabled,
			Messages: c.convertChatAlertsMessages(entity.StreamOnline.Messages),
			Cooldown: int32(entity.StreamOnline.Cooldown),
		},
		StreamOffline: &modules_chat_alerts.ChatAlertsStreamOffline{
			Enabled:  entity.StreamOffline.Enabled,
			Messages: c.convertChatAlertsMessages(entity.StreamOffline.Messages),
			Cooldown: int32(entity.StreamOffline.Cooldown),
		},
		ChatCleared: &modules_chat_alerts.ChatAlertsChatCleared{
			Enabled:  entity.ChatCleared.Enabled,
			Messages: c.convertChatAlertsMessages(entity.ChatCleared.Messages),
			Cooldown: int32(entity.ChatCleared.Cooldown),
		},
		Ban: &modules_chat_alerts.ChatAlertsBan{
			Enabled:           entity.Ban.Enabled,
			Messages:          c.convertChatAlertsCountedMessages(entity.Ban.Messages),
			Cooldown:          int32(entity.Ban.Cooldown),
			IgnoreTimeoutFrom: entity.Ban.IgnoreTimeoutFrom,
		},
		ChannelUnbanRequestCreate: &modules_chat_alerts.ChatAlertsChannelUnbanRequestCreate{
			Enabled:  entity.UnbanRequestCreate.Enabled,
			Messages: c.convertChatAlertsMessages(entity.UnbanRequestCreate.Messages),
			Cooldown: int32(entity.UnbanRequestCreate.Cooldown),
		},
		ChannelUnbanRequestResolve: &modules_chat_alerts.ChatAlertsChannelUnbanRequestResolve{
			Enabled:  entity.UnbanRequestResolve.Enabled,
			Messages: c.convertChatAlertsMessages(entity.UnbanRequestResolve.Messages),
			Cooldown: int32(entity.UnbanRequestResolve.Cooldown),
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
			Cooldown: int(req.Followers.Cooldown),
		},
		Raids: model.ChatAlertsRaids{
			Enabled:  req.Raids.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Raids.Messages),
			Cooldown: int(req.Raids.Cooldown),
		},
		Donations: model.ChatAlertsDonations{
			Enabled:  req.Donations.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Donations.Messages),
			Cooldown: int(req.Donations.Cooldown),
		},
		Subscribers: model.ChatAlertsSubscribers{
			Enabled:  req.Subscribers.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Subscribers.Messages),
			Cooldown: int(req.Subscribers.Cooldown),
		},
		Cheers: model.ChatAlertsCheers{
			Enabled:  req.Cheers.Enabled,
			Messages: c.chatAlertsRequestedCountedToDb(req.Cheers.Messages),
			Cooldown: int(req.Cheers.Cooldown),
		},
		Redemptions: model.ChatAlertsRedemptions{
			Enabled:  req.Redemptions.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.Redemptions.Messages),
			Cooldown: int(req.Redemptions.Cooldown),
		},
		FirstUserMessage: model.ChatAlertsFirstUserMessage{
			Enabled:  req.FirstUserMessage.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.FirstUserMessage.Messages),
			Cooldown: int(req.FirstUserMessage.Cooldown),
		},
		StreamOnline: model.ChatAlertsStreamOnline{
			Enabled:  req.StreamOnline.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.StreamOnline.Messages),
			Cooldown: int(req.StreamOnline.Cooldown),
		},
		StreamOffline: model.ChatAlertsStreamOffline{
			Enabled:  req.StreamOffline.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.StreamOffline.Messages),
			Cooldown: int(req.StreamOffline.Cooldown),
		},
		ChatCleared: model.ChatAlertsChatCleared{
			Enabled:  req.ChatCleared.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.ChatCleared.Messages),
			Cooldown: int(req.ChatCleared.Cooldown),
		},
		Ban: model.ChatAlertsBan{
			Enabled:           req.Ban.Enabled,
			Messages:          c.chatAlertsRequestedCountedToDb(req.Ban.Messages),
			Cooldown:          int(req.Ban.Cooldown),
			IgnoreTimeoutFrom: req.Ban.IgnoreTimeoutFrom,
		},
		UnbanRequestCreate: model.ChatAlertsUnbanRequestCreate{
			Enabled:  req.ChannelUnbanRequestCreate.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.ChannelUnbanRequestCreate.Messages),
			Cooldown: int(req.ChannelUnbanRequestCreate.Cooldown),
		},
		UnbanRequestResolve: model.ChatAlertsUnbanRequestResolve{
			Enabled:  req.ChannelUnbanRequestResolve.Enabled,
			Messages: c.chatAlertsRequestedToDb(req.ChannelUnbanRequestResolve.Messages),
			Cooldown: int(req.ChannelUnbanRequestResolve.Cooldown),
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
