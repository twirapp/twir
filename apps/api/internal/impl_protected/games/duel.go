package games

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/games"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var duelType = "duel"

func (c *Games) GamesGetDuelSettings(
	ctx context.Context,
	_ *emptypb.Empty,
) (*games.DuelSettingsResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get selected dashboard: %w", err)
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, duelType).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	settings := model.ChannelModulesSettingsDuel{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, err
	}

	return &games.DuelSettingsResponse{
		UserCooldown:    settings.UserCooldown,
		GlobalCooldown:  settings.GlobalCooldown,
		TimeoutSeconds:  settings.TimeoutSeconds,
		StartMessage:    settings.StartMessage,
		ResultMessage:   settings.ResultMessage,
		Enabled:         settings.Enabled,
		SecondsToAccept: settings.SecondsToAccept,
		PointsPerWin:    settings.PointsPerWin,
		PointsPerLose:   settings.PointsPerLose,
		BothDiePercent:  settings.BothDiePercent,
	}, nil
}

func (c *Games) GamesUpdateDuelSettings(
	ctx context.Context,
	req *games.UpdateDuelSettings,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get selected dashboard: %w", err)
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, duelType).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}

	settings := model.ChannelModulesSettingsDuel{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, err
	}

	settings.UserCooldown = req.UserCooldown
	settings.GlobalCooldown = req.GlobalCooldown
	settings.TimeoutSeconds = req.TimeoutSeconds
	settings.StartMessage = req.StartMessage
	settings.ResultMessage = req.ResultMessage
	settings.Enabled = req.Enabled
	settings.SecondsToAccept = req.SecondsToAccept
	settings.PointsPerWin = req.PointsPerWin
	settings.PointsPerLose = req.PointsPerLose
	settings.BothDiePercent = req.BothDiePercent

	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		entity.ID = uuid.New().String()
		entity.ChannelId = dashboardId
		entity.Type = duelType
	}

	entity.Settings = settingsJson

	txErr := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			duelCommand := model.ChannelsCommands{}
			if err := tx.
				WithContext(ctx).
				Where(`"channelId" = ? and "defaultName" = ?`, dashboardId, "duel").
				First(&duelCommand).
				Error; err != nil {
				return err
			}

			duelCommand.Enabled = req.Enabled

			if err := tx.
				WithContext(ctx).
				Save(&entity).
				Error; err != nil {
				return err
			}

			if err := tx.
				WithContext(ctx).
				Save(&duelCommand).
				Error; err != nil {
				return err
			}

			return nil
		},
	)

	if txErr != nil {
		return nil, fmt.Errorf("cannot update duel settings: %w", txErr)
	}

	return &emptypb.Empty{}, nil
}
