package seventv

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/avast/retry-go/v4"
	dispatchtypes "github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/dispatch_types"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	seventvapi "github.com/twirapp/twir/libs/integrations/seventv/api"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels/model"
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
					var channelModel model.Channel

					if err := c.gorm.WithContext(ctx).
						Model(&model.Channel{}).
						Select("twitch_user_id", "kick_user_id").
						Where("id = ?", channel).
						Scan(&channelModel).Error; err != nil {
						return fmt.Errorf("failed to fetch channel %s: %w", channel, err)
					}

					var platforms []platformentity.Platform
					if channelModel.TwitchUserID != nil {
						platforms = append(platforms, platformentity.PlatformTwitch)
					}
					if channelModel.KickUserID != nil {
						platforms = append(platforms, platformentity.PlatformKick)
					}

					if len(platforms) == 0 {
						return fmt.Errorf("channel %s has no connected platform", channel)
					}

					for _, platform := range platforms {
						var userIDStr string
						switch platform {
						case platformentity.PlatformKick:
							userIDStr = channelModel.KickUserID.String()
						case platformentity.PlatformTwitch:
							userIDStr = channelModel.TwitchUserID.String()
						default:
							return fmt.Errorf("unsupported platform %q for channel %s", platform, channel)
						}

						platformUser, err := c.usersRepo.GetByID(ctx, userIDStr)
						if err != nil {
							return fmt.Errorf(
								"failed to fetch user for channel %s platform %s: %w",
								channel, platform, err,
							)
						}

						var profile any
						switch platform {
						case platformentity.PlatformKick:
							profile, err = c.sevenTvApiClient.GetProfileByKickId(ctx, platformUser.PlatformID)
						case platformentity.PlatformTwitch:
							profile, err = c.sevenTvApiClient.GetProfileByTwitchId(ctx, platformUser.PlatformID)
						}

						if err != nil {
							return fmt.Errorf(
								"failed to fetch 7TV profile for channel %s platform %s: %w",
								channel, platform, err,
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
								Platform:   string(platform),
							},
						)
						if activeEmoteSetID != "" {
							c.emoteSetToChannelID[activeEmoteSetID] = channel
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
