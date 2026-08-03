package mcp_oauth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appentity "github.com/twirapp/twir/apps/api-gql/internal/entity"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

func TestService_VerifyAccessToken_returnsAuthorizedChannel(t *testing.T) {
	// Given
	ctx := context.Background()
	repository := newFakeRepository()
	client := testClient()
	repository.clients[client.ClientID] = client
	token := testToken(client.ClientID, "access-token", "refresh-token", entity.AllScopes())
	repository.putToken(token)
	service := newTestService(t, repository, true, appentity.User{ID: token.UserID.String()})
	channel := channelentity.Channel{
		ID:       token.ChannelID,
		Bindings: []channelplatformentity.ChannelPlatform{{ID: uuid.New()}},
	}
	service.channels = fakeChannels{channel: channel}

	// When
	grant, err := service.VerifyAccessToken(ctx, "access-token")

	// Then
	if err != nil {
		t.Fatalf("VerifyAccessToken() error = %v", err)
	}
	if grant.Channel.ID != channel.ID || len(grant.Channel.Bindings) != len(channel.Bindings) {
		t.Fatalf("grant channel = %#v, want %#v", grant.Channel, channel)
	}
	if grant.ApprovingUserID != token.UserID {
		t.Fatalf("grant approving user = %s, want %s", grant.ApprovingUserID, token.UserID)
	}
}

func TestService_ProtectedResourceMetadataURL_usesCanonicalOrigin(t *testing.T) {
	// Given
	service := newTestService(t, newFakeRepository(), true, appentity.User{})

	// When
	metadataURL := service.ProtectedResourceMetadataURL()

	// Then
	const want = "https://twir.example/.well-known/oauth-protected-resource/api/mcp"
	if metadataURL != want {
		t.Fatalf("protected resource metadata URL = %q, want %q", metadataURL, want)
	}
}
