package users

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/entities/platform"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type usersChannelLookupStub struct {
	channel channelsmodel.Channel
}

func (s usersChannelLookupStub) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return s.channel, nil
}

type usersRepositoryStub struct {
	usersrepository.Repository
	user usersmodel.User
}

func (s usersRepositoryStub) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return s.user, nil
}

type usersCaptureTransport struct {
	calls  int
	method string
	path   string
	query  url.Values
}

func (t *usersCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.calls++
	t.method = req.Method
	t.path = req.URL.Path
	t.query = req.URL.Query()

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		Request:    req,
	}, nil
}

func newUsersServiceTestClient(t *testing.T, transport http.RoundTripper) *helix.Client {
	t.Helper()

	client, err := helix.NewClient(&helix.Options{
		ClientID: "test-client",
		HTTPClient: &http.Client{
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("new Helix client: %v", err)
	}

	return client
}

func usersMixedTwitchBindings(ownerID uuid.UUID) []channelplatformsmodel.ChannelPlatform {
	return []channelplatformsmodel.ChannelPlatform{
		{
			Platform:          platform.PlatformVKVideoLive,
			UserID:            uuid.New(),
			PlatformChannelID: "vk-channel",
			BotConfig:         json.RawMessage(`{`),
		},
		{
			Platform:          platform.PlatformKick,
			UserID:            uuid.New(),
			PlatformChannelID: "kick-channel",
			BotConfig:         json.RawMessage(`{`),
		},
		{
			Platform:          platform.PlatformTwitch,
			UserID:            ownerID,
			PlatformChannelID: "twitch-channel",
			BotConfig:         json.RawMessage(`{`),
		},
	}
}

func TestGetChannelUserInfoUsesSelectedTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()
	transport := &usersCaptureTransport{}
	var clientOwnerID uuid.UUID
	service := &Service{
		channelService: usersChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: usersMixedTwitchBindings(ownerID),
		}},
		usersRepository: usersRepositoryStub{user: usersmodel.User{
			ID:         userID,
			PlatformID: "viewer-platform-user",
		}},
		loadChannelUserInfo: func(context.Context, ChannelUserInfoInput) (deprecatedgormmodel.Users, error) {
			return deprecatedgormmodel.Users{ID: userID.String()}, nil
		},
		newUserClient: func(_ context.Context, gotOwnerID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = gotOwnerID
			return newUsersServiceTestClient(t, transport), nil
		},
	}

	info, err := service.GetChannelUserInfo(context.Background(), ChannelUserInfoInput{
		ChannelID: channelID.String(),
		UserID:    userID.String(),
	})
	if err != nil {
		t.Fatalf("get channel user info: %v", err)
	}
	if info.ID != userID.String() {
		t.Fatalf("info user ID = %q, want %q", info.ID, userID.String())
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.method != http.MethodGet || transport.path != "/helix/channels/followers" {
		t.Fatalf("request = %s %s, want GET /helix/channels/followers", transport.method, transport.path)
	}
	if got := transport.query.Get("broadcaster_id"); got != "twitch-channel" {
		t.Fatalf("broadcaster ID = %q, want selected Twitch channel", got)
	}
	if got := transport.query.Get("user_id"); got != "viewer-platform-user" {
		t.Fatalf("follow user ID = %q, want viewer platform user", got)
	}
}

func TestGetChannelUserInfoRejectsMissingTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	userID := uuid.New()
	loaderCalls := 0
	clientCalls := 0
	service := &Service{
		channelService: usersChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: usersMixedTwitchBindings(uuid.New())[:2],
		}},
		loadChannelUserInfo: func(context.Context, ChannelUserInfoInput) (deprecatedgormmodel.Users, error) {
			loaderCalls++
			return deprecatedgormmodel.Users{}, nil
		},
		newUserClient: func(context.Context, uuid.UUID) (*helix.Client, error) {
			clientCalls++
			return nil, nil
		},
	}

	_, err := service.GetChannelUserInfo(context.Background(), ChannelUserInfoInput{
		ChannelID: channelID.String(),
		UserID:    userID.String(),
	})
	if err == nil {
		t.Fatal("expected missing Twitch binding error")
	}
	if loaderCalls != 0 {
		t.Fatalf("database loader calls = %d, want 0", loaderCalls)
	}
	if clientCalls != 0 {
		t.Fatalf("user client calls = %d, want 0", clientCalls)
	}
}

func TestGetChannelUserInfoUsesEmptySelectedTwitchProviderID(t *testing.T) {
	channelID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()
	bindings := usersMixedTwitchBindings(ownerID)
	bindings[2].PlatformChannelID = ""
	transport := &usersCaptureTransport{}
	loaderCalls := 0
	var clientOwnerID uuid.UUID
	service := &Service{
		channelService: usersChannelLookupStub{channel: channelsmodel.Channel{Bindings: bindings}},
		usersRepository: usersRepositoryStub{user: usersmodel.User{
			ID:         userID,
			PlatformID: "viewer-platform-user",
		}},
		loadChannelUserInfo: func(context.Context, ChannelUserInfoInput) (deprecatedgormmodel.Users, error) {
			loaderCalls++
			return deprecatedgormmodel.Users{ID: userID.String()}, nil
		},
		newUserClient: func(_ context.Context, gotOwnerID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = gotOwnerID
			return newUsersServiceTestClient(t, transport), nil
		},
	}

	if _, err := service.GetChannelUserInfo(context.Background(), ChannelUserInfoInput{
		ChannelID: channelID.String(),
		UserID:    userID.String(),
	}); err != nil {
		t.Fatalf("get channel user info: %v", err)
	}
	if loaderCalls != 1 {
		t.Fatalf("database loader calls = %d, want 1", loaderCalls)
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", transport.calls)
	}
}
