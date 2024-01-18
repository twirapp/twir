package overlays

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/overlays_chat"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Overlays) chatOverlayDbToGrpc(s model.ChatOverlaySettings) *overlays_chat.Settings {
	id := s.ID.String()

	return &overlays_chat.Settings{
		Id:                  &id,
		MessageHideTimeout:  s.MessageHideTimeout,
		MessageShowDelay:    s.MessageShowDelay,
		Preset:              s.Preset,
		FontFamily:          s.FontFamily,
		FontSize:            s.FontSize,
		FontWeight:          s.FontWeight,
		FontStyle:           s.FontStyle,
		HideCommands:        s.HideCommands,
		HideBots:            s.HideBots,
		ShowBadges:          s.ShowBadges,
		ShowAnnounceBadge:   s.ShowAnnounceBadge,
		TextShadowColor:     s.TextShadowColor,
		TextShadowSize:      s.TextShadowSize,
		ChatBackgroundColor: s.ChatBackgroundColor,
		Direction:           s.Direction,
		PaddingContainer:    s.PaddingContainer,
	}
}

func (c *Overlays) OverlayChatGet(
	ctx context.Context,
	req *overlays_chat.GetRequest,
) (*overlays_chat.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChatOverlaySettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetId()).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	return c.chatOverlayDbToGrpc(entity), nil
}

func (c *Overlays) OverlayChatCreate(
	ctx context.Context,
	req *overlays_chat.Settings,
) (*overlays_chat.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChatOverlaySettings{
		ID:                  uuid.New(),
		MessageHideTimeout:  req.GetMessageHideTimeout(),
		MessageShowDelay:    req.GetMessageShowDelay(),
		Preset:              req.GetPreset(),
		FontFamily:          req.GetFontFamily(),
		FontSize:            req.GetFontSize(),
		FontWeight:          req.GetFontWeight(),
		FontStyle:           req.GetFontStyle(),
		HideCommands:        req.GetHideCommands(),
		HideBots:            req.GetHideBots(),
		ShowBadges:          req.GetShowBadges(),
		ShowAnnounceBadge:   req.GetShowAnnounceBadge(),
		TextShadowColor:     req.GetTextShadowColor(),
		TextShadowSize:      req.GetTextShadowSize(),
		ChatBackgroundColor: req.GetChatBackgroundColor(),
		Direction:           req.GetDirection(),
		PaddingContainer:    req.GetPaddingContainer(),
		ChannelID:           dashboardId,
	}

	if err := c.Db.
		WithContext(ctx).
		Create(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot create settings: %w", err)
	}

	return c.chatOverlayDbToGrpc(entity), nil
}

func (c *Overlays) OverlayChatUpdate(
	ctx context.Context,
	req *overlays_chat.UpdateRequest,
) (*overlays_chat.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	settings := req.GetSettings()
	if settings == nil {
		return nil, fmt.Errorf("settings are required")
	}

	if req.GetId() == "" {
		return nil, fmt.Errorf("id is required")
	}

	entity := model.ChatOverlaySettings{
		ID:        uuid.MustParse(req.GetId()),
		ChannelID: dashboardId,
	}

	if err := c.Db.
		WithContext(ctx).
		Find(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot create settings: %w", err)
	}

	entity.MessageHideTimeout = req.GetSettings().GetMessageHideTimeout()
	entity.MessageShowDelay = req.GetSettings().GetMessageShowDelay()
	entity.Preset = req.GetSettings().GetPreset()
	entity.FontFamily = req.GetSettings().GetFontFamily()
	entity.FontSize = req.GetSettings().GetFontSize()
	entity.FontWeight = req.GetSettings().GetFontWeight()
	entity.FontStyle = req.GetSettings().GetFontStyle()
	entity.HideCommands = req.GetSettings().GetHideCommands()
	entity.HideBots = req.GetSettings().GetHideBots()
	entity.ShowBadges = req.GetSettings().GetShowBadges()
	entity.ShowAnnounceBadge = req.GetSettings().GetShowAnnounceBadge()
	entity.TextShadowColor = req.GetSettings().GetTextShadowColor()
	entity.TextShadowSize = req.GetSettings().GetTextShadowSize()
	entity.ChatBackgroundColor = req.GetSettings().GetChatBackgroundColor()
	entity.Direction = req.GetSettings().GetDirection()
	entity.PaddingContainer = req.GetSettings().GetPaddingContainer()

	if err := c.Db.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, fmt.Errorf("cannot update settings: %w", err)
	}

	if _, err := c.Grpc.Websockets.RefreshChatOverlaySettings(
		ctx,
		&websockets.RefreshChatSettingsRequest{
			ChannelId: dashboardId,
			Id:        entity.ID.String(),
		},
	); err != nil {
		c.Logger.Error("cannot refresh chat overlay settings", slog.Any("err", err))
	}

	return c.chatOverlayDbToGrpc(entity), nil
}

func (c *Overlays) OverlayChatDelete(
	ctx context.Context,
	req *overlays_chat.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetId() == "" {
		return nil, fmt.Errorf("id is required")
	}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channel_id" = ? AND "id" = ?`, dashboardId, req.GetId()).
		Delete(&model.ChatOverlaySettings{}).
		Error; err != nil {
		return nil, fmt.Errorf("cannot delete settings: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (c *Overlays) OverlayChatGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_chat.GetAllResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChatOverlaySettings
	if err := c.Db.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		Order("created_at asc").
		Find(&entities).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	settings := make([]*overlays_chat.Settings, 0, len(entities))
	for _, entity := range entities {
		settings = append(settings, c.chatOverlayDbToGrpc(entity))
	}

	return &overlays_chat.GetAllResponse{Settings: settings}, nil
}
