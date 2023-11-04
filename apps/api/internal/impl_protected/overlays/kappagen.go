package overlays

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/overlays_kappagen"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

const kappagenOverlayType = "kappagen_overlay"

func (c *Overlays) kappagenDbToGrpc(s model.ChatOverlaySettings) *overlays_kappagen.Settings {
	return &overlays_kappagen.Settings{}
}

func (c *Overlays) kappagenGrpcToDb(s *overlays_kappagen.Settings) model.ChatOverlaySettings {
	return model.ChatOverlaySettings{}
}

func (c *Overlays) OverlayKappaGenGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_kappagen.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, kappagenOverlayType).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	parsedSettings := model.ChatOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	return c.kappagenDbToGrpc(parsedSettings), nil
}

func (c *Overlays) OverlayKappaGenUpdate(
	ctx context.Context,
	req *overlays_kappagen.Settings,
) (*overlays_kappagen.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, kappagenOverlayType).
		Find(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	if entity.ID == "" {
		entity.ID = uuid.NewString()
		entity.ChannelId = dashboardId
		entity.Type = kappagenOverlayType
	}

	parsedSettings := c.kappagenGrpcToDb(req)
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

	c.Grpc.Websockets.RefreshKappagenOverlaySettings(
		ctx, &websockets.RefreshKappagenOverlaySettingsRequest{
			ChannelId: dashboardId,
		},
	)

	return c.kappagenDbToGrpc(newSettings), nil
}
