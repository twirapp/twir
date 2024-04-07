package overlays

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/overlays"
	"github.com/twirapp/twir/libs/api/messages/overlays_now_playing"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func convertEntityToProto(entity model.ChannelOverlayNowPlaying) *overlays_now_playing.Settings {
	return &overlays_now_playing.Settings{
		Id:              entity.ID.String(),
		Preset:          entity.Preset.String(),
		ChannelId:       entity.ChannelID,
		FontFamily:      entity.FontFamily,
		FontWeight:      entity.FontWeight,
		BackgroundColor: entity.BackgroundColor,
		ShowImage:       entity.ShowImage,
		HideTimeout:     lo.ToPtr(int32(entity.HideTimeout.Int64)),
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
			req.GetId(),
			dashboardId,
		).
		First(&overlay).Error; err != nil {
		return nil, err
	}

	overlay.Preset = overlays.ChannelOverlayNowPlayingPreset(req.GetPreset())
	overlay.FontFamily = req.GetFontFamily()
	overlay.FontWeight = req.GetFontWeight()
	overlay.BackgroundColor = req.GetBackgroundColor()
	overlay.ShowImage = req.GetShowImage()
	overlay.HideTimeout = null.IntFrom(int64(req.GetHideTimeout()))

	if err := c.Db.
		WithContext(ctx).
		Save(&overlay).Error; err != nil {
		return nil, err
	}

	c.Grpc.Websockets.RefreshOverlaySettings(
		ctx,
		&websockets.RefreshOverlaysRequest{
			ChannelId:   dashboardId,
			OverlayName: websockets.RefreshOverlaySettingsName_NOW_PLAYING,
			OverlayId:   lo.ToPtr(overlay.ID.String()),
		},
	)

	return convertEntityToProto(overlay), nil
}

func (c *Overlays) OverlaysNowPlayingDelete(
	ctx context.Context,
	req *overlays_now_playing.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.Db.
		WithContext(ctx).
		Where(
			"channel_id = ? AND id = ?",
			dashboardId,
			req.GetId(),
		).
		Delete(&model.ChannelOverlayNowPlaying{}).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Overlays) OverlaysNowPlayingCreate(
	ctx context.Context,
	req *overlays_now_playing.CreateRequest,
) (*overlays_now_playing.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	overlay := model.ChannelOverlayNowPlaying{
		ID:              uuid.New(),
		Preset:          overlays.ChannelOverlayNowPlayingPresetTransparent,
		FontFamily:      req.GetFontFamily(),
		FontWeight:      req.GetFontWeight(),
		BackgroundColor: req.GetBackgroundColor(),
		ChannelID:       dashboardId,
		ShowImage:       req.GetShowImage(),
		HideTimeout:     null.IntFrom(int64(req.GetHideTimeout())),
	}

	if err := c.Db.
		WithContext(ctx).
		Create(&overlay).Error; err != nil {
		return nil, err
	}

	return convertEntityToProto(overlay), nil
}
