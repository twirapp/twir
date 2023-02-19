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
	"strings"
)

func handleGet(channelId string, services types.Services) []model.Event {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	events := []model.Event{}
	err := db.
		Where(`"channelId" = ?`, channelId).
		Preload("Operations").
		Find(&events).Error
	if err != nil {
		logger.Error(err)
	}

	return events
}

func handlePost(channelId string, dto *eventDto) (*model.Event, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	newEvent := &model.Event{
		ID:          uuid.NewV4().String(),
		ChannelID:   channelId,
		Type:        dto.Type,
		RewardID:    null.NewString(*dto.RewardID, *dto.RewardID != ""),
		CommandID:   null.NewString(*dto.CommandID, *dto.CommandID != ""),
		KeywordID:   null.NewString(*dto.KeywordID, *dto.KeywordID != ""),
		Description: null.StringFrom(dto.Description),
		Enabled:     true,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newEvent).Error; err != nil {
			return err
		}

		for i, operation := range dto.Operations {
			newOperation := &model.EventOperation{
				ID:            uuid.NewV4().String(),
				Type:          operation.Type,
				Delay:         operation.Delay,
				EventID:       newEvent.ID,
				Input:         null.StringFrom(strings.TrimSpace(*operation.Input)),
				Repeat:        operation.Repeat,
				Order:         i,
				UseAnnounce:   *operation.UseAnnounce,
				TimeoutTime:   operation.TimeoutTime,
				ObsTargetName: null.StringFrom(operation.ObsTargetName),
			}

			if err := tx.Create(newOperation).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	db.Where(`"id" = ?`, newEvent.ID).Preload("Operations").Find(newEvent)
	return newEvent, nil
}

func handleUpdate(channelId, eventId string, dto *eventDto) (*model.Event, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	event := model.Event{}
	err := db.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Find(&event).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "event not found")
	}

	event.RewardID = null.StringFromPtr(dto.RewardID)
	event.CommandID = null.StringFromPtr(dto.CommandID)
	event.CommandID = null.StringFromPtr(dto.KeywordID)
	event.Description = null.StringFrom(dto.Description)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err = db.Save(&event).Error; err != nil {
			return err
		}

		if err = db.Where(`"eventId" = ?`, event.ID).Delete(&model.EventOperation{}).Error; err != nil {
			return err
		}

		for i, operation := range dto.Operations {
			newOperation := model.EventOperation{
				ID:            uuid.NewV4().String(),
				Type:          operation.Type,
				Delay:         operation.Delay,
				EventID:       event.ID,
				Input:         null.StringFrom(strings.TrimSpace(*operation.Input)),
				Repeat:        operation.Repeat,
				Order:         i,
				UseAnnounce:   *operation.UseAnnounce,
				TimeoutTime:   operation.TimeoutTime,
				ObsTargetName: null.StringFrom(operation.ObsTargetName),
			}

			if err := tx.Save(&newOperation).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	db.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Preload("Operations").Find(&event)
	return &event, nil
}

func handleDelete(channelId, eventId string) error {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	event := model.Event{}
	err := db.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Find(&event).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return fiber.NewError(http.StatusNotFound, "event not found")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Where(`"eventId" = ?`, event.ID).Delete(&model.EventOperation{}).Error
		if err != nil {
			return err
		}

		err = db.Delete(&event).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func handlePatch(
	channelId, eventId string,
	dto *eventPatchDto,
	services types.Services,
) (*model.Event, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	event := model.Event{}
	err := db.
		Where(`"id" = ? and "channelId" = ?`, eventId, channelId).
		Preload("Operations").
		Find(&event).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "event not found")
	}

	event.Enabled = *dto.Enabled
	err = services.DB.Model(&event).Select("*").Updates(&event).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update event")
	}

	return &event, nil
}
