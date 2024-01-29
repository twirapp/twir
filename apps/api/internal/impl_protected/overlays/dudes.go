package overlays

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/overlays_dudes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func convertDudesEntityToGrpc(entity model.ChannelsOverlaysDudes) *overlays_dudes.Settings {
	return &overlays_dudes.Settings{
		Id: lo.ToPtr(entity.ID.String()),
		DudeSettings: &overlays_dudes.DudeSettings{
			Color:         entity.DudeColor,
			MaxLifeTime:   entity.DudeMaxLifeTime,
			Gravity:       entity.DudeGravity,
			Scale:         entity.DudeScale,
			SoundsEnabled: entity.DudeSoundsEnabled,
			SoundsVolume:  entity.DudeSoundsVolume,
		},
		MessageBoxSettings: &overlays_dudes.MessageBoxSettings{
			BorderRadius: entity.MessageBoxBorderRadius,
			BoxColor:     entity.MessageBoxBoxColor,
			FontFamily:   entity.MessageBoxFontFamily,
			FontSize:     entity.MessageBoxFontSize,
			Padding:      entity.MessageBoxPadding,
			ShowTime:     entity.MessageBoxShowTime,
			Fill:         entity.MessageBoxFill,
		},
		NameBoxSettings: &overlays_dudes.NameBoxSettings{
			FontFamily:         entity.NameBoxFontFamily,
			FontSize:           entity.NameBoxFontSize,
			Fill:               entity.NameBoxFill,
			LineJoin:           entity.NameBoxLineJoin,
			StrokeThickness:    entity.NameBoxStrokeThickness,
			Stroke:             entity.NameBoxStroke,
			FillGradientStops:  entity.NameBoxFillGradientStops,
			FillGradientType:   entity.NameBoxFillGradientType,
			FontStyle:          entity.NameBoxFontStyle,
			FontVariant:        entity.NameBoxFontVariant,
			FontWeight:         entity.NameBoxFontWeight,
			DropShadow:         entity.NameBoxDropShadow,
			DropShadowAlpha:    entity.NameBoxDropShadowAlpha,
			DropShadowAngle:    entity.NameBoxDropShadowAngle,
			DropShadowBlur:     entity.NameBoxDropShadowBlur,
			DropShadowDistance: entity.NameBoxDropShadowDistance,
			DropShadowColor:    entity.NameBoxDropShadowColor,
		},
	}
}

func convertDudesGrpcToDb(settings *overlays_dudes.Settings) model.ChannelsOverlaysDudes {
	return model.ChannelsOverlaysDudes{
		DudeColor:                 settings.GetDudeSettings().GetColor(),
		DudeMaxLifeTime:           settings.GetDudeSettings().GetMaxLifeTime(),
		DudeGravity:               settings.GetDudeSettings().GetGravity(),
		DudeScale:                 settings.GetDudeSettings().GetScale(),
		DudeSoundsEnabled:         settings.GetDudeSettings().GetSoundsEnabled(),
		DudeSoundsVolume:          settings.GetDudeSettings().GetSoundsVolume(),
		MessageBoxBorderRadius:    settings.GetMessageBoxSettings().GetBorderRadius(),
		MessageBoxBoxColor:        settings.GetMessageBoxSettings().GetBoxColor(),
		MessageBoxFontFamily:      settings.GetMessageBoxSettings().GetFontFamily(),
		MessageBoxFontSize:        settings.GetMessageBoxSettings().GetFontSize(),
		MessageBoxPadding:         settings.GetMessageBoxSettings().GetPadding(),
		MessageBoxShowTime:        settings.GetMessageBoxSettings().GetShowTime(),
		MessageBoxFill:            settings.GetMessageBoxSettings().GetFill(),
		NameBoxFontFamily:         settings.GetNameBoxSettings().GetFontFamily(),
		NameBoxFontSize:           settings.GetNameBoxSettings().GetFontSize(),
		NameBoxFill:               settings.GetNameBoxSettings().GetFill(),
		NameBoxLineJoin:           settings.GetNameBoxSettings().GetLineJoin(),
		NameBoxStrokeThickness:    settings.GetNameBoxSettings().GetStrokeThickness(),
		NameBoxStroke:             settings.GetNameBoxSettings().GetStroke(),
		NameBoxFillGradientStops:  settings.GetNameBoxSettings().GetFillGradientStops(),
		NameBoxFillGradientType:   settings.GetNameBoxSettings().GetFillGradientType(),
		NameBoxFontStyle:          settings.GetNameBoxSettings().GetFontStyle(),
		NameBoxFontVariant:        settings.GetNameBoxSettings().GetFontVariant(),
		NameBoxFontWeight:         settings.GetNameBoxSettings().GetFontWeight(),
		NameBoxDropShadow:         settings.GetNameBoxSettings().GetDropShadow(),
		NameBoxDropShadowAlpha:    settings.GetNameBoxSettings().GetDropShadowAlpha(),
		NameBoxDropShadowAngle:    settings.GetNameBoxSettings().GetDropShadowAngle(),
		NameBoxDropShadowBlur:     settings.GetNameBoxSettings().GetDropShadowBlur(),
		NameBoxDropShadowDistance: settings.GetNameBoxSettings().GetDropShadowDistance(),
		NameBoxDropShadowColor:    settings.GetNameBoxSettings().GetDropShadowColor(),
	}
}

func (c *Overlays) OverlayDudesGet(
	ctx context.Context,
	req *overlays_dudes.GetRequest,
) (*overlays_dudes.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var entity model.ChannelsOverlaysDudes
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.GetId(),
	).First(&entity).
		Error; err != nil {
		return nil, err
	}

	return convertDudesEntityToGrpc(entity), nil
}

func (c *Overlays) OverlayDudesGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_dudes.GetAllResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	entities := []model.ChannelsOverlaysDudes{}
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ?",
		dashboardId,
	).First(&entities).
		Error; err != nil {
		return nil, err
	}

	convertedEntities := make([]*overlays_dudes.Settings, 0, len(entities))
	for _, entity := range entities {
		convertedEntities = append(convertedEntities, convertDudesEntityToGrpc(entity))
	}

	return &overlays_dudes.GetAllResponse{Settings: convertedEntities}, nil
}

func (c *Overlays) OverlayDudesCreate(
	ctx context.Context,
	req *overlays_dudes.Settings,
) (*overlays_dudes.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	entity := convertDudesGrpcToDb(req)
	entity.ID = uuid.New()
	entity.ChannelID = dashboardId

	if err := c.Db.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	return convertDudesEntityToGrpc(entity), nil
}

func (c *Overlays) OverlayDudesUpdate(
	ctx context.Context,
	req *overlays_dudes.UpdateRequest,
) (*overlays_dudes.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var existedEntity model.ChannelsOverlaysDudes
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? and id = ?", dashboardId,
		req.GetId(),
	).First(&existedEntity).Error; err != nil {
		return nil, err
	}

	reqModel := convertDudesGrpcToDb(req.GetSettings())
	newEntity := reqModel
	newEntity.ID = existedEntity.ID
	newEntity.ChannelID = existedEntity.ChannelID

	if err := c.Db.WithContext(ctx).Save(&newEntity).Error; err != nil {
		return nil, err
	}

	return convertDudesEntityToGrpc(newEntity), nil
}

func (c *Overlays) OverlayDudesDelete(
	ctx context.Context,
	req *overlays_dudes.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? and id = ?", dashboardId,
		req.GetId(),
	).Delete(&model.ChannelsOverlaysDudes{}).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
