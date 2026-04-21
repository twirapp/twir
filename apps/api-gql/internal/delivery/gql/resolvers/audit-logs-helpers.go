package resolvers

import (
	"context"

	"github.com/google/uuid"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func resolveUserProfile(ctx context.Context, r *Resolver, userID string) (*gqlmodel.TwirUserTwitchInfo, string, error) {
	user, err := r.deps.UsersRepository.GetByID(ctx, userID)
	if err != nil {
		if err == usersmodel.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	if user.Platform == platformentity.PlatformTwitch {
		profile, err := data_loader.GetHelixUserById(ctx, user.PlatformID)
		if err != nil {
			return nil, "", err
		}
		return profile, string(user.Platform), nil
	}

	if user.Login == "" {
		return &gqlmodel.TwirUserTwitchInfo{
			ID:       user.PlatformID,
			Login:    "[twir] not found",
			NotFound: true,
		}, string(user.Platform), nil
	}

	return &gqlmodel.TwirUserTwitchInfo{
		ID:              user.PlatformID,
		Login:           user.Login,
		DisplayName:     user.DisplayName,
		ProfileImageURL: user.Avatar,
		Description:     "",
	}, string(user.Platform), nil
}

func resolveChannelProfile(ctx context.Context, r *Resolver, channelID string) (*gqlmodel.TwirUserTwitchInfo, string, error) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, "", nil
	}

	channel, err := r.deps.ChannelsRepository.GetByID(ctx, parsedID)
	if err != nil {
		if err == channelsrepo.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	if channel.IsNil() {
		return nil, "", nil
	}

	var userID string
	var platform string

	if channel.TwitchUserID != nil {
		userID = channel.TwitchUserID.String()
		platform = "twitch"
	} else if channel.KickUserID != nil {
		userID = channel.KickUserID.String()
		platform = "kick"
	} else {
		return nil, "", nil
	}

	profile, _, err := resolveUserProfile(ctx, r, userID)
	return profile, platform, err
}
