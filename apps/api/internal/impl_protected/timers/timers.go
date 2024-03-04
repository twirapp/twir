package timers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/timers"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Timers struct {
	*impl_deps.Deps
}

func (c *Timers) convertEntity(entity *model.ChannelsTimers) *timers.Timer {
	return &timers.Timer{
		Id:                       entity.ID,
		ChannelId:                entity.ChannelID,
		Name:                     entity.Name,
		Enabled:                  entity.Enabled,
		TimeInterval:             entity.TimeInterval,
		MessageInterval:          entity.MessageInterval,
		LastTriggerMessageNumber: entity.LastTriggerMessageNumber,
		Responses: lo.Map(
			entity.Responses,
			func(r *model.ChannelsTimersResponses, _ int) *timers.Timer_Response {
				return &timers.Timer_Response{
					Id:         r.ID,
					Text:       r.Text,
					IsAnnounce: r.IsAnnounce,
					TimerId:    r.TimerID,
				}
			},
		),
	}
}

func (c *Timers) TimersGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*timers.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []*model.ChannelsTimers
	if err := c.Db.
		WithContext(ctx).
		Preload("Responses").
		Where(`"channelId" = ?`, dashboardId).
		Group(`"id"`).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	return &timers.GetResponse{
		Timers: lo.Map(
			entities, func(t *model.ChannelsTimers, _ int) *timers.Timer {
				return c.convertEntity(t)
			},
		),
	}, nil
}

func (c *Timers) TimersUpdate(
	ctx context.Context,
	request *timers.UpdateRequest,
) (*timers.Timer, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelsTimers{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "id" = ?`, dashboardId, request.Id).
		Preload("Responses").
		Find(entity).Error; err != nil {
		return nil, err
	}

	entity.Name = request.Timer.Name
	entity.Enabled = request.Timer.Enabled
	entity.LastTriggerMessageNumber = 0
	entity.MessageInterval = request.Timer.MessageInterval
	entity.TimeInterval = request.Timer.TimeInterval

	txErr := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.
				Where(`"timerId" = ?`, entity.ID).
				Delete(&model.ChannelsTimersResponses{}).
				Error; err != nil {
				return err
			}

			entity.Responses = lo.Map(
				request.Timer.Responses,
				func(r *timers.CreateData_Response, _ int) *model.ChannelsTimersResponses {
					return &model.ChannelsTimersResponses{
						ID:         uuid.New().String(),
						Text:       r.Text,
						IsAnnounce: r.IsAnnounce,
						TimerID:    entity.ID,
					}
				},
			)

			err := tx.Save(entity).Error
			return err
		},
	)
	if txErr != nil {
		return nil, txErr
	}

	timersRequest := timersbusservice.AddOrRemoveTimerRequest{
		TimerID: entity.ID,
	}
	err := c.Bus.Timers.RemoveTimer.Publish(timersRequest)
	if err != nil {
		c.Logger.Error("cannot remove timer from queue", slog.Any("err", err))
		return nil, fmt.Errorf("cannot remove timer from queue: %w", err)
	}

	if entity.Enabled {
		err := c.Bus.Timers.AddTimer.Publish(timersRequest)
		if err != nil {
			c.Logger.Error("cannot add timer to queue", slog.Any("err", err))
			return nil, fmt.Errorf("cannot add timer to queue: %w", err)
		}
	}

	return c.convertEntity(entity), nil
}

func (c *Timers) TimersDelete(
	ctx context.Context,
	request *timers.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	if err := c.Bus.Timers.RemoveTimer.Publish(
		timersbusservice.AddOrRemoveTimerRequest{TimerID: request.Id},
	); err != nil {
		c.Logger.Error("cannot remove timer from queue", slog.Any("err", err))
	}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.Id).
		Delete(&model.ChannelsTimers{}).
		Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Timers) TimersCreate(
	ctx context.Context,
	request *timers.CreateRequest,
) (*timers.Timer, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var currentCount int64
	if err := c.Db.Model(&model.ChannelsTimers{}).Where(
		`"channelId" = ?`,
		dashboardId,
	).Count(&currentCount).Error; err != nil {
		return nil, fmt.Errorf("cannot get timers count time: %w", err)
	}

	if currentCount >= 10 {
		return nil, fmt.Errorf("you cannot create more than 10 timers")
	}

	entity := &model.ChannelsTimers{
		ID:              uuid.New().String(),
		ChannelID:       dashboardId,
		Name:            request.Data.Name,
		Enabled:         true,
		TimeInterval:    request.Data.TimeInterval,
		MessageInterval: request.Data.MessageInterval,
		Responses: lo.Map(
			request.Data.Responses,
			func(r *timers.CreateData_Response, _ int) *model.ChannelsTimersResponses {
				return &model.ChannelsTimersResponses{
					ID:         uuid.New().String(),
					Text:       r.Text,
					IsAnnounce: r.IsAnnounce,
				}
			},
		),
	}

	if err := c.Db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	timersReq := timersbusservice.AddOrRemoveTimerRequest{TimerID: entity.ID}
	if entity.Enabled {
		c.Bus.Timers.AddTimer.Publish(timersReq)
	} else {
		c.Bus.Timers.RemoveTimer.Publish(timersReq)
	}

	return c.convertEntity(entity), nil
}

func (c *Timers) TimersEnableOrDisable(
	ctx context.Context,
	request *timers.PatchRequest,
) (*timers.Timer, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelsTimers{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "id" = ?`, dashboardId, request.Id).
		Preload("Responses").
		Find(entity).Error; err != nil {
		return nil, err
	}

	entity.Enabled = request.Enabled
	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	timersRequest := timersbusservice.AddOrRemoveTimerRequest{TimerID: entity.ID}
	if entity.Enabled {
		c.Bus.Timers.AddTimer.Publish(timersRequest)
	} else {
		c.Bus.Timers.RemoveTimer.Publish(timersRequest)
	}

	return c.convertEntity(entity), nil
}
