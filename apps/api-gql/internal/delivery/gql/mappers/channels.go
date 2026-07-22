package mappers

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func MapChannelModelToGqlPublicUser(
	c model.Channel,
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
