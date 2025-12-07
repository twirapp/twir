package twitch

import (
	"context"
	"sort"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/repositories/users/model"
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
		channelIds := lo.Map(
			channels, func(channel helix.Channel, _ int) string {
				return channel.ID
			},
		)

		existingUsers, err := c.usersRepository.GetManyByIDS(
			ctx, users.GetManyInput{
				IDs: channelIds,
			},
		)
		if err != nil {
			return nil, err
		}

		existingUserIds := lo.Map(
			existingUsers, func(user model.User, _ int) string {
				return user.ID
			},
		)

		channels = lo.Filter(
			channels, func(channel helix.Channel, _ int) bool {
				return lo.Contains(existingUserIds, channel.ID)
			},
		)
	}

	return channels, nil
}
