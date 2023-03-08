package timers

import (
	"context"

	"github.com/samber/lo"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"

	"github.com/satont/tsuwari/apps/api/internal/types"
)

func handleGet(channelId string, services *types.Services) []model.ChannelsTimers {
	timers, err := services.TimersService.FindManyByChannelId(channelId)
	if err != nil {
		return nil
	}

	return timers
}

func handlePost(
	channelId string,
	dto *timerDto,
	services *types.Services,
) (*model.ChannelsTimers, error) {
	responses := lo.Map(dto.Responses, func(r responseDto, _ int) model.ChannelsTimersResponses {
		return model.ChannelsTimersResponses{
			Text:       r.Text,
			IsAnnounce: *r.IsAnnounce,
		}
	})

	timer, err := services.TimersService.Create(model.ChannelsTimers{
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

	services.Grpc.Timers.AddTimerToQueue(context.Background(), &timers.Request{
		TimerId: timer.ID,
	})

	return timer, nil
}

func handleDelete(timerId string, services *types.Services) error {
	err := services.TimersService.Delete(timerId)
	if err != nil {
		return err
	}

	services.Grpc.Timers.RemoveTimerFromQueue(context.Background(), &timers.Request{
		TimerId: timerId,
	})

	return nil
}

func handlePut(
	timerId string,
	dto *timerDto,
	services *types.Services,
) (*model.ChannelsTimers, error) {
	responses := lo.Map(dto.Responses, func(r responseDto, _ int) model.ChannelsTimersResponses {
		return model.ChannelsTimersResponses{
			Text:       r.Text,
			IsAnnounce: *r.IsAnnounce,
		}
	})

	timer, err := services.TimersService.Update(
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
		services.Grpc.Timers.AddTimerToQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	} else {
		services.Grpc.Timers.RemoveTimerFromQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	}

	return timer, nil
}

func handlePatch(
	timerId string,
	dto *timerPatchDto,
	services *types.Services,
) (*model.ChannelsTimers, error) {
	updatedTimer, err := services.TimersService.SetEnabled(timerId, *dto.Enabled)
	if err != nil {
		return nil, err
	}

	if updatedTimer.Enabled {
		services.Grpc.Timers.AddTimerToQueue(context.Background(), &timers.Request{
			TimerId: timerId,
		})
	} else {
		services.Grpc.Timers.RemoveTimerFromQueue(context.Background(), &timers.Request{
			TimerId: timerId,
		})
	}

	return updatedTimer, nil
}
