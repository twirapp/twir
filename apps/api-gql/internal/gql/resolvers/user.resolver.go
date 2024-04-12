package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"slices"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// AuthenticatedUser is the resolver for the authenticatedUser field.
func (r *queryResolver) AuthenticatedUser(ctx context.Context) (
	*gqlmodel.AuthenticatedUser,
	error,
) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("not authenticated: %w", err)
	}

	authedUser := &gqlmodel.AuthenticatedUser{
		ID:                user.ID,
		IsBotAdmin:        user.IsBotAdmin,
		IsBanned:          user.IsBanned,
		HideOnLandingPage: user.HideOnLandingPage,
		TwitchProfile:     &gqlmodel.TwirUserTwitchInfo{},
		APIKey:            user.ApiKey,
	}

	if slices.Contains(GetPreloads(ctx), "twitchProfile") {
		twitchProfile, err := r.cachedTwitchClient.GetUserById(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		authedUser.TwitchProfile = &gqlmodel.TwirUserTwitchInfo{
			Login:           twitchProfile.Login,
			DisplayName:     twitchProfile.DisplayName,
			ProfileImageURL: twitchProfile.ProfileImageURL,
			Description:     twitchProfile.Description,
		}
	}

	if user.Channel != nil {
		authedUser.IsEnabled = &user.Channel.IsEnabled
		authedUser.IsBotModerator = &user.Channel.IsBotMod
		authedUser.BotID = &user.Channel.BotID
	}

	return authedUser, nil
}
