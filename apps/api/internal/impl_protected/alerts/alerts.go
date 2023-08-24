package alerts

import (
	"context"
	"log/slog"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/alerts"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Alerts struct {
	*impl_deps.Deps
}

func (c *Alerts) convertEntity(entity model.ChannelAlert) *alerts.Alert {
	return &alerts.Alert{
		Id:          entity.ID,
		Name:        entity.Name,
		AudioId:     entity.AudioID.Ptr(),
		AudioVolume: int32(entity.AudioVolume),
		CommandIds:  entity.CommandIDS,
		RewardIds:   entity.RewardIDS,
	}
}

func (c *Alerts) AlertsGetAll(ctx context.Context, _ *emptypb.Empty) (
	*alerts.GetAllResponse,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []model.ChannelAlert

	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ?",
		dashboardId,
	).Order("name desc").Find(&entities).Error; err != nil {
		c.Logger.Error(
			"cannot get alerts",
			slog.String("channelId", dashboardId),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &alerts.GetAllResponse{
		Alerts: lo.Map(
			entities, func(item model.ChannelAlert, _ int) *alerts.Alert {
				return c.convertEntity(item)
			},
		),
	}, nil
}

func (c *Alerts) AlertsCreate(ctx context.Context, req *alerts.CreateRequest) (
	*alerts.Alert,
	error,
) {
	if utf8.RuneCountInString(req.Name) > 30 {
		return nil, twirp.NewError(twirp.OutOfRange, "Alert name is too long")
	}

	dashboardId := ctx.Value("dashboardId").(string)
	entity := model.ChannelAlert{
		ID:          uuid.New().String(),
		ChannelID:   dashboardId,
		Name:        req.Name,
		AudioID:     null.StringFromPtr(req.AudioId),
		AudioVolume: int(req.AudioVolume),
		CommandIDS:  req.CommandIds,
		RewardIDS:   req.RewardIds,
	}

	if err := c.Db.WithContext(ctx).Create(&entity).Error; err != nil {
		c.Logger.Error(
			"cannot create alert",
			slog.String("channelId", dashboardId),
			slog.Any("err", err),
		)
		return nil, err
	}
	return c.convertEntity(entity), nil
}

func (c *Alerts) AlertsDelete(ctx context.Context, req *alerts.RemoveRequest) (
	*emptypb.Empty,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	if err := c.Db.WithContext(ctx).Where(
		"id = ? AND channel_id = ?",
		req.Id,
		dashboardId,
	).Delete(&model.ChannelAlert{}).
		Error; err != nil {
		c.Logger.Error(
			"cannot create alert",
			slog.String("channelId", dashboardId),
			slog.String("id", req.Id),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Alerts) AlertsUpdate(ctx context.Context, req *alerts.UpdateRequest) (
	*alerts.Alert,
	error,
) {
	if len(req.Name) > 30 {
		return nil, twirp.NewError(twirp.OutOfRange, "Alert name is too long")
	}

	dashboardId := ctx.Value("dashboardId").(string)
	entity := model.ChannelAlert{}
	if err := c.Db.WithContext(ctx).Where(
		"id = ? AND channel_id = ?",
		req.Id,
		dashboardId,
	).Find(&entity).Error; err != nil {
		c.Logger.Error(
			"cannot find alert",
			slog.String("channelId", dashboardId),
			slog.String("id", req.Id),
			slog.Any("err", err),
		)
		return nil, err
	}

	entity.Name = req.Name
	entity.AudioID = null.StringFromPtr(req.AudioId)
	entity.AudioVolume = int(req.AudioVolume)
	entity.CommandIDS = req.CommandIds
	entity.RewardIDS = req.RewardIds

	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		c.Logger.Error(
			"cannot update alert",
			slog.String("channelId", dashboardId),
			slog.String("id", req.Id),
			slog.Any("err", err),
		)
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Alerts) AlertsGetOne(ctx context.Context, req *alerts.GetOneRequest) (
	*alerts.Alert,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := model.ChannelAlert{}
	if err := c.Db.WithContext(ctx).Where(
		"id = ? AND channel_id = ?",
		req.Id,
		dashboardId,
	).Find(&entity).
		Error; err != nil {
		c.Logger.Error(
			"cannot get alert",
			slog.String("channelId", dashboardId),
			slog.String("id", req.Id),
			slog.Any("err", err),
		)
		return nil, err
	}

	if entity.ID == "" {
		return nil, twirp.NewError(twirp.NotFound, "Alert not found")
	}

	return c.convertEntity(entity), nil
}
