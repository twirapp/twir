package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type SearchCategoriesInput struct {
	Query string
}

type GetCategoriesInput struct {
	IDs []string
}

type SetChannelInformationInput struct {
	ChannelID  string
	CategoryID *string
	Title      *string
}

func (s *Service) SearchCategories(ctx context.Context, input SearchCategoriesInput) ([]helix.SearchCategory, error) {
	if input.Query == "" {
		return []helix.SearchCategory{}, nil
	}

	twitchClient, err := s.createAppClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create twitch client: %w", err)
	}

	resp, err := twitchClient.Search.SearchCategories(ctx, helix.SearchCategoriesRequest{
		Query: input.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("unexpected error when fetching categories: %w", err)
	}

	return resp.Data, nil
}

func (s *Service) GetCategories(ctx context.Context, input GetCategoriesInput) ([]helix.Game, error) {
	if len(input.IDs) == 0 {
		return []helix.Game{}, nil
	}

	twitchClient, err := s.createAppClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create twitch client: %w", err)
	}

	resp, err := twitchClient.Games.GetGames(ctx, helix.GetGamesRequest{
		IDs: input.IDs,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get categories: %w", err)
	}

	return resp.Data, nil
}

func (s *Service) SetChannelInformation(ctx context.Context, input SetChannelInformationInput) error {
	if input.ChannelID == "" {
		return fmt.Errorf("userId is required")
	}

	if input.CategoryID == nil && input.Title == nil {
		return fmt.Errorf("at least one of categoryID or title must be provided")
	}

	parsedID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := s.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return fmt.Errorf("channel not found or twitch not connected")
	}

	twitchBinding, found := channel.Binding(platformentity.PlatformTwitch)
	if !found || twitchBinding.UserID == uuid.Nil {
		return fmt.Errorf("channel not found or twitch not connected")
	}

	twitchClient, err := s.createUserClient(ctx, twitchBinding.UserID)
	if err != nil {
		return fmt.Errorf("cannot create twitch client for user %s: %w", input.ChannelID, err)
	}

	params := helix.ModifyChannelInformationRequest{
		BroadcasterID: twitchBinding.PlatformChannelID,
	}

	if input.CategoryID != nil {
		params.GameID = input.CategoryID
	}

	if input.Title != nil {
		params.Title = input.Title
	}

	_, err = twitchClient.Channels.ModifyChannelInformation(ctx, params)
	if err != nil {
		return fmt.Errorf("cannot update channel information: %w", err)
	}

	return nil
}
