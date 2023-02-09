package grpc_impl

import (
	"errors"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *EventsGrpcImplementation) processEvent(channelId string, data Data, eventType string) error {
	dbEntity := &model.Event{}

	err := c.services.DB.
		Where(`"channelId" = ? AND "type" = ?`, channelId, eventType).
		Preload("Operations").
		Find(dbEntity).
		Error

	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return err
	}

	if dbEntity == nil || dbEntity.ID == "" {
		return errors.New("event not found")
	}

	c.processOperations(channelId, dbEntity.Operations, data)

	return nil
}
