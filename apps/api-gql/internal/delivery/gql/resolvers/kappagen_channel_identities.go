package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

func mapKappagenChannelIdentitiesForResolver(ctx context.Context, r *Resolver, rawChannelID string) ([]gqlmodel.KappagenChannelIdentity, error) {
	channelID, err := uuid.Parse(rawChannelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	identities, err := r.deps.ChannelsService.GetPlatformIdentities(ctx, channelID)
	if err != nil {
		return nil, err
	}

	result := make([]gqlmodel.KappagenChannelIdentity, 0, len(identities))
	for _, identity := range identities {
		platform, err := mappers.EntityPlatformToGraphQL(identity.Platform)
		if err != nil {
			return nil, err
		}

		result = append(result, gqlmodel.KappagenChannelIdentity{
			Platform: platform,
			ID:       identity.ID,
		})
	}

	return result, nil
}
