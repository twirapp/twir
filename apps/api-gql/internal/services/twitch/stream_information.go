package twitch

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/libs/oauth"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type SetStreamInformationInput struct {
	ChannelID  string
	ActorID    uuid.UUID
	Platform   platformentity.Platform
	CategoryID *string
	Title      *string
}

type kickUserTokenRequester func(context.Context, uuid.UUID) (oauth.Credential, error)
type kickStreamInformationUpdater func(context.Context, string, *string, *string) error

func (s *Service) SetStreamInformation(ctx context.Context, input SetStreamInformationInput) error {
	if input.CategoryID == nil && input.Title == nil {
		return fmt.Errorf("at least one of categoryID or title must be provided")
	}

	switch input.Platform {
	case platformentity.PlatformTwitch:
		return s.SetChannelInformation(ctx, SetChannelInformationInput{
			ChannelID:  input.ChannelID,
			CategoryID: input.CategoryID,
			Title:      input.Title,
		})
	case platformentity.PlatformKick:
		return s.setKickStreamInformation(ctx, input)
	case platformentity.PlatformVKVideoLive:
		return fmt.Errorf("unsupported platform: %s", input.Platform)
	default:
		return fmt.Errorf("unsupported platform: %s", input.Platform)
	}
}

func (s *Service) setKickStreamInformation(ctx context.Context, input SetStreamInformationInput) error {
	channelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := s.channelService.GetChannelByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return fmt.Errorf("channel not found or Kick not connected")
	}

	kickBinding, found := channel.Binding(platformentity.PlatformKick)
	if !found || kickBinding.UserID == uuid.Nil {
		return fmt.Errorf("channel not found or Kick not connected")
	}
	if input.ActorID == uuid.Nil || kickBinding.UserID != input.ActorID {
		return fmt.Errorf("Kick stream information can only be edited by the channel's Kick account")
	}

	requestToken := s.requestKickUserToken
	if requestToken == nil {
		requestToken = s.defaultKickUserTokenRequester
	}
	token, err := requestToken(ctx, kickBinding.UserID)
	if err != nil {
		return fmt.Errorf("request Kick user token: %w", err)
	}
	if !hasKickChannelWriteScope(token.Scopes) {
		return fmt.Errorf("Kick stream information requires the channel:write scope; reconnect the Kick account")
	}

	update := s.updateKickStreamInformation
	if update == nil {
		update = s.defaultKickStreamInformationUpdater
	}
	if err := update(ctx, token.AccessToken, input.Title, input.CategoryID); err != nil {
		return fmt.Errorf("update Kick stream information: %w", err)
	}

	return nil
}

func (s *Service) defaultKickUserTokenRequester(
	ctx context.Context,
	userID uuid.UUID,
) (oauth.Credential, error) {
	return s.kickUserTokens.Token(ctx, userID)
}

func (s *Service) defaultKickStreamInformationUpdater(
	ctx context.Context,
	accessToken string,
	title *string,
	categoryID *string,
) error {
	if s.kickProvider == nil {
		return fmt.Errorf("Kick provider is not configured")
	}

	return s.kickProvider.UpdateStreamInformation(ctx, accessToken, title, categoryID)
}

func hasKickChannelWriteScope(scopes []string) bool {
	for _, scope := range scopes {
		if strings.TrimSpace(scope) == "channel:write" {
			return true
		}
	}

	return false
}

func (s *Service) SearchKickCategories(ctx context.Context, query string) ([]gokick.CategoryResponse, error) {
	if len(strings.TrimSpace(query)) < 3 {
		return []gokick.CategoryResponse{}, nil
	}
	if s.kickProvider == nil {
		return nil, fmt.Errorf("Kick provider is not configured")
	}

	response, err := s.kickAppTokens.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("request Kick app token: %w", err)
	}

	categories, err := s.kickProvider.SearchCategories(ctx, response.AccessToken, query)
	if err != nil {
		return nil, fmt.Errorf("search Kick categories: %w", err)
	}

	return categories, nil
}
