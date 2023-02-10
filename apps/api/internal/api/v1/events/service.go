package events

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func handleGet(channelId string, services types.Services) []model.Event {
	db := do.MustInvoke[gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	events := []model.Event{}
	err := db.Where(`"channelId" = ?`, channelId).Find(&events).Error
	if err != nil {
		logger.Error(err)
	}

	return events
}

func handlePost(channelId string, dto *eventDto) error {
	db := do.MustInvoke[gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	newEvent := model.Event{
		ID:          uuid.NewV4().String(),
		ChannelID:   channelId,
		Type:        dto.Type,
		RewardID:    null.StringFromPtr(dto.RewardID),
		CommandID:   null.StringFromPtr(dto.CommandID),
		Description: null.StringFromPtr(dto.Description),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newEvent).Error; err != nil {
			return err
		}

		for i, operation := range dto.Operations {
			newOperation := model.EventOperation{
				ID:      uuid.NewV4().String(),
				Type:    operation.Type,
				Delay:   null.IntFromPtr(operation.Delay),
				EventID: newEvent.ID,
				Input:   null.StringFromPtr(operation.Input),
				Repeat:  null.IntFromPtr(operation.Repeat),
				Order:   i,
			}

			if err := tx.Create(&newOperation).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func handleUpdate(channelId, eventId string, dto *eventDto) error {
	db := do.MustInvoke[gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	event := model.Event{}
	err := db.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return fiber.NewError(http.StatusNotFound, "event not found")
	}

	event.RewardID = null.StringFromPtr(dto.RewardID)
	event.CommandID = null.StringFromPtr(dto.CommandID)
	event.Description = null.StringFromPtr(dto.Description)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err = db.Save(&event).Error; err != nil {
			return err
		}

		if err = db.Where(`"eventId" = ?`, event.ID).Delete(&model.EventOperation{}).Error; err != nil {
			return err
		}

		for i, operation := range dto.Operations {
			newOperation := model.EventOperation{
				ID:      uuid.NewV4().String(),
				Type:    operation.Type,
				Delay:   null.IntFromPtr(operation.Delay),
				EventID: event.ID,
				Input:   null.StringFromPtr(operation.Input),
				Repeat:  null.IntFromPtr(operation.Repeat),
				Order:   i,
			}

			if err := tx.Create(&newOperation).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
