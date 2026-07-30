package mappers

import (
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func GraphQLPlatformToEntity(platform gqlmodel.Platform) (platformentity.Platform, error) {
	switch platform {
	case gqlmodel.PlatformTwitch:
		return platformentity.PlatformTwitch, nil
	case gqlmodel.PlatformKick:
		return platformentity.PlatformKick, nil
	case gqlmodel.PlatformVkVideoLive:
		return platformentity.PlatformVKVideoLive, nil
	case gqlmodel.PlatformYoutube:
		return platformentity.PlatformYouTube, nil
	default:
		return "", fmt.Errorf("unknown graphql platform: %s", platform)
	}
}

func GraphQLPlatformsToEntities(platforms []gqlmodel.Platform) ([]platformentity.Platform, error) {
	if len(platforms) == 0 {
		return nil, nil
	}

	result := make([]platformentity.Platform, 0, len(platforms))
	for _, p := range platforms {
		mapped, err := GraphQLPlatformToEntity(p)
		if err != nil {
			return nil, err
		}

		result = append(result, mapped)
	}

	return result, nil
}

func EntityPlatformToGraphQL(platform platformentity.Platform) (gqlmodel.Platform, error) {
	switch platform {
	case platformentity.PlatformTwitch:
		return gqlmodel.PlatformTwitch, nil
	case platformentity.PlatformKick:
		return gqlmodel.PlatformKick, nil
	case platformentity.PlatformVKVideoLive:
		return gqlmodel.PlatformVkVideoLive, nil
	case platformentity.PlatformYouTube:
		return gqlmodel.PlatformYoutube, nil
	default:
		return "", fmt.Errorf("unknown entity platform: %s", platform)
	}
}

func EntityPlatformsToGraphQL(platforms []platformentity.Platform) []gqlmodel.Platform {
	result := make([]gqlmodel.Platform, 0, len(platforms))
	for _, p := range platforms {
		mapped, err := EntityPlatformToGraphQL(p)
		if err != nil {
			continue
		}

		result = append(result, mapped)
	}

	return result
}
