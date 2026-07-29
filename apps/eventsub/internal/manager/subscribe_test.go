package manager

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/kvizyx/twitchy/helix"
)

type subscriptionRoundTripper func(*http.Request) (*http.Response, error)

func (f subscriptionRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func newSubscriptionTestClient(t *testing.T, credential helix.Credential, requestHandler subscriptionRoundTripper) *helix.Client {
	t.Helper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: requestHandler}),
		helix.WithStaticToken(credential),
	)
	if err != nil {
		t.Fatalf("new Helix client: %v", err)
	}
	return client
}

func TestGetConditionForTopicUsesBotIDForChatTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelChatMessage, "broadcaster-123", "bot-456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	chatCondition, ok := condition.(eventsub.ChannelChatMessageCondition)
	if !ok {
		t.Fatalf("expected ChannelChatMessageCondition, got %T", condition)
	}

	if chatCondition.BroadcasterUserId != "broadcaster-123" {
		t.Fatalf("expected broadcaster_user_id=broadcaster-123, got %q", chatCondition.BroadcasterUserId)
	}

	if chatCondition.UserId != "bot-456" {
		t.Fatalf("expected user_id=bot-456, got %q", chatCondition.UserId)
	}
}

func TestGetConditionForTopicUsesBotIDForModeratorTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelModerate, "broadcaster-123", "bot-456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	moderateCondition, ok := condition.(eventsub.ChannelModerateV2Condition)
	if !ok {
		t.Fatalf("expected ChannelModerateV2Condition, got %T", condition)
	}

	if moderateCondition.BroadcasterUserId != "broadcaster-123" {
		t.Fatalf("expected broadcaster_user_id=broadcaster-123, got %q", moderateCondition.BroadcasterUserId)
	}

	if moderateCondition.ModeratorUserId != "bot-456" {
		t.Fatalf("expected moderator_user_id=bot-456, got %q", moderateCondition.ModeratorUserId)
	}
}

func TestGetConditionForTopicUsesBotIDForFollowTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelFollow, "broadcaster-123", "bot-456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	followCondition, ok := condition.(eventsub.ChannelFollowCondition)
	if !ok {
		t.Fatalf("expected ChannelFollowCondition, got %T", condition)
	}

	if followCondition.BroadcasterUserId != "broadcaster-123" {
		t.Fatalf("expected broadcaster_user_id=broadcaster-123, got %q", followCondition.BroadcasterUserId)
	}

	if followCondition.ModeratorUserId != "bot-456" {
		t.Fatalf("expected moderator_user_id=bot-456, got %q", followCondition.ModeratorUserId)
	}
}

func TestSubscribeWithLimitsUsesAppClientAndPreservesConduitRequest(t *testing.T) {
	var requests []helix.CreateEventSubSubscriptionRequest
	appClient := newSubscriptionTestClient(t, helix.Credential{
		AccessToken: "app-token",
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassApp,
	}, func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost || request.URL.Path != "/helix/eventsub/subscriptions" {
			t.Fatalf("request = %s %s", request.Method, request.URL.Path)
		}
		if request.Header.Get("Authorization") != "Bearer app-token" {
			t.Fatalf("authorization = %q", request.Header.Get("Authorization"))
		}
		var body helix.CreateEventSubSubscriptionRequest
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatalf("decode subscription request: %v", err)
		}
		requests = append(requests, body)
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		}, nil
	})
	manager := &Manager{
		newAppTwitchClient: func(context.Context) (*helix.Client, error) { return appClient, nil },
		newBotTwitchClient: func(context.Context, string) (*helix.Client, error) {
			t.Fatal("bot client selected for conduit subscription")
			return nil, nil
		},
	}

	err := manager.SubscribeWithLimits(
		context.Background(),
		eventsub.EventTypeChannelChatMessage,
		eventsub.ConduitTransport{Method: "conduit", ConduitId: "conduit-123"},
		"1",
		"broadcaster-123",
		"bot-456",
	)
	if err != nil {
		t.Fatalf("subscribe with limits: %v", err)
	}
	if len(requests) != 1 {
		t.Fatalf("subscription requests = %d, want 1", len(requests))
	}
	request := requests[0]
	if request.Type != eventsub.EventTypeChannelChatMessage.String() || request.Version != "1" {
		t.Fatalf("subscription type/version = %q/%q", request.Type, request.Version)
	}
	if request.Condition["broadcaster_user_id"] != "broadcaster-123" || request.Condition["user_id"] != "bot-456" {
		t.Fatalf("subscription condition = %#v", request.Condition)
	}
	if request.Transport.Method != helix.EventSubTransportConduit || request.Transport.ConduitID != "conduit-123" {
		t.Fatalf("subscription transport = %#v", request.Transport)
	}
}

func TestSubscribeWithLimitsUsesBotClientForWebsocketTransport(t *testing.T) {
	botClient := newSubscriptionTestClient(t, helix.Credential{
		AccessToken: "bot-token",
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassUser,
		UserID:      "bot-456",
		Scopes:      []helix.AuthorizationScope{helix.ScopeChannelReadSubscriptions},
	}, func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Authorization") != "Bearer bot-token" {
			t.Fatalf("authorization = %q", request.Header.Get("Authorization"))
		}
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		}, nil
	})
	manager := &Manager{
		newAppTwitchClient: func(context.Context) (*helix.Client, error) {
			t.Fatal("app client selected for websocket subscription")
			return nil, nil
		},
		newBotTwitchClient: func(_ context.Context, botID string) (*helix.Client, error) {
			if botID != "bot-456" {
				t.Fatalf("bot ID = %q", botID)
			}
			return botClient, nil
		},
	}

	err := manager.SubscribeWithLimits(
		context.Background(),
		eventsub.EventTypeStreamOnline,
		eventsub.WebsocketTransport{Method: "websocket", SessionId: "session-123"},
		"1",
		"bot-456",
		"bot-456",
	)
	if err != nil {
		t.Fatalf("subscribe with limits: %v", err)
	}
}
