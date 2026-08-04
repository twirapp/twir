package resolvers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
)

func TestAuthenticatedUserRegenerateChannelAPIKeyUpdatesSelectedDashboard(t *testing.T) {
	dashboardID := uuid.New()
	sessions := channelAPIKeySessions{dashboardID: dashboardID.String()}
	repository := &channelAPIKeyRepository{}
	resolver := &mutationResolver{Resolver: &Resolver{deps: Deps{
		Sessions:           &sessions,
		ChannelsRepository: repository,
	}}}

	apiKey, err := resolver.AuthenticatedUserRegenerateChannelAPIKey(context.Background())
	if err != nil {
		t.Fatalf("AuthenticatedUserRegenerateChannelAPIKey() error = %v", err)
	}
	if apiKey == "" {
		t.Fatal("generated API key is empty")
	}
	if repository.channelID != dashboardID {
		t.Fatalf("updated channel ID = %s, want %s", repository.channelID, dashboardID)
	}
	if repository.apiKey == nil || *repository.apiKey != apiKey {
		t.Fatalf("updated API key = %v, want %q", repository.apiKey, apiKey)
	}
}

type channelAPIKeySessions struct {
	dashboardID string
}

func (*channelAPIKeySessions) GetAuthenticatedUserModel(context.Context) (*model.Users, error) {
	return nil, nil
}

func (*channelAPIKeySessions) GetCurrentPlatform(context.Context) (string, error) {
	return "", nil
}

func (s *channelAPIKeySessions) GetSelectedDashboard(context.Context) (string, error) {
	return s.dashboardID, nil
}

func (*channelAPIKeySessions) SetSessionSelectedDashboard(context.Context, string) error {
	return nil
}

func (*channelAPIKeySessions) SessionLogout(context.Context) error {
	return nil
}

func (*channelAPIKeySessions) GetChannelFromApiKey(context.Context) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

type channelAPIKeyRepository struct {
	channelsrepository.Repository

	channelID uuid.UUID
	apiKey    *string
}

func (r *channelAPIKeyRepository) Update(
	_ context.Context,
	channelID uuid.UUID,
	input channelsrepository.UpdateInput,
) (channelentity.Channel, error) {
	r.channelID = channelID
	r.apiKey = input.ApiKey
	return channelentity.Channel{ID: channelID, ApiKey: input.ApiKey}, nil
}
