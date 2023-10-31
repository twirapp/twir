package modules

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/modules_chat_overlay"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

const chatOverlayType = "chat_overlay"

func (c *Modules) chatOverlayDbToGrpc(s model.ChatOverlaySettings) *modules_chat_overlay.Settings {
	return &modules_chat_overlay.Settings{
		MessageHideTimeout: s.MessageHideTimeout,
		MessageShowDelay:   s.MessageShowDelay,
		Preset:             s.Preset,
		FontSize:           s.FontSize,
		HideCommands:       s.HideCommands,
		HideBots:           s.HideBots,
		FontFamily:         s.FontFamily,
		ShowBadges:         s.ShowBadges,
		ShowAnnounceBadge:  s.ShowAnnounceBadge,
	}
}

func (c *Modules) chatOverlayGrpcToDb(s *modules_chat_overlay.Settings) model.ChatOverlaySettings {
	return model.ChatOverlaySettings{
		MessageHideTimeout: s.MessageHideTimeout,
		MessageShowDelay:   s.MessageShowDelay,
		Preset:             s.Preset,
		FontSize:           s.FontSize,
		HideCommands:       s.HideCommands,
		HideBots:           s.HideBots,
		FontFamily:         s.FontFamily,
		ShowBadges:         s.ShowBadges,
		ShowAnnounceBadge:  s.ShowAnnounceBadge,
	}
}

func (c *Modules) ModulesChatOverlayGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*modules_chat_overlay.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, chatOverlayType).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	parsedSettings := model.ChatOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	return c.chatOverlayDbToGrpc(parsedSettings), nil
}

func (c *Modules) ModulesChatOverlayUpdate(
	ctx context.Context,
	req *modules_chat_overlay.Settings,
) (*modules_chat_overlay.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, chatOverlayType).
		Find(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	if entity.ID == "" {
		entity.ID = uuid.NewString()
		entity.ChannelId = dashboardId
		entity.Type = chatOverlayType
	}

	parsedSettings := c.chatOverlayGrpcToDb(req)
	settingsJson, err := json.Marshal(parsedSettings)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal settings: %w", err)
	}
	entity.Settings = settingsJson
	if err := c.Db.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot update settings: %w", err)
	}

	newSettings := model.ChatOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &newSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	c.Grpc.Websockets.RefreshChatSettings(
		ctx, &websockets.RefreshChatSettingsRequest{
			ChannelId: dashboardId,
		},
	)

	return c.chatOverlayDbToGrpc(newSettings), nil
}
