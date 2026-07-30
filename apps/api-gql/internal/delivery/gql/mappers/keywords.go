package mappers

import (
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func KeywordsFrom(k entity.Keyword) (gqlmodel.Keyword, error) {
	platforms := make([]gqlmodel.Platform, 0, len(k.Platforms))
	for _, p := range k.Platforms {
		mappedPlatform, err := EntityPlatformToGraphQL(p)
		if err != nil {
			return gqlmodel.Keyword{}, fmt.Errorf("map keyword platform: %w", err)
		}

		platforms = append(platforms, mappedPlatform)
	}

	return gqlmodel.Keyword{
		ID:                  k.ID,
		Text:                k.Text,
		Response:            &k.Response,
		Enabled:             k.Enabled,
		Cooldown:            k.Cooldown,
		IsReply:             k.IsReply,
		IsRegularExpression: k.IsRegular,
		UsageCount:          k.Usages,
		RolesIds:            k.RolesIDs,
		Platforms:           platforms,
	}, nil
}
