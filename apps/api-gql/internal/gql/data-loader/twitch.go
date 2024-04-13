package data_loader

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

func (c *DataLoader) getHelixUsersByIds(ctx context.Context, ids []string) (
	[]*gqlmodel.TwirUserTwitchInfo,
	[]error,
) {
	users, err := c.cachedTwitchClient.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	mappedUsers := make([]*gqlmodel.TwirUserTwitchInfo, len(users))
	for i, user := range users {
		mappedUsers[i] = &gqlmodel.TwirUserTwitchInfo{
			ID:              user.ID,
			Login:           user.Login,
			DisplayName:     user.DisplayName,
			ProfileImageURL: user.ProfileImageURL,
			Description:     user.Description,
		}
	}

	return mappedUsers, nil
}

func GetHelixUser(ctx context.Context, userID string) (*gqlmodel.TwirUserTwitchInfo, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserLoader.Load(ctx, userID)
}

func GetHelixUsers(ctx context.Context, userIDs []string) ([]*gqlmodel.TwirUserTwitchInfo, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.helixUserLoader.LoadAll(ctx, userIDs)
}
