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
	users, err := c.cachedTwitchClient.GetUsersByIds(ctx, ids)
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
					Login:       "[Twir] Twitch Banned",
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
