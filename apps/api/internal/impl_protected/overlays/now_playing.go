package overlays

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/overlays"
	"github.com/twirapp/twir/libs/api/messages/overlays_now_playing"
	"google.golang.org/protobuf/types/known/emptypb"
)

func convertEntityToProto(entity model.ChannelOverlayNowPlaying) *overlays_now_playing.Settings {
	return &overlays_now_playing.Settings{
		Id:        entity.ID.String(),
		Preset:    entity.Preset.String(),
		ChannelId: entity.ChannelID,
	}
}

func (c *Overlays) OverlaysNowPlayingGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_now_playing.GetAllResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelOverlayNowPlaying
	if err := c.Db.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			dashboardId,
		).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	settings := make([]*overlays_now_playing.Settings, 0, len(entities))
	for _, overlay := range entities {
		settings = append(settings, convertEntityToProto(overlay))
	}

	return &overlays_now_playing.GetAllResponse{
		Settings: settings,
	}, nil
}

func (c *Overlays) OverlaysNowPlayingUpdate(
	ctx context.Context,
	req *overlays_now_playing.UpdateRequest,
) (*overlays_now_playing.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	overlay := model.ChannelOverlayNowPlaying{}
	if err := c.Db.
		WithContext(ctx).
		Where(
			"id = ? AND channel_id = ?",
			dashboardId,
			req.GetId(),
		).
		First(&overlay).Error; err != nil {
		return nil, err
	}

	overlay.Preset = overlays.ChannelOverlayNowPlayingPreset(req.GetPreset())
	if err := c.Db.
		WithContext(ctx).
		Save(&overlay).Error; err != nil {
		return nil, err
	}

	return convertEntityToProto(overlay), nil
}

func (c *Overlays) OverlaysNowPlayingDelete(
	ctx context.Context,
	_ *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.Db.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			dashboardId,
		).
		Delete(&model.ChannelOverlayNowPlaying{}).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Overlays) OverlaysNowPlayingCreate(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_now_playing.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	overlay := model.ChannelOverlayNowPlaying{
		ID:        uuid.New(),
		ChannelID: dashboardId,
		Preset:    overlays.ChannelOverlayNowPlayingPresetTransparent,
	}

	if err := c.Db.
		WithContext(ctx).
		Create(&overlay).Error; err != nil {
		return nil, err
	}

	return convertEntityToProto(overlay), nil
}
