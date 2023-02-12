package grpc_impl

import (
	"errors"
	"github.com/satont/tsuwari/apps/events/internal"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *EventsGrpcImplementation) processEvent(channelId string, data internal.Data, eventType string) error {
	dbEntities := []model.Event{}

	err := c.services.DB.
		Where(`"channelId" = ? AND "type" = ?`, channelId, eventType).
		Preload("Operations").
		Find(&dbEntities).
		Error

	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return err
	}

	for _, entity := range dbEntities {
		if entity.ID == "" {
			return errors.New("event not found")
		}

		if entity.Type == "COMMAND_USED" &&
			data.CommandID != "" &&
			entity.CommandID.Valid &&
			data.CommandID != entity.CommandID.String {
			continue
		}

		if entity.Type == "REDEMPTION_CREATED" &&
			data.RewardID != "" &&
			entity.RewardID.Valid &&
			data.RewardID != entity.RewardID.String {
			continue
		}

		go c.processOperations(channelId, entity.Operations, data)
	}

	return nil
}
