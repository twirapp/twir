package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *greetingResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.Greeting) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.UserID)
}

// GreetingsCreate is the resolver for the greetingsCreate field.
func (r *mutationResolver) GreetingsCreate(ctx context.Context, opts gqlmodel.GreetingsCreateInput) (*gqlmodel.Greeting, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	newGreeting, err := r.greetingsService.Create(
		ctx,
		greetings.CreateInput{
			ChannelID: dashboardId,
			ActorID:   user.ID,
			UserID:    opts.UserID,
			Enabled:   opts.Enabled,
			Text:      opts.Text,
			IsReply:   opts.IsReply,
			Processed: false,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create greeting: %w", err)
	}

	converted := mappers.GreetingEntityTo(newGreeting)
	return &converted, nil
}

// GreetingsUpdate is the resolver for the greetingsUpdate field.
func (r *mutationResolver) GreetingsUpdate(ctx context.Context, id uuid.UUID, opts gqlmodel.GreetingsUpdateInput) (*gqlmodel.Greeting, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	newGreeting, err := r.greetingsService.Update(
		ctx,
		id,
		greetings.UpdateInput{
			ChannelID: dashboardId,
			ActorID:   user.ID,
			UserID:    opts.UserID.Value(),
			Enabled:   opts.Enabled.Value(),
			Text:      opts.Text.Value(),
			IsReply:   opts.IsReply.Value(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot update greeting: %w", err)
	}

	converted := mappers.GreetingEntityTo(newGreeting)
	return &converted, nil
}

// GreetingsRemove is the resolver for the greetingsRemove field.
func (r *mutationResolver) GreetingsRemove(ctx context.Context, id uuid.UUID) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	err = r.greetingsService.Delete(
		ctx, greetings.DeleteInput{
			ChannelID: dashboardId,
			ActorID:   user.ID,
			ID:        id,
		},
	)
	if err != nil {
		return false, fmt.Errorf("cannot delete greeting: %w", err)
	}

	return true, nil
}

// Greetings is the resolver for the greetings field.
func (r *queryResolver) Greetings(ctx context.Context) ([]gqlmodel.Greeting, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entities, err := r.greetingsService.GetManyByChannelID(ctx, dashboardId)
	if err != nil {
		return nil, fmt.Errorf("cannot get greetings: %w", err)
	}

	result := make([]gqlmodel.Greeting, len(entities))
	for i, greeting := range entities {
		result[i] = mappers.GreetingEntityTo(greeting)
	}

	return result, nil
}

// Greeting returns graph.GreetingResolver implementation.
func (r *Resolver) Greeting() graph.GreetingResolver { return &greetingResolver{r} }

type greetingResolver struct{ *Resolver }
