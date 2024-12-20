package dataloader

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func (c *dataLoader) getTwitchCategoriesByIDs(ctx context.Context, ids []string) (
	[]*gqlmodel.TwitchCategory,
	[]error,
) {
	result := make([]*gqlmodel.TwitchCategory, len(ids))

	games, err := c.deps.CachedTwitchClient.GetGames(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	for i, id := range ids {
		for _, game := range games {
			if game.ID == id {
				result[i] = &gqlmodel.TwitchCategory{
					ID:        game.ID,
					Name:      game.Name,
					BoxArtURL: game.BoxArtURL,
				}
				break
			}
		}
	}

	return result, nil
}

func GetTwitchCategoryByID(ctx context.Context, userID string) (*gqlmodel.TwitchCategory, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.twitchCategoriesByIdLoader.Load(ctx, userID)
}

func GetTwitchCategoriesByIDs(ctx context.Context, userIDs []string) (
	[]*gqlmodel.TwitchCategory,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.twitchCategoriesByIdLoader.LoadAll(ctx, userIDs)
}
