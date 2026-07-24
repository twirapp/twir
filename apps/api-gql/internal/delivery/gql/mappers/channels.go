package mappers

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func MapChannelModelToGqlPublicUser(
	c channelentity.Channel,
	profiles map[uuid.UUID]usersmodel.User,
) *gqlmodel.TwirPublicUser {
	u := &gqlmodel.TwirPublicUser{
		ID:                c.ID,
		HideOnLandingPage: false,
		TwitchProfile:     nil,
		KickProfile:       nil,
	}

	for _, binding := range c.Bindings {
		profile, ok := profiles[binding.UserID]
		if !ok {
			continue
		}

		switch binding.Platform {
		case platformentity.PlatformTwitch:
			u.TwitchProfile = &gqlmodel.TwirUserTwitchInfo{
				ID:              profile.PlatformID,
				Login:           profile.Login,
				DisplayName:     profile.DisplayName,
				ProfileImageURL: profile.Avatar,
				Description:     "",
				NotFound:        false,
			}
		case platformentity.PlatformKick:
			kickProfile := &gqlmodel.KickProfile{
				ID:          profile.PlatformID,
				Slug:        profile.Login,
				DisplayName: profile.DisplayName,
			}
			if profile.Avatar != "" {
				kickProfile.ProfilePicture = &profile.Avatar
			}
			u.KickProfile = kickProfile
		}
	}

	return u
}
