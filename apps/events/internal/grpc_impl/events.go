package grpc_impl

import (
	"errors"
	"github.com/satont/twir/apps/events/internal"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *EventsGrpcImplementation) processEvent(channelId string, data internal.Data, eventType model.EventType) error {
	var dbEntities []model.Event

	err := c.services.DB.
		Where(`"channelId" = ? AND "type" = ? AND "enabled" = ?`, channelId, eventType, true).
		Preload("Channel").
		Preload("Operations").
		Preload("Operations.Filters").
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

		if entity.Channel != nil && (!entity.Channel.IsEnabled || entity.Channel.IsBanned) {
			continue
		}

		if entity.Type == model.EventTypeCommandUsed &&
			data.CommandID != "" &&
			entity.CommandID.Valid &&
			data.CommandID != entity.CommandID.String {
			continue
		}

		if entity.Type == model.EventTypeRedemptionCreated &&
			data.RewardID != "" &&
			entity.RewardID.Valid &&
			data.RewardID != entity.RewardID.String {
			continue
		}

		if entity.Type == model.EventTypeKeywordMatched &&
			data.RewardID != "" &&
			entity.RewardID.Valid &&
			data.KeywordID != entity.KeywordID.String {
			continue
		}

		go c.processOperations(channelId, entity, data)
	}

	return nil
}
