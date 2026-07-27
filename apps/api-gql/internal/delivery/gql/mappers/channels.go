package mappers

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func PlatformProfileToLinkedAccount(platform platformentity.Platform, profile usersmodel.User) gqlmodel.LinkedAccount {
	displayName := profile.DisplayName
	if displayName == "" {
		displayName = profile.Login
	}

	account := gqlmodel.LinkedAccount{
		Platform:            platform.String(),
		PlatformUserID:      profile.PlatformID,
		PlatformLogin:       profile.Login,
		PlatformDisplayName: displayName,
	}
	if profile.Avatar != "" {
		account.PlatformAvatar = &profile.Avatar
	}

	return account
}

func MapChannelModelToGqlPublicUser(
	c channelentity.Channel,
	profiles map[uuid.UUID]usersmodel.User,
	requested *platformentity.Platform,
) *gqlmodel.TwirPublicUser {
	u := &gqlmodel.TwirPublicUser{
		ID:                c.ID,
		HideOnLandingPage: false,
		Profile:           nil,
	}

	var selected *channelplatformentity.ChannelPlatform
	for i := range c.Bindings {
		if requested != nil && c.Bindings[i].Platform == *requested {
			selected = &c.Bindings[i]
			break
		}
	}
	if selected == nil && len(c.Bindings) > 0 {
		selected = &c.Bindings[0]
	}
	if selected != nil {
		if profile, ok := profiles[selected.UserID]; ok {
			account := PlatformProfileToLinkedAccount(selected.Platform, profile)
			u.Profile = &account
		}
	}

	return u
}
