package twitch

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
)

type SearchCategoriesInput struct {
	Query   string
}

type GetCategoriesInput struct {
	IDs    []string
}

type SetChannelInformationInput struct {
	ChannelID     string
	CategoryID *string
	Title      *string
}

func (s *Service) SearchCategories(ctx context.Context, input SearchCategoriesInput) ([]helix.Category, error) {
	if input.Query == "" {
		return []helix.Category{}, nil
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		s.config,
		s.twirBus,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create twitch client: %w", err)
	}

	resp, err := twitchClient.SearchCategories(&helix.SearchCategoriesParams{
		Query: input.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("unexpected error when fetching categories: %w", err)
	}

	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("cannot get categories: %s", resp.ErrorMessage)
	}

	return resp.Data.Categories, nil
}

func (s *Service) GetCategories(ctx context.Context, input GetCategoriesInput) ([]helix.Game, error) {
	if len(input.IDs) == 0 {
		return []helix.Game{}, nil
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		s.config,
		s.twirBus,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create twitch client: %w", err)
	}

	resp, err := twitchClient.GetGames(&helix.GamesParams{
		IDs: input.IDs,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get categories: %w", err)
	}

	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("cannot get categories: %s", resp.ErrorMessage)
	}

	return resp.Data.Games, nil
}

func (s *Service) SetChannelInformation(ctx context.Context, input SetChannelInformationInput) error {
	if input.ChannelID == "" {
		return fmt.Errorf("userId is required")
	}

	if input.CategoryID == nil && input.Title == nil {
		return fmt.Errorf("at least one of categoryID or title must be provided")
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		input.ChannelID,
		s.config,
		s.twirBus,
	)
	if err != nil {
		return fmt.Errorf("cannot create twitch client for user %s: %w", input.ChannelID, err)
	}

	params := &helix.EditChannelInformationParams{
		BroadcasterID: input.ChannelID,
	}

	if input.CategoryID != nil {
		params.GameID = *input.CategoryID
	}

	if input.Title != nil {
		params.Title = *input.Title
	}

	resp, err := twitchClient.EditChannelInformation(params)
	if err != nil {
		return fmt.Errorf("cannot update channel information: %w", err)
	}

	if resp.ErrorMessage != "" {
		return fmt.Errorf("cannot update channel information: %s", resp.ErrorMessage)
	}

	return nil
}
