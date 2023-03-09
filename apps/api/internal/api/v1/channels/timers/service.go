package timers

import (
	"context"

	"github.com/samber/lo"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"
)

func (c *Timers) getService(channelId string) []model.ChannelsTimers {
	timers, err := c.services.TimersService.FindManyByChannelId(channelId)
	if err != nil {
		return nil
	}

	return timers
}

func (c *Timers) postService(
	channelId string,
	dto *timerDto,
) (*model.ChannelsTimers, error) {
	responses := lo.Map(dto.Responses, func(r responseDto, _ int) model.ChannelsTimersResponses {
		return model.ChannelsTimersResponses{
			Text:       r.Text,
			IsAnnounce: *r.IsAnnounce,
		}
	})

	timer, err := c.services.TimersService.Create(model.ChannelsTimers{
		ChannelID:                channelId,
		Name:                     dto.Name,
		Enabled:                  *dto.Enabled,
		TimeInterval:             int32(dto.TimeInterval),
		MessageInterval:          int32(dto.MessageInterval),
		LastTriggerMessageNumber: 0,
	}, responses)
	if err != nil {
		return nil, err
	}

	c.services.Grpc.Timers.AddTimerToQueue(context.Background(), &timers.Request{
		TimerId: timer.ID,
	})

	return timer, nil
}

func (c *Timers) deleteService(timerId string) error {
	err := c.services.TimersService.Delete(timerId)
	if err != nil {
		return err
	}

	c.services.Grpc.Timers.RemoveTimerFromQueue(context.Background(), &timers.Request{
		TimerId: timerId,
	})

	return nil
}

func (c *Timers) putService(
	timerId string,
	dto *timerDto,
) (*model.ChannelsTimers, error) {
	responses := lo.Map(dto.Responses, func(r responseDto, _ int) model.ChannelsTimersResponses {
		return model.ChannelsTimersResponses{
			Text:       r.Text,
			IsAnnounce: *r.IsAnnounce,
		}
	})

	timer, err := c.services.TimersService.Update(
		timerId,
		model.ChannelsTimers{
			Name:            dto.Name,
			MessageInterval: int32(dto.MessageInterval),
			TimeInterval:    int32(dto.TimeInterval),
			Enabled:         *dto.Enabled,
		},
		responses,
	)
	if err != nil {
		return nil, err
	}

	if timer.Enabled {
		c.services.Grpc.Timers.AddTimerToQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	} else {
		c.services.Grpc.Timers.RemoveTimerFromQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	}

	return timer, nil
}

func (c *Timers) patchService(
	timerId string,
	dto *timerPatchDto,
) (*model.ChannelsTimers, error) {
	updatedTimer, err := c.services.TimersService.SetEnabled(timerId, *dto.Enabled)
	if err != nil {
		return nil, err
	}

	if updatedTimer.Enabled {
		c.services.Grpc.Timers.AddTimerToQueue(context.Background(), &timers.Request{
			TimerId: timerId,
		})
	} else {
		c.services.Grpc.Timers.RemoveTimerFromQueue(context.Background(), &timers.Request{
			TimerId: timerId,
		})
	}

	return updatedTimer, nil
}
