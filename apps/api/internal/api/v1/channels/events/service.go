package events

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (c *Events) getService(channelId string) []model.Event {
	events := []model.Event{}
	err := c.services.Gorm.
		Where(`"channelId" = ?`, channelId).
		Preload("Operations").
		Find(&events).Error
	if err != nil {
		c.services.Logger.Error(err)
	}

	return events
}

func (c *Events) postService(channelId string, dto *eventDto) (*model.Event, error) {
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

	err := c.services.Gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newEvent).Error; err != nil {
			return err
		}

		for i, operation := range dto.Operations {
			newOperation := &model.EventOperation{
				ID:          uuid.NewV4().String(),
				Type:        operation.Type,
				Delay:       operation.Delay,
				EventID:     newEvent.ID,
				Input:       null.StringFrom(strings.TrimSpace(*operation.Input)),
				Repeat:      operation.Repeat,
				Order:       i,
				UseAnnounce: *operation.UseAnnounce,
				TimeoutTime: operation.TimeoutTime,
				Target:      null.StringFrom(operation.Target),
			}

			if err := tx.Create(newOperation).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	c.services.Gorm.Where(`"id" = ?`, newEvent.ID).Preload("Operations").Find(newEvent)
	return newEvent, nil
}

func (c *Events) putService(channelId, eventId string, dto *eventDto) (*model.Event, error) {
	event := model.Event{}
	err := c.services.Gorm.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Find(&event).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "event not found")
	}

	event.RewardID = null.StringFromPtr(dto.RewardID)
	event.CommandID = null.StringFromPtr(dto.CommandID)
	event.KeywordID = null.StringFromPtr(dto.KeywordID)
	event.Description = null.StringFrom(dto.Description)

	err = c.services.Gorm.Transaction(func(tx *gorm.DB) error {
		if err = c.services.Gorm.Save(&event).Error; err != nil {
			return err
		}

		if err = c.services.Gorm.Where(`"eventId" = ?`, event.ID).Delete(&model.EventOperation{}).Error; err != nil {
			return err
		}

		for i, operation := range dto.Operations {
			newOperation := model.EventOperation{
				ID:          uuid.NewV4().String(),
				Type:        operation.Type,
				Delay:       operation.Delay,
				EventID:     event.ID,
				Input:       null.StringFrom(strings.TrimSpace(*operation.Input)),
				Repeat:      operation.Repeat,
				Order:       i,
				UseAnnounce: *operation.UseAnnounce,
				TimeoutTime: operation.TimeoutTime,
				Target:      null.StringFrom(operation.Target),
			}

			if err := tx.Save(&newOperation).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	c.services.Gorm.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Preload("Operations").Find(&event)
	return &event, nil
}

func (c *Events) deleteService(channelId, eventId string) error {
	event := model.Event{}
	err := c.services.Gorm.Where(`"id" = ? and "channelId" = ?`, eventId, channelId).Find(&event).Error
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return fiber.NewError(http.StatusNotFound, "event not found")
	}

	err = c.services.Gorm.Transaction(func(tx *gorm.DB) error {
		err = tx.Where(`"eventId" = ?`, event.ID).Delete(&model.EventOperation{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&event).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func (c *Events) patchService(
	channelId,
	eventId string,
	dto *eventPatchDto,
) (*model.Event, error) {
	event := model.Event{}
	err := c.services.Gorm.
		Where(`"id" = ? and "channelId" = ?`, eventId, channelId).
		Preload("Operations").
		Find(&event).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if event.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "event not found")
	}

	event.Enabled = *dto.Enabled
	err = c.services.Gorm.Model(&event).Select("*").Updates(&event).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update event")
	}

	return &event, nil
}
