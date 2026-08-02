package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func TestNewServerRegistersToolSchemas(t *testing.T) {
	handler := &Handler{}
	server := handler.newServer(scope{Channel: channelentity.Channel{ID: uuid.New()}})
	if server == nil {
		t.Fatal("expected MCP server")
	}
}

func TestHandlerRequiresAPIKeyHeader(t *testing.T) {
	handler := &Handler{}
	request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
	if !strings.Contains(response.Body.String(), "Api-Key header is required") {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestHandlerAcceptsLegacyUserAPIKey(t *testing.T) {
	userID := uuid.New()
	channel := channelentity.Channel{ID: uuid.New(), Bindings: []channelplatformentity.ChannelPlatform{{
		Platform: platform.PlatformTwitch,
		UserID:   userID,
	}}}

	channels := channelservice.NewChannelService(
		&mcpLegacyChannelsRepository{channel: channel},
		&buscore.Bus{},
		config.Config{},
		nil,
		nil,
	)
	usersService := users.New(users.Opts{
		UsersRepository: &mcpLegacyUsersRepository{user: usersmodel.User{ID: userID, Platform: platform.PlatformTwitch}},
		ChannelService:  channels,
	})
	handler := &Handler{
		deps: Deps{Channels: channels, Users: usersService},
		transport: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestScope, ok := r.Context().Value(contextKey{}).(scope)
			if !ok {
				t.Fatal("request scope was not set")
			}
			if requestScope.Channel.ID != channel.ID {
				t.Fatalf("channel ID = %s, want %s", requestScope.Channel.ID, channel.ID)
			}
			w.WriteHeader(http.StatusNoContent)
		}),
	}
	request := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	request.Header.Set("Api-Key", "legacy-user-api-key")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Deprecation"); got != "true" {
		t.Fatalf("Deprecation header = %q, want true", got)
	}
}

type mcpLegacyChannelsRepository struct {
	channelsrepository.Repository

	channel channelentity.Channel
}

func (*mcpLegacyChannelsRepository) GetByApiKey(context.Context, string) (channelentity.Channel, error) {
	return channelentity.Nil, channelsrepository.ErrNotFound
}

func (r *mcpLegacyChannelsRepository) GetByBindingUserID(
	context.Context,
	platform.Platform,
	uuid.UUID,
) (channelentity.Channel, error) {
	return r.channel, nil
}

type mcpLegacyUsersRepository struct {
	usersrepository.Repository

	user usersmodel.User
}

func (r *mcpLegacyUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return r.user, nil
}
