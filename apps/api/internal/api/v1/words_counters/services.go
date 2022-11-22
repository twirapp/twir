package words_counters

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelWordCounter, error) {
	counters := []model.ChannelWordCounter{}

	if err := services.DB.Where(`"channelId" = ?`, channelId).Find(&counters).Error; err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "unexpected internal error")
	}

	return counters, nil
}

func handlePost(channelId string, dto *wordsCountersDto, services types.Services) (*model.ChannelWordCounter, error) {
	phrase := strings.ToLower(dto.Phrase)
	existedCounter := model.ChannelWordCounter{}
	err := services.DB.
		Where(`"channelId" = ? AND "phrase" = ?`, channelId, phrase).
		Find(&existedCounter).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedCounter.ID != "" {
		return nil, fiber.NewError(http.StatusConflict, "counter with that phrase already exists")
	}

	newCounter := &model.ChannelWordCounter{
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		Phrase:    phrase,
		Counter:   0,
		Enabled:   *dto.Enabled,
	}
	err = services.DB.Save(newCounter).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return newCounter, nil
}

func handlePut(channelId, counterId string, dto *wordsCountersDto, services types.Services) (*model.ChannelWordCounter, error) {
	existedCounter := &model.ChannelWordCounter{}
	err := services.DB.
		Where(`"channelId" = ? AND "id" = ?`, channelId, counterId).
		Find(existedCounter).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedCounter.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "counter with that id not found")
	}

	existedCounter.Phrase = strings.ToLower(dto.Phrase)
	existedCounter.Counter = dto.Counter
	existedCounter.Enabled = *dto.Enabled

	err = services.DB.Save(existedCounter).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return existedCounter, nil
}

func handleDelete(channelId, counterId string, services types.Services) error {
	existedCounter := model.ChannelWordCounter{}
	err := services.DB.
		Where(`"channelId" = ? AND "id" = ?`, channelId, counterId).
		Find(&existedCounter).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedCounter.ID == "" {
		return fiber.NewError(http.StatusNotFound, "counter with that id not found")
	}

	err = services.DB.Delete(&existedCounter).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete counter due internal error")
	}

	return nil
}
