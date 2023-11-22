package overlays

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/overlays_be_right_back"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

const brbOverlayType = "be_right_back_overlay"

func (c *Overlays) brbDbToGrpc(
	s model.
		ChannelModulesSettingsBeRightBack,
) *overlays_be_right_back.Settings {
	return &overlays_be_right_back.Settings{
		Text: s.Text,
		Late: &overlays_be_right_back.Settings_Late{
			Enabled:         s.Late.Enabled,
			Text:            s.Late.Text,
			DisplayBrbTime:  s.Late.DisplayBrbTime,
			DisplayLateTime: s.Late.DisplayLateTime,
		},
		BackgroundColor: s.BackgroundColor,
		FontSize:        s.FontSize,
		FontColor:       s.FontColor,
	}
}

func (c *Overlays) brbGrpcToDb(s *overlays_be_right_back.Settings) model.
	ChannelModulesSettingsBeRightBack {
	return model.ChannelModulesSettingsBeRightBack{
		Text: s.Text,
		Late: model.ChannelModulesSettingsBeRightBackLate{
			Enabled:         s.Late.Enabled,
			Text:            s.Late.Text,
			DisplayBrbTime:  s.Late.DisplayBrbTime,
			DisplayLateTime: s.Late.DisplayLateTime,
		},
		BackgroundColor: s.BackgroundColor,
		FontSize:        s.FontSize,
		FontColor:       s.FontColor,
	}
}

func (c *Overlays) OverlayBeRightBackGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_be_right_back.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, brbOverlayType).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	parsedSettings := model.ChannelModulesSettingsBeRightBack{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	return c.brbDbToGrpc(parsedSettings), nil
}

func (c *Overlays) OverlayBeRightBackUpdate(
	ctx context.Context,
	req *overlays_be_right_back.Settings,
) (*overlays_be_right_back.Settings, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, brbOverlayType).
		Find(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	if entity.ID == "" {
		entity.ID = uuid.NewString()
		entity.ChannelId = dashboardId
		entity.Type = brbOverlayType
	}

	parsedSettings := c.brbGrpcToDb(req)
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

	newSettings := model.ChannelModulesSettingsBeRightBack{}
	if err := json.Unmarshal(entity.Settings, &newSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	c.Grpc.Websockets.RefreshBrbSettings(
		ctx, &websockets.RefreshBrbSettingsRequest{
			ChannelId: dashboardId,
		},
	)

	return c.brbDbToGrpc(newSettings), nil
}
