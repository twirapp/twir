package mappers

import (
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func PlatformsToStrings(platforms []platformentity.Platform) []string {
	if platforms == nil {
		return []string{}
	}
	result := make([]string, len(platforms))
	for i, p := range platforms {
		result[i] = p.String()
	}
	return result
}

func StringsToPlatforms(ss []string) []platformentity.Platform {
	if ss == nil {
		return []platformentity.Platform{}
	}
	result := make([]platformentity.Platform, len(ss))
	for i, s := range ss {
		result[i] = platformentity.Platform(s)
	}
	return result
}

func GraphQLPlatformToEntity(platform gqlmodel.Platform) (platformentity.Platform, error) {
	switch platform {
	case gqlmodel.PlatformTwitch:
		return platformentity.PlatformTwitch, nil
	case gqlmodel.PlatformKick:
		return platformentity.PlatformKick, nil
	case gqlmodel.PlatformVkVideoLive:
		return platformentity.PlatformVKVideoLive, nil
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
	default:
		return "", fmt.Errorf("unknown entity platform: %s", platform)
	}
}
