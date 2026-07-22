package twitch

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
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type twitchChannelLookupStub struct {
	channel channelsmodel.Channel
}

func (s twitchChannelLookupStub) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return s.channel, nil
}

type twitchCaptureTransport struct {
	calls  int
	method string
	path   string
	query  url.Values
}

func (t *twitchCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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

func newTwitchServiceTestClient(t *testing.T, transport http.RoundTripper) *helix.Client {
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

func mixedTwitchBindings(ownerID uuid.UUID) []channelplatformsmodel.ChannelPlatform {
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

func TestGetChannelChatBadgesUsesSelectedTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	ownerID := uuid.New()
	transport := &twitchCaptureTransport{}
	appClientCalls := 0
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: mixedTwitchBindings(ownerID),
		}},
		newAppClient: func(context.Context) (*helix.Client, error) {
			appClientCalls++
			return newTwitchServiceTestClient(t, transport), nil
		},
	}

	if _, err := service.GetChannelChatBadges(context.Background(), channelID.String()); err != nil {
		t.Fatalf("get channel chat badges: %v", err)
	}
	if appClientCalls != 1 || transport.calls != 1 {
		t.Fatalf("app client calls = %d, HTTP calls = %d, want 1 each", appClientCalls, transport.calls)
	}
	if transport.method != http.MethodGet || transport.path != "/helix/chat/badges" {
		t.Fatalf("request = %s %s, want GET /helix/chat/badges", transport.method, transport.path)
	}
	if got := transport.query.Get("broadcaster_id"); got != "twitch-channel" {
		t.Fatalf("broadcaster ID = %q, want selected Twitch channel", got)
	}
}

func TestGetChannelChatBadgesSkipsMissingTwitchBinding(t *testing.T) {
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: mixedTwitchBindings(uuid.Nil)[:2],
		}},
		newAppClient: func(context.Context) (*helix.Client, error) {
			t.Fatal("app client must not be created without a Twitch binding")
			return nil, nil
		},
	}

	badges, err := service.GetChannelChatBadges(context.Background(), uuid.NewString())
	if err != nil {
		t.Fatalf("get channel chat badges: %v", err)
	}
	if badges != nil {
		t.Fatalf("badges = %#v, want nil", badges)
	}
}

func TestGetChannelChatBadgesUsesEmptySelectedTwitchProviderID(t *testing.T) {
	ownerID := uuid.New()
	bindings := mixedTwitchBindings(ownerID)
	bindings[2].PlatformChannelID = ""
	transport := &twitchCaptureTransport{}
	appClientCalls := 0
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{Bindings: bindings}},
		newAppClient: func(context.Context) (*helix.Client, error) {
			appClientCalls++
			return newTwitchServiceTestClient(t, transport), nil
		},
	}

	if _, err := service.GetChannelChatBadges(context.Background(), uuid.NewString()); err != nil {
		t.Fatalf("get channel chat badges: %v", err)
	}
	if appClientCalls != 1 || transport.calls != 1 {
		t.Fatalf("app client calls = %d, HTTP calls = %d, want 1 each", appClientCalls, transport.calls)
	}
}

func TestSetChannelInformationUsesSelectedTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	ownerID := uuid.New()
	transport := &twitchCaptureTransport{}
	var clientOwnerID uuid.UUID
	categoryID := "game-id"
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: mixedTwitchBindings(ownerID),
		}},
		newUserClient: func(_ context.Context, userID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = userID
			return newTwitchServiceTestClient(t, transport), nil
		},
	}

	if err := service.SetChannelInformation(context.Background(), SetChannelInformationInput{
		ChannelID:  channelID.String(),
		CategoryID: &categoryID,
	}); err != nil {
		t.Fatalf("set channel information: %v", err)
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.method != http.MethodPatch || transport.path != "/helix/channels" {
		t.Fatalf("request = %s %s, want PATCH /helix/channels", transport.method, transport.path)
	}
	if got := transport.query.Get("broadcaster_id"); got != "twitch-channel" {
		t.Fatalf("broadcaster ID = %q, want selected Twitch channel", got)
	}
}

func TestSetChannelInformationRejectsMissingTwitchBinding(t *testing.T) {
	categoryID := "game-id"
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: mixedTwitchBindings(uuid.Nil)[:2],
		}},
		newUserClient: func(context.Context, uuid.UUID) (*helix.Client, error) {
			t.Fatal("user client must not be created without a Twitch binding")
			return nil, nil
		},
	}

	err := service.SetChannelInformation(context.Background(), SetChannelInformationInput{
		ChannelID:  uuid.NewString(),
		CategoryID: &categoryID,
	})
	if err == nil {
		t.Fatal("expected missing Twitch binding error")
	}
}

func TestSetChannelInformationUsesEmptySelectedTwitchProviderID(t *testing.T) {
	ownerID := uuid.New()
	bindings := mixedTwitchBindings(ownerID)
	bindings[2].PlatformChannelID = ""
	transport := &twitchCaptureTransport{}
	var clientOwnerID uuid.UUID
	categoryID := "game-id"
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{Bindings: bindings}},
		newUserClient: func(_ context.Context, userID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = userID
			return newTwitchServiceTestClient(t, transport), nil
		},
	}

	if err := service.SetChannelInformation(context.Background(), SetChannelInformationInput{
		ChannelID:  uuid.NewString(),
		CategoryID: &categoryID,
	}); err != nil {
		t.Fatalf("set channel information: %v", err)
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", transport.calls)
	}
}

func TestGetRewardsByChannelIDUsesSelectedTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	ownerID := uuid.New()
	transport := &twitchCaptureTransport{}
	var clientOwnerID uuid.UUID
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: mixedTwitchBindings(ownerID),
		}},
		newUserClient: func(_ context.Context, userID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = userID
			return newTwitchServiceTestClient(t, transport), nil
		},
	}

	result, err := service.GetRewardsByChannelID(context.Background(), channelID.String())
	if err != nil {
		t.Fatalf("get rewards: %v", err)
	}
	if !result.IsPartnerOrAffiliate {
		t.Fatal("expected a successful Twitch rewards response")
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.method != http.MethodGet || transport.path != "/helix/channel_points/custom_rewards" {
		t.Fatalf("request = %s %s, want GET /helix/channel_points/custom_rewards", transport.method, transport.path)
	}
	if got := transport.query.Get("broadcaster_id"); got != "twitch-channel" {
		t.Fatalf("broadcaster ID = %q, want selected Twitch channel", got)
	}
}

func TestGetRewardsByChannelIDSkipsMissingTwitchBinding(t *testing.T) {
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: mixedTwitchBindings(uuid.Nil)[:2],
		}},
		newUserClient: func(context.Context, uuid.UUID) (*helix.Client, error) {
			t.Fatal("user client must not be created without a Twitch binding")
			return nil, nil
		},
	}

	result, err := service.GetRewardsByChannelID(context.Background(), uuid.NewString())
	if err != nil {
		t.Fatalf("get rewards: %v", err)
	}
	if result.IsPartnerOrAffiliate || result.Rewards != nil {
		t.Fatalf("result = %#v, want zero result", result)
	}
}

func TestGetRewardsByChannelIDUsesEmptySelectedTwitchProviderID(t *testing.T) {
	ownerID := uuid.New()
	bindings := mixedTwitchBindings(ownerID)
	bindings[2].PlatformChannelID = ""
	transport := &twitchCaptureTransport{}
	var clientOwnerID uuid.UUID
	service := &Service{
		channelService: twitchChannelLookupStub{channel: channelsmodel.Channel{Bindings: bindings}},
		newUserClient: func(_ context.Context, userID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = userID
			return newTwitchServiceTestClient(t, transport), nil
		},
	}

	result, err := service.GetRewardsByChannelID(context.Background(), uuid.NewString())
	if err != nil {
		t.Fatalf("get rewards: %v", err)
	}
	if !result.IsPartnerOrAffiliate {
		t.Fatal("expected a successful Twitch rewards response")
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", transport.calls)
	}
}
