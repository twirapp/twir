package overlays

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/overlays_chat"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

const chatOverlayType = "chat_overlay"

func (c *Overlays) chatOverlayDbToGrpc(s model.ChatOverlaySettings) *overlays_chat.Settings {
	return &overlays_chat.Settings{
		MessageHideTimeout:  s.MessageHideTimeout,
		MessageShowDelay:    s.MessageShowDelay,
		Preset:              s.Preset,
		FontSize:            s.FontSize,
		HideCommands:        s.HideCommands,
		HideBots:            s.HideBots,
		FontFamily:          s.FontFamily,
		ShowBadges:          s.ShowBadges,
		ShowAnnounceBadge:   s.ShowAnnounceBadge,
		TextShadowColor:     s.TextShadowColor,
		TextShadowSize:      s.TextShadowSize,
		ChatBackgroundColor: s.ChatBackgroundColor,
		Direction:           s.Direction,
	}
}

func (c *Overlays) chatOverlayGrpcToDb(s *overlays_chat.Settings) model.ChatOverlaySettings {
	return model.ChatOverlaySettings{
		MessageHideTimeout:  s.MessageHideTimeout,
		MessageShowDelay:    s.MessageShowDelay,
		Preset:              s.Preset,
		FontSize:            s.FontSize,
		HideCommands:        s.HideCommands,
		HideBots:            s.HideBots,
		FontFamily:          s.FontFamily,
		ShowBadges:          s.ShowBadges,
		ShowAnnounceBadge:   s.ShowAnnounceBadge,
		TextShadowColor:     s.TextShadowColor,
		TextShadowSize:      s.TextShadowSize,
		ChatBackgroundColor: s.ChatBackgroundColor,
		Direction:           s.Direction,
	}
}

func (c *Overlays) OverlayChatGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_chat.Settings, error) {
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

func (c *Overlays) OverlayChatUpdate(
	ctx context.Context,
	req *overlays_chat.Settings,
) (*overlays_chat.Settings, error) {
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

	c.Grpc.Websockets.RefreshChatOverlaySettings(
		ctx, &websockets.RefreshChatSettingsRequest{
			ChannelId: dashboardId,
		},
	)

	return c.chatOverlayDbToGrpc(newSettings), nil
}
