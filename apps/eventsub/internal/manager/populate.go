package manager

import (
	"context"
	"log/slog"
	"sync"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *Manager) populateChannels() error {
	requestContext := context.Background()
	var channels []model.Channels
	err := c.gorm.Where(
		`"channels"."isEnabled" = ? AND "User"."is_banned" = ? AND "channels"."isTwitchBanned" = ?`,
		true,
		false,
		false,
	).Joins("User").Find(&channels).Error
	if err != nil {
		return err
	}

	var topics []model.EventsubTopic
	if err := c.gorm.WithContext(requestContext).Find(&topics).Error; err != nil {
		return err
	}

	channelsWg := sync.WaitGroup{}

	for _, channel := range channels {
		channelsWg.Add(1)

		channel := channel

		go func() {
			defer channelsWg.Done()
			err := c.SubscribeToNeededEvents(
				requestContext,
				topics,
				channel.ID,
				channel.BotID,
			)
			if err != nil {
				c.logger.Error(
					"failed to subscribe to needed events",
					slog.Any("err", err),
				)
			}
		}()
	}

	channelsWg.Wait()

	return nil
}
