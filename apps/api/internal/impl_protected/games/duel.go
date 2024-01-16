package games

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/games"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var duelType = "duel"

func (c *Games) GamesGetDuelSettings(
	ctx context.Context,
	_ *emptypb.Empty,
) (*games.DuelSettingsResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, twirp.NewErrorf(twirp.NotFound, "cannot get selected dashboard: %w", err)
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, duelType).
		First(&entity).
		Error; err != nil {
		return nil, twirp.NewErrorf(twirp.NotFound, "cannot get duel settings: %w", err)
	}

	settings := model.ChannelModulesSettingsDuel{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, twirp.NewErrorf(twirp.Internal, "cannot parse duel settings: %w", err)
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
		BothDieMessage:  settings.BothDieMessage,
	}, nil
}

func (c *Games) GamesUpdateDuelSettings(
	ctx context.Context,
	req *games.UpdateDuelSettings,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
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

	if entity.ID == "" {
		entity.ID = uuid.NewString()
		entity.ChannelId = dashboardId
		entity.Type = duelType
	}

	settings := model.ChannelModulesSettingsDuel{}
	if entity.Settings != nil {
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return nil, err
		}
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
	settings.BothDieMessage = req.BothDieMessage

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
			if err := tx.
				WithContext(ctx).
				Model(&model.ChannelsCommands{}).
				Where(`"channelId" = ? and "defaultName" = ?`, dashboardId, "duel").
				Update("enabled", req.Enabled).
				Error; err != nil {
				return err
			}

			if err := tx.
				WithContext(ctx).
				Model(&model.ChannelsCommands{}).
				Where(`"channelId" = ? and "defaultName" = ?`, dashboardId, "duel accept").
				Update("enabled", req.Enabled).
				Error; err != nil {
				return err
			}

			if err := tx.
				WithContext(ctx).
				Save(&entity).
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
