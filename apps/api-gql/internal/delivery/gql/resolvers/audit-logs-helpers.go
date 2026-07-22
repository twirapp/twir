package resolvers

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func resolveUserProfile(ctx context.Context, r *Resolver, userID string) (*gqlmodel.TwirUserTwitchInfo, string, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, "", nil
	}

	user, err := r.deps.UsersRepository.GetByID(ctx, parsedUserID)
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

	channel, err := r.deps.ChannelService.GetChannelByID(ctx, parsedID)
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

	for _, candidate := range platformentity.All() {
		binding, found := channelbinding.Find(channel, candidate)
		if !found {
			continue
		}

		userID = binding.UserID.String()
		platform = candidate.String()
		break
	}
	if userID == "" {
		return nil, "", nil
	}

	profile, _, err := resolveUserProfile(ctx, r, userID)
	return profile, platform, err
}
