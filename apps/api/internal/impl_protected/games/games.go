package games

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/games"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Games struct {
	*impl_deps.Deps
}

func (c *Games) GamesGetEightBallSettings(
	ctx context.Context,
	_ *emptypb.Empty,
) (*games.EightBallSettingsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = '8ball'`, dashboardId).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	settings := model.EightBallSettings{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, err
	}

	return &games.EightBallSettingsResponse{
		Answers: settings.Answers,
		Enabled: settings.Enabled,
	}, nil
}

const maxAnswers = 25

func (c *Games) GamesUpdateEightBallSettings(
	ctx context.Context,
	req *games.UpdateEightBallSettings,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	if len(req.Answers) > maxAnswers {
		return nil, twirp.NewError("400", fmt.Sprintf("Max answers is %v", maxAnswers))
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = '8ball'`, dashboardId).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}

	if entity.ID == "" {
		entity.ID = uuid.New().String()
		entity.ChannelId = dashboardId
		entity.Type = "8ball"
	}

	settings := model.EightBallSettings{
		Answers: req.Answers,
		Enabled: req.Enabled,
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
