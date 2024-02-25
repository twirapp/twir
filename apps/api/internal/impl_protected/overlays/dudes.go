package overlays

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/overlays_dudes"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func convertDudesEntityToGrpc(entity model.ChannelsOverlaysDudes) *overlays_dudes.Settings {
	return &overlays_dudes.Settings{
		Id: lo.ToPtr(entity.ID.String()),
		DudeSettings: &overlays_dudes.DudeSettings{
			Color:          entity.DudeColor,
			EyesColor:      entity.DudeEyesColor,
			CosmeticsColor: entity.DudeCosmeticsColor,
			MaxLifeTime:    entity.DudeMaxLifeTime,
			Gravity:        entity.DudeGravity,
			Scale:          entity.DudeScale,
			SoundsEnabled:  entity.DudeSoundsEnabled,
			SoundsVolume:   entity.DudeSoundsVolume,
			VisibleName:    entity.DudeVisibleName,
			GrowTime:       entity.DudeGrowTime,
			GrowMaxScale:   entity.DudeGrowMaxScale,
			MaxOnScreen:    entity.DudeMaxOnScreen,
			DefaultSprite:  entity.DudeDefaultSprite,
		},
		MessageBoxSettings: &overlays_dudes.MessageBoxSettings{
			Enabled:      entity.MessageBoxEnabled,
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
		IgnoreSettings: &overlays_dudes.IgnoreSettings{
			IgnoreCommands: entity.IgnoreCommands,
			IgnoreUsers:    entity.IgnoreUsers,
			Users:          entity.IgnoredUsers,
		},
		SpitterEmoteSettings: &overlays_dudes.SpitterEmoteSettings{
			Enabled: entity.SpitterEmoteEnabled,
		},
	}
}

func convertDudesGrpcToDb(settings *overlays_dudes.Settings) model.ChannelsOverlaysDudes {
	return model.ChannelsOverlaysDudes{
		DudeColor:                 settings.GetDudeSettings().GetColor(),
		DudeEyesColor:             settings.GetDudeSettings().GetEyesColor(),
		DudeCosmeticsColor:        settings.GetDudeSettings().GetCosmeticsColor(),
		DudeMaxLifeTime:           settings.GetDudeSettings().GetMaxLifeTime(),
		DudeGravity:               settings.GetDudeSettings().GetGravity(),
		DudeScale:                 settings.GetDudeSettings().GetScale(),
		DudeSoundsEnabled:         settings.GetDudeSettings().GetSoundsEnabled(),
		DudeSoundsVolume:          settings.GetDudeSettings().GetSoundsVolume(),
		DudeVisibleName:           settings.GetDudeSettings().GetVisibleName(),
		DudeGrowTime:              settings.GetDudeSettings().GetGrowTime(),
		DudeGrowMaxScale:          settings.GetDudeSettings().GetGrowMaxScale(),
		DudeDefaultSprite:         settings.GetDudeSettings().GetDefaultSprite(),
		DudeMaxOnScreen:           settings.GetDudeSettings().GetMaxOnScreen(),
		MessageBoxEnabled:         settings.GetMessageBoxSettings().GetEnabled(),
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
		IgnoreCommands:            settings.GetIgnoreSettings().GetIgnoreCommands(),
		IgnoreUsers:               settings.GetIgnoreSettings().GetIgnoreUsers(),
		IgnoredUsers:              append(pq.StringArray{}, settings.GetIgnoreSettings().GetUsers()...),
		SpitterEmoteEnabled:       settings.GetSpitterEmoteSettings().GetEnabled(),
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
	if err := c.Db.WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		Order("created_at asc").
		Find(&entities).
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
	entity.CreatedAt = time.Now()

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
	newEntity.CreatedAt = existedEntity.CreatedAt

	if err := c.Db.WithContext(ctx).Save(&newEntity).Error; err != nil {
		return nil, err
	}

	if _, err := c.Grpc.Websockets.RefreshOverlaySettings(
		ctx,
		&websockets.RefreshOverlaysRequest{
			ChannelId:   dashboardId,
			OverlayName: websockets.RefreshOverlaySettingsName_DUDES,
			OverlayId:   lo.ToPtr(existedEntity.ID.String()),
		},
	); err != nil {
		c.Logger.Error("cannot refresh chat overlay settings", slog.Any("err", err))
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
