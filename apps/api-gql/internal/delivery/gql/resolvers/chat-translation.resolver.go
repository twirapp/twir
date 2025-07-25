package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"errors"

	ulid "github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_translation"
	repo "github.com/twirapp/twir/libs/repositories/chat_translation"
)

// ChatTranslationCreate is the resolver for the chatTranslationCreate field.
func (r *mutationResolver) ChatTranslationCreate(ctx context.Context, input gqlmodel.ChatTranslationCreateInput) (*gqlmodel.ChatTranslation, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	createInput := chat_translation.CreateInput{
		ChannelID:         dashboardID,
		ActorID:           user.ID,
		Enabled:           input.Enabled,
		TargetLanguage:    input.TargetLanguage,
		ExcludedLanguages: input.ExcludedLanguages,
		UseItalic:         input.UseItalic,
		ExcludedUsersIDs:  input.ExcludedUsersIDs,
	}

	translation, err := r.deps.ChatTranslationService.Create(ctx, createInput)
	if err != nil {
		return nil, err
	}

	result := mappers.ChatTranslationEntityTo(translation)
	return &result, nil
}

// ChatTranslationUpdate is the resolver for the chatTranslationUpdate field.
func (r *mutationResolver) ChatTranslationUpdate(ctx context.Context, id string, input gqlmodel.ChatTranslationUpdateInput) (*gqlmodel.ChatTranslation, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	updateInput := chat_translation.UpdateInput{
		ActorID:   user.ID,
		ChannelID: dashboardID,
	}

	if input.Enabled.IsSet() {
		updateInput.Enabled = input.Enabled.Value()
	}

	if input.TargetLanguage.IsSet() {
		updateInput.TargetLanguage = input.TargetLanguage.Value()
	}

	if input.ExcludedLanguages.IsSet() {
		langs := input.ExcludedLanguages.Value()
		updateInput.ExcludedLanguages = &langs
	}
	if input.ExcludedUsersIDs.IsSet() {
		ids := input.ExcludedUsersIDs.Value()
		updateInput.ExcludedUsersIDs = &ids
	}

	if input.UseItalic.IsSet() {
		updateInput.UseItalic = input.UseItalic.Value()
	}

	parsedId, err := ulid.Parse(id)
	if err != nil {
		return nil, err
	}

	translation, err := r.deps.ChatTranslationService.Update(ctx, parsedId, updateInput)
	if err != nil {
		return nil, err
	}

	result := mappers.ChatTranslationEntityTo(translation)
	return &result, nil
}

// ChatTranslationDelete is the resolver for the chatTranslationDelete field.
func (r *mutationResolver) ChatTranslationDelete(ctx context.Context, id string) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	parsedId, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}

	if err := r.deps.ChatTranslationService.Delete(
		ctx, chat_translation.DeleteInput{
			ID:        parsedId,
			ChannelID: dashboardID,
			ActorID:   user.ID,
		},
	); err != nil {
		return false, err
	}

	return true, nil
}

// ChatTranslation is the resolver for the chatTranslation field.
func (r *queryResolver) ChatTranslation(ctx context.Context) (*gqlmodel.ChatTranslation, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	translation, err := r.deps.ChatTranslationService.GetByChannelID(ctx, dashboardID)
	if err != nil {
		if errors.Is(err, repo.ErrSettingsNotFound) {
			return nil, nil
		}
		return nil, err
	}

	result := mappers.ChatTranslationEntityTo(translation)
	return &result, nil
}
