package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func ChannelPlatformBindingToGraphQL(
	binding channelplatformentity.ChannelPlatform,
	profile usersmodel.User,
	capabilities platformentity.Capabilities,
) (gqlmodel.ChannelPlatformBinding, error) {
	platform, err := EntityPlatformToGraphQL(binding.Platform)
	if err != nil {
		return gqlmodel.ChannelPlatformBinding{}, err
	}

	result := gqlmodel.ChannelPlatformBinding{
		ID:                  binding.ID,
		Platform:            platform,
		UserID:              binding.UserID,
		PlatformChannelID:   binding.PlatformChannelID,
		Enabled:             binding.Enabled,
		PlatformUserID:      profile.PlatformID,
		PlatformLogin:       profile.Login,
		PlatformDisplayName: profile.DisplayName,
		Capabilities:        PlatformCapabilitiesToGraphQL(capabilities),
	}
	if profile.Avatar != "" {
		result.PlatformAvatar = &profile.Avatar
	}

	return result, nil
}

func ChannelPlatformOptionToGraphQL(option channelplatformservice.Option) (gqlmodel.ChannelPlatformOption, error) {
	platform, err := EntityPlatformToGraphQL(option.Platform)
	if err != nil {
		return gqlmodel.ChannelPlatformOption{}, err
	}

	return gqlmodel.ChannelPlatformOption{
		Platform:     platform,
		Capabilities: PlatformCapabilitiesToGraphQL(option.Capabilities),
	}, nil
}

func PlatformCapabilitiesToGraphQL(capabilities platformentity.Capabilities) []gqlmodel.PlatformCapability {
	result := make([]gqlmodel.PlatformCapability, 0, len(capabilities))
	for _, capability := range capabilities {
		result = append(result, gqlmodel.PlatformCapability{Name: string(capability)})
	}

	return result
}
