package dataloader

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/cache/twitch"
)

func (c *dataLoader) getHelixUsersByNames(ctx context.Context, names []string) (
	[]*gqlmodel.TwirUserTwitchInfo,
	[]error,
) {
	nonEmptyNames := lo.Filter(
		names, func(id string, _ int) bool {
			return id != ""
		},
	)

	users, err := c.deps.CachedTwitchClient.GetUsersByNames(ctx, nonEmptyNames)
	if err != nil {
		return nil, []error{err}
	}

	mappedUsers := make([]*gqlmodel.TwirUserTwitchInfo, 0, len(users))

	for _, name := range names {
		user, ok := lo.Find(
			users, func(item twitch.TwitchUser) bool {
				return item.Login == name
			},
		)
		if !ok {
			mappedUsers = append(
				mappedUsers, &gqlmodel.TwirUserTwitchInfo{
					ID:          name,
					Login:       "[twir] not found",
					DisplayName: "[Twir] Not Found",
					NotFound:    true,
				},
			)
			continue
		}

		mappedUsers = append(
			mappedUsers, &gqlmodel.TwirUserTwitchInfo{
				ID:              user.ID,
				Login:           user.Login,
				DisplayName:     user.DisplayName,
				ProfileImageURL: user.ProfileImageURL,
				Description:     user.Description,
			},
		)
	}

	return mappedUsers, nil
}

func GetHelixUserByName(ctx context.Context, userName string) (
	*gqlmodel.TwirUserTwitchInfo,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserByNameLoader.Load(ctx, userName)
}

func GetHelixUsersByName(ctx context.Context, userName []string) (
	[]*gqlmodel.TwirUserTwitchInfo,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserByNameLoader.LoadAll(ctx, userName)
}
