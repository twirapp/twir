package seventv

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/avast/retry-go/v4"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emotes_store"
	dispatchtypes "github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/dispatch_types"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	seventvapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	"github.com/twirapp/twir/libs/logger"
	"gorm.io/gorm"
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
					logger.Error(err),
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

	for channelID, platformEntries := range data {
		for _, entry := range platformEntries {
			for _, emoteSetSubType := range emoteSetSubTypes {
				platform := "TWITCH"
				if strings.EqualFold(entry.Platform, "kick") {
					platform = "KICK"
				}

				subMessages = append(
					subMessages,
					connSubscription{
						subType: fmt.Sprintf("emote_set.%s", emoteSetSubType),
						conditions: map[string]string{
							"readCtx":  "channel",
							"id":       channelID,
							"platform": platform,
						},
					},
				)
			}

			if entry.EmoteSetID != "" {
				subMessages = append(
					subMessages,
					connSubscription{
						subType: string(dispatchtypes.UpdateEmoteSet),
						conditions: map[string]string{
							"object_id": entry.EmoteSetID,
						},
					},
				)
			}
		}
	}

	return subMessages
}

type channelWithEmoteSet struct {
	EmoteSetID string
	Platform   string
}

type channelsWithEmotesSetsIds map[string][]channelWithEmoteSet

type channelBindingData struct {
	Platform   platformentity.Platform `gorm:"column:platform"`
	PlatformID string                  `gorm:"column:platform_channel_id"`
}

func buildChannelBindingsQuery(db *gorm.DB, ctx context.Context, channelID string) *gorm.DB {
	return db.
		WithContext(ctx).
		Table("channel_platforms AS cp").
		Select("cp.platform", "cp.platform_channel_id").
		Where("cp.channel_id = ?", channelID).
		Where("cp.platform IN ?", []platformentity.Platform{
			platformentity.PlatformTwitch,
			platformentity.PlatformKick,
		})
}

func (c *Service) getChannelsWithEmotesSets(
	ctx context.Context,
	channelsIDs ...string,
) (channelsWithEmotesSetsIds, error) {
	channelsWithEmoteSets := make(channelsWithEmotesSetsIds)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, channel := range channelsIDs {
		if c.registeredChannelIDs[channel] {
			continue
		}

		channelsWithEmoteSets[channel] = nil
		c.registeredChannelIDs[channel] = true

		wg.Add(1)

		go func(channel string) {
			defer wg.Done()

			retry.Do(
				func() error {
					var bindings []channelBindingData
					if err := buildChannelBindingsQuery(c.gorm, ctx, channel).
						Scan(&bindings).Error; err != nil {
						return fmt.Errorf("failed to fetch channel %s: %w", channel, err)
					}

					if len(bindings) == 0 {
						return fmt.Errorf("channel %s has no connected platform", channel)
					}

					for _, binding := range bindings {
						var profile any
						var err error
						switch binding.Platform {
						case platformentity.PlatformKick:
							profile, err = c.sevenTvApiClient.GetProfileByKickId(ctx, binding.PlatformID)
						case platformentity.PlatformTwitch:
							profile, err = c.sevenTvApiClient.GetProfileByTwitchId(ctx, binding.PlatformID)
						default:
							return fmt.Errorf("unsupported platform %q for channel %s", binding.Platform, channel)
						}

						if err != nil {
							return fmt.Errorf(
								"failed to fetch 7TV profile for channel %s platform %s: %w",
								channel, binding.Platform, err,
							)
						}

						if profile == nil {
							continue
						}

						var activeEmoteSetID string
						switch p := profile.(type) {
						case *seventvapi.GetProfileByTwitchIdResponse:
							if p.Users.UserByConnection != nil && p.Users.UserByConnection.Style.ActiveEmoteSet != nil {
								activeEmoteSetID = p.Users.UserByConnection.Style.ActiveEmoteSet.Id
							}
						case *seventvapi.GetProfileByKickIdResponse:
							if p.Users.UserByConnection != nil && p.Users.UserByConnection.Style.ActiveEmoteSet != nil {
								activeEmoteSetID = p.Users.UserByConnection.Style.ActiveEmoteSet.Id
							}
						default:
							return fmt.Errorf("unexpected profile type %T for channel %s", profile, channel)
						}

						mu.Lock()
						channelsWithEmoteSets[channel] = append(
							channelsWithEmoteSets[channel],
							channelWithEmoteSet{
								EmoteSetID: activeEmoteSetID,
								Platform:   string(binding.Platform),
							},
						)
						if activeEmoteSetID != "" {
							c.emoteSetToChannelID[activeEmoteSetID] = emotes_store.ChannelKey{
								Platform: binding.Platform,
								ID:       binding.PlatformID,
							}
						}
						mu.Unlock()
					}

					return nil
				},
				retry.RetryIf(func(err error) bool { return err != nil }),
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
		if len(ws.Instance.subscriptions)+1 < ws.Instance.maxCapacity &&
			len(ws.Instance.createConnUrl()) < 1500 {
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
