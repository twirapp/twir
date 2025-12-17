package bus_listener

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

func (c *BusListener) banUsers(
	ctx context.Context,
	req []bots.BanRequest,
) error {
	uniqueChannelsIdsMap := make(map[string]struct{}, len(req))
	for _, r := range req {
		if r.ChannelID == r.UserID {
			continue
		}

		if _, ok := uniqueChannelsIdsMap[r.ChannelID]; !ok {
			uniqueChannelsIdsMap[r.ChannelID] = struct{}{}
		}
	}

	if len(uniqueChannelsIdsMap) == 0 {
		return nil
	}

	uniqueChannelsIds := make([]string, 0, len(uniqueChannelsIdsMap))
	for k := range uniqueChannelsIdsMap {
		uniqueChannelsIds = append(uniqueChannelsIds, k)
	}

	var channelsEntities []model.Channels
	if err := c.gorm.WithContext(ctx).Where(
		`"id" IN ?`,
		uniqueChannelsIds,
	).Find(&channelsEntities).Error; err != nil {
		c.logger.Error("cannot get channels entities", logger.Error(err))
		return err
	}

	timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	var wg sync.WaitGroup

	var collectedErrors []error

	for _, r := range req {
		var channelEntity *model.Channels
		for _, channel := range channelsEntities {
			if channel.ID == r.ChannelID {
				channelEntity = &channel
				break
			}
		}

		if channelEntity == nil {
			continue
		}

		if r.ChannelID == r.UserID || channelEntity.BotID == r.UserID {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := c.twitchActions.Ban(
				timeoutCtx,
				twitchactions.BanOpts{
					Duration:       r.BanTime,
					Reason:         r.Reason,
					BroadcasterID:  r.ChannelID,
					UserID:         r.UserID,
					ModeratorID:    channelEntity.BotID,
					IsModerator:    r.IsModerator,
					AddModAfterBan: r.AddModAfterBan,
				},
			); err != nil {
				c.logger.Error("cannot ban user", logger.Error(err))
				collectedErrors = append(collectedErrors, err)
			}
		}()
	}

	wg.Wait()

	if len(collectedErrors) > 0 {
		var gigaError string
		for _, err := range collectedErrors {
			gigaError += err.Error() + " "
		}

		return errors.New(gigaError)
	}

	return nil
}
