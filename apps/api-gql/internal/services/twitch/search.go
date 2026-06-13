package twitch

import (
	"context"
	"sort"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func (c *Service) SearchByName(ctx context.Context, query string) ([]helix.Channel, error) {
	if query == "" {
		return nil, nil
	}

	channels, err := c.cachedTwitchClient.SearchChannels(ctx, query)
	if err != nil {
		return nil, err
	}

	return channels, err
}

type SearchChannelsInput struct {
	Query    string
	TwirOnly bool
}

func (c *Service) SearchChannels(ctx context.Context, input SearchChannelsInput) (
	[]helix.Channel,
	error,
) {
	if input.Query == "" {
		return []helix.Channel{}, nil
	}

	channels, err := c.cachedTwitchClient.SearchChannels(ctx, input.Query)
	if err != nil {
		return nil, err
	}

	// Sort channels by relevance
	sort.Slice(
		channels, func(i, j int) bool {
			name1 := channels[i].BroadcasterLogin
			name2 := channels[j].BroadcasterLogin

			containsName1 := strings.Contains(strings.ToLower(name1), strings.ToLower(input.Query))
			containsName2 := strings.Contains(strings.ToLower(name2), strings.ToLower(input.Query))

			if containsName1 && !containsName2 {
				return true
			} else if !containsName1 && containsName2 {
				return false
			} else {
				return name1 < name2
			}
		},
	)

	if input.TwirOnly {
		existingPlatformIds := make([]string, 0, len(channels))
		for _, channel := range channels {
			user, err := c.usersRepository.GetByPlatformID(ctx, platformentity.PlatformTwitch, channel.ID)
			if err == nil {
				existingPlatformIds = append(existingPlatformIds, user.PlatformID)
			}
		}

		channels = lo.Filter(
			channels, func(channel helix.Channel, _ int) bool {
				return lo.Contains(existingPlatformIds, channel.ID)
			},
		)
	}

	return channels, nil
}
