package timers

import (
	"context"
	"errors"
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) []model.ChannelsTimers {
	timers := []model.ChannelsTimers{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Preload("Responses").Find(&timers).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil
	}
	return timers
}

func handlePost(
	channelId string,
	dto *timerDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timer := model.ChannelsTimers{
		ID:                       uuid.NewV4().String(),
		ChannelID:                channelId,
		Name:                     dto.Name,
		Enabled:                  *dto.Enabled,
		TimeInterval:             int32(dto.TimeInterval),
		MessageInterval:          int32(dto.MessageInterval),
		LastTriggerMessageNumber: 0,
	}

	err := services.DB.Save(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create timer")
	}

	timerResponses := []model.ChannelsTimersResponses{}

	for _, t := range dto.Responses {
		response := model.ChannelsTimersResponses{
			ID:         uuid.NewV4().String(),
			Text:       t.Text,
			IsAnnounce: *t.IsAnnounce,
			TimerID:    timer.ID,
		}

		err := services.DB.Save(&response).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			services.DB.Where(`"id" = ?`, timer.ID).Delete(&model.ChannelsTimers{})

			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"cannot create timer responses",
			)
		}

		timerResponses = append(timerResponses, response)
	}

	services.TimersGrpc.AddTimerToQueue(context.Background(), &timers.Request{
		TimerId: timer.ID,
	})

	timer.Responses = &timerResponses
	return &timer, nil
}

func handleDelete(timerId string, services types.Services) error {
	timer := model.ChannelsTimers{}
	err := services.DB.Where("id = ?", timerId).First(&timer).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(http.StatusNotFound, "timer not found")
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "timer not found")
	}

	err = services.DB.Delete(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete timer")
	}

	if timer.Enabled {
		services.TimersGrpc.RemoveTimerFromQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	}

	return nil
}

func handlePut(
	channelId string,
	timerId string,
	dto *timerDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timer := model.ChannelsTimers{}

	err := services.DB.Where(`"channelId" = ? AND "id" = ?`, channelId, timerId).
		Preload("Responses").
		First(&timer).
		Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(401, "timer with that id not found")
	}

	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	timer.Enabled = *dto.Enabled
	timer.LastTriggerMessageNumber = 0
	timer.MessageInterval = int32(dto.MessageInterval)
	timer.Name = dto.Name
	timer.TimeInterval = int32(dto.TimeInterval)

	err = services.DB.Select("*").Updates(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update timer")
	}

	err = services.DB.Transaction(func(tx *gorm.DB) error {
		for _, response := range *timer.Responses {
			err = tx.Where("id = ?", response.ID).Delete(&model.ChannelsTimersResponses{}).Error
			if err != nil {
				services.Logger.Sugar().Error(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"internal error when deleting responses",
		)
	}

	newResponses := []model.ChannelsTimersResponses{}
	for _, response := range dto.Responses {
		r := model.ChannelsTimersResponses{
			ID:         uuid.NewV4().String(),
			Text:       response.Text,
			IsAnnounce: *response.IsAnnounce,
			TimerID:    timer.ID,
		}
		services.DB.Save(&r)
		newResponses = append(newResponses, r)
	}

	timer.Responses = &newResponses

	if timer.Enabled {
		services.TimersGrpc.AddTimerToQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	} else {
		services.TimersGrpc.RemoveTimerFromQueue(context.Background(), &timers.Request{
			TimerId: timer.ID,
		})
	}

	return &timer, nil
}

func handlePatch(
	channelId, timerId string,
	dto *timerPatchDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timer := model.ChannelsTimers{}
	err := services.DB.Where("id = ?", timerId).First(&timer).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(http.StatusNotFound, "timer not found")
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "timer not found")
	}

	timer.Enabled = *dto.Enabled
	err = services.DB.Select("*").Updates(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update timer")
	}

	if err = services.DB.Select("*").Preload("Responses").Find(&timer).Error; err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &timer, nil
}
