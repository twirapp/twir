package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func MapChannelModelToGqlPublicUser(c model.Channel) *gqlmodel.TwirPublicUser {
	u := &gqlmodel.TwirPublicUser{
		ID:                c.ID,
		HideOnLandingPage: false,
		TwitchProfile:     nil,
		KickProfile:       nil,
	}

	if c.TwitchUser != nil {
		u.TwitchProfile = &gqlmodel.TwirUserTwitchInfo{
			ID:              c.TwitchUser.PlatformID,
			Login:           c.TwitchUser.Login,
			DisplayName:     c.TwitchUser.DisplayName,
			ProfileImageURL: c.TwitchUser.Avatar,
			Description:     "",
			NotFound:        false,
		}
	}

	if c.KickUser != nil {
		u.KickProfile = &gqlmodel.KickProfile{
			ID:             c.KickUser.PlatformID,
			Slug:           c.KickUser.Login,
			DisplayName:    c.KickUser.DisplayName,
			ProfilePicture: &c.KickUser.Avatar,
		}
	}

	return u
}
