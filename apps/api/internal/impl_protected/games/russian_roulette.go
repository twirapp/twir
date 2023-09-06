package games

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/games"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var rouletteType = "russian_roulette"

func (c *Games) GamesGetRouletteSettings(
	ctx context.Context,
	_ *emptypb.Empty,
) (*games.RussianRoulleteSettingsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, rouletteType).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	settings := model.RussianRouletteSetting{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, err
	}

	return &games.RussianRoulleteSettingsResponse{
		Enabled:              settings.Enabled,
		CanBeUsedByModerator: settings.CanBeUsedByModerators,
		TimeoutSeconds:       int32(settings.TimeoutSeconds),
		DecisionSeconds:      int32(settings.DecisionSeconds),
		ChargedBullets:       int32(settings.ChargedBullets),
		InitMessage:          settings.InitMessage,
		SurviveMessage:       settings.SurviveMessage,
		DeathMessage:         settings.DeathMessage,
	}, nil
}

var maxTimeoutTime = 24 * 7 * 2 * time.Hour

func (c *Games) GamesUpdateRouletteSettings(
	ctx context.Context,
	req *games.UpdateRussianRoulleteSettings,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	if req.TimeoutSeconds > int32(maxTimeoutTime.Seconds()) {
		return nil, twirp.NewError("400", "Max timeout is 14 days")
	}

	if req.DecisionSeconds > 60 {
		return nil, twirp.NewError("400", "Max decision time is 60 seconds")
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, rouletteType).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}

	if entity.ID == "" {
		entity.ID = uuid.New().String()
		entity.ChannelId = dashboardId
		entity.Type = rouletteType
	}

	settings := model.RussianRouletteSetting{
		Enabled:               req.Enabled,
		CanBeUsedByModerators: req.CanBeUsedByModerator,
		TimeoutSeconds:        int(req.TimeoutSeconds),
		DecisionSeconds:       int(req.DecisionSeconds),
		ChargedBullets:        int(req.ChargedBullets),
		InitMessage:           req.InitMessage,
		SurviveMessage:        req.SurviveMessage,
		DeathMessage:          req.DeathMessage,
	}

	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	entity.Settings = settingsJson

	if err := c.Db.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
