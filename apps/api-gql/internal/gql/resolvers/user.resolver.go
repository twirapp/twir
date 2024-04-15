package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *authenticatedUserResolver) TwitchProfile(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) (*gqlmodel.TwirUserTwitchInfo, error) {
	user, err := data_loader.GetHelixUser(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &gqlmodel.TwirUserTwitchInfo{
		Login:           user.Login,
		DisplayName:     user.DisplayName,
		ProfileImageURL: user.ProfileImageURL,
		Description:     user.Description,
	}, nil
}

// AuthenticatedUser is the resolver for the authenticatedUser field.
func (r *queryResolver) AuthenticatedUser(ctx context.Context) (
	*gqlmodel.AuthenticatedUser,
	error,
) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("not authenticated: %w", err)
	}

	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	authedUser := &gqlmodel.AuthenticatedUser{
		ID:                  user.ID,
		IsBotAdmin:          user.IsBotAdmin,
		IsBanned:            user.IsBanned,
		HideOnLandingPage:   user.HideOnLandingPage,
		TwitchProfile:       &gqlmodel.TwirUserTwitchInfo{},
		APIKey:              user.ApiKey,
		SelectedDashboardID: dashboardId,
	}

	if user.Channel != nil {
		authedUser.IsEnabled = &user.Channel.IsEnabled
		authedUser.IsBotModerator = &user.Channel.IsBotMod
		authedUser.BotID = &user.Channel.BotID
	}

	return authedUser, nil
}

// AuthenticatedUser returns graph.AuthenticatedUserResolver implementation.
func (r *Resolver) AuthenticatedUser() graph.AuthenticatedUserResolver {
	return &authenticatedUserResolver{r}
}

type authenticatedUserResolver struct{ *Resolver }
