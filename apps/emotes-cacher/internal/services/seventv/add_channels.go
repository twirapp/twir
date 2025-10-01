package seventv

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/avast/retry-go/v4"
	dispatchtypes "github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/dispatch_types"
)

func (c *Service) AddChannels(ctx context.Context, channelsIDs ...string) error {
	if len(channelsIDs) == 0 {
		return nil
	}

	channelsWithEmoteSets, err := c.getChannelsWithEmotesSets(ctx, channelsIDs...)
	if err != nil {
		return fmt.Errorf(
			"failed to get channels with emotes sets: %w",
			err,
		)
	}

	subMessages := c.createSubMessages(channelsWithEmoteSets)

	for _, msg := range subMessages {
		instance, err := c.getOrCreateSocketInstance()
		if err != nil {
			return fmt.Errorf(
				"failed to get or create socket instance: %w",
				err,
			)
		}

		for {
			if err := instance.Instance.subscribe(msg); err != nil {
				c.logger.Error(
					"failed to subscribe to 7TV websocket",
					slog.Any("err", err),
				)
				continue
			}

			break
		}
	}

	return nil
}

var emoteSetSubTypes = []string{"create", "delete"}

func (c *Service) createSubMessages(data channelsWithEmotesSetsIds) []connSubscription {
	subMessages := make([]connSubscription, 0, len(data))

	for channelId, emoteSetId := range data {
		for _, emoteSetSubType := range emoteSetSubTypes {
			subMessages = append(
				subMessages,
				connSubscription{
					subType: fmt.Sprintf("emote_set.%s", emoteSetSubType),
					conditions: map[string]string{
						"readCtx":  "channel",
						"id":       channelId,
						"platform": "TWITCH",
					},
				},
			)
		}

		if emoteSetId != "" {
			subMessages = append(
				subMessages,
				connSubscription{
					subType: string(dispatchtypes.UpdateEmoteSet),
					conditions: map[string]string{
						"object_id": emoteSetId,
					},
				},
			)
		}

	}

	return subMessages
}

type channelsWithEmotesSetsIds map[string]string

func (c *Service) getChannelsWithEmotesSets(
	ctx context.Context,
	channelsIDs ...string,
) (channelsWithEmotesSetsIds, error) {
	channelsWithEmoteSets := make(channelsWithEmotesSetsIds)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, channel := range channelsIDs {
		if _, ok := c.registeredChannelsWithEmoteSetId[channel]; ok {
			continue
		}

		channelsWithEmoteSets[channel] = ""
		c.registeredChannelsWithEmoteSetId[channel] = ""

		wg.Add(1)

		go func(channel string) {
			defer wg.Done()
			retry.Do(
				func() error {
					profile, err := c.sevenTvApiClient.GetProfileByTwitchId(ctx, channel)
					if err != nil {
						return fmt.Errorf(
							"failed to fetch profile for channel %s: %w",
							channel,
							err,
						)
					}
					if profile == nil || profile.Users.UserByConnection == nil || profile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
						return nil
					}

					mu.Lock()
					channelsWithEmoteSets[channel] = profile.Users.UserByConnection.Style.ActiveEmoteSet.Id
					c.registeredChannelsWithEmoteSetId[channel] = profile.Users.UserByConnection.Style.ActiveEmoteSet.Id
					mu.Unlock()

					return nil
				},
				retry.RetryIf(
					func(err error) bool {
						return err != nil
					},
				),
				retry.Attempts(10),
			)
		}(channel)
	}

	wg.Wait()

	return channelsWithEmoteSets, nil
}

var socketsMu sync.Mutex

func (c *Service) getOrCreateSocketInstance() (*socketInstance, error) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	var instance *socketInstance

	// find free socket instance
	for _, ws := range c.sockets {
		if len(ws.Instance.subscriptions)+1 < ws.Instance.maxCapacity {
			instance = ws
			break
		}
	}

	if instance == nil {
		newConn := newConn(
			c.onMessage,
			350,
		)
		instance = &socketInstance{
			Instance: newConn,
			ShardID:  uint8(len(c.sockets)),
		}

		c.sockets = append(
			c.sockets,
			instance,
		)
	}

	return instance, nil
}
