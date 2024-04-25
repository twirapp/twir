package data_loader

import (
	"context"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/cache/twitch"
)

func (c *DataLoader) getHelixUsersByIds(ctx context.Context, ids []string) (
	[]*helix.User,
	[]error,
) {
	nonEmptyIds := lo.Filter(
		ids, func(id string, _ int) bool {
			return id != ""
		},
	)

	users, err := c.cachedTwitchClient.GetUsersByIds(ctx, nonEmptyIds)
	if err != nil {
		return nil, []error{err}
	}

	mappedUsers := make([]*helix.User, 0, len(users))

	for _, id := range ids {
		user, ok := lo.Find(
			users, func(item twitch.TwitchUser) bool {
				return item.ID == id
			},
		)
		if !ok {
			mappedUsers = append(
				mappedUsers, &helix.User{
					ID:          id,
					Login:       "[twir] twitch banned",
					DisplayName: "[Twir] Twitch Banned",
				},
			)
			continue
		}

		mappedUsers = append(mappedUsers, &user.User)
	}

	return mappedUsers, nil
}

func GetHelixUser(ctx context.Context, userID string) (*helix.User, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserLoader.Load(ctx, userID)
}

func GetHelixUsers(ctx context.Context, userIDs []string) ([]*helix.User, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserLoader.LoadAll(ctx, userIDs)
}
