package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twirapp/twir/libs/integrations/vk"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestTransportProvidesDevAPIWebSocketTokensToCentrifugo(t *testing.T) {
	// Given
	const (
		currentChannelURL = "https://live.vkvideo.ru/current-channel"
		opaqueChatChannel = "opaque-chat:current/with+symbols"
	)
	binding := testBinding()
	oauthTokens := &recordingTokenProvider{}
	requests := make([]string, 0, 4)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer fixture-access-token" {
			t.Errorf("DevAPI authorization = %q", got)
		}
		requests = append(requests, request.Method+" "+request.URL.Path)

		switch request.URL.Path {
		case "/v1/current_user":
			_, _ = fmt.Fprintf(writer, `{"data":{"channel":{"url":%q}}}`, currentChannelURL)
		case "/v1/channels":
			if got := request.Method; got != http.MethodPost {
				t.Errorf("DevAPI channels method = %q, want POST", got)
			}
			_, _ = fmt.Fprintf(writer, `{"data":{"channels":[{"channel":{"url":%q,"web_socket_channels":{"chat":%q}}}]}}`, currentChannelURL, opaqueChatChannel)
		case "/v1/websocket/token":
			_, _ = writer.Write([]byte(`{"data":{"token":"fixture-connection-jwt"}}`))
		case "/v1/websocket/subscription_token":
			if got := request.URL.Query().Get("channels"); got != opaqueChatChannel {
				t.Errorf("DevAPI subscription channel = %q, want %q", got, opaqueChatChannel)
			}
			_, _ = fmt.Fprintf(writer, `{"data":{"channel_tokens":[{"channel":%q,"token":"fixture-subscription-jwt"}]}}`, opaqueChatChannel)
		default:
			t.Errorf("DevAPI path = %q", request.URL.Path)
		}
	}))
	defer server.Close()

	client, err := vk.NewWebSocketTokenClient(vk.VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create DevAPI WebSocket token client: %v", err)
	}
	connection := &recordingConnection{}
	transport := newTestTransport(
		t,
		devAPIWebSocketTokenProvider{oauthTokens: oauthTokens, client: client},
		&recordingPublisher{},
		&recordingPublisher{},
		connection,
		usersmodel.User{},
	)

	// When
	err = transport.Subscribe(context.Background(), binding)

	// Then
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	if connection.channel != opaqueChatChannel {
		t.Errorf("realtime channel = %q, want discovered channel %q", connection.channel, opaqueChatChannel)
	}
	if connection.connectionToken != "fixture-connection-jwt" {
		t.Errorf("connection callback token = %q", connection.connectionToken)
	}
	if connection.subscriptionToken != "fixture-subscription-jwt" {
		t.Errorf("subscription callback token = %q", connection.subscriptionToken)
	}
	if connection.connectionToken == connection.subscriptionToken {
		t.Error("connection and subscription callback tokens must be distinct")
	}
	if got := oauthTokens.UserIDs(); len(got) != 3 || got[0] != binding.UserID || got[1] != binding.UserID || got[2] != binding.UserID {
		t.Errorf("OAuth token user IDs = %v, want binding user three times", got)
	}
	wantRequests := []string{
		"GET /v1/current_user",
		"POST /v1/channels",
		"GET /v1/websocket/token",
		"GET /v1/websocket/subscription_token",
	}
	if fmt.Sprint(requests) != fmt.Sprint(wantRequests) {
		t.Errorf("DevAPI requests = %v, want %v", requests, wantRequests)
	}

	if err := transport.Unsubscribe(context.Background(), binding); err != nil {
		t.Fatalf("unsubscribe: %v", err)
	}
}

func TestTransportDoesNotCreateConnectionWhenChatChannelDiscoveryFails(t *testing.T) {
	// Given
	discoveryErr := errors.New("discovery failed")
	connection := &recordingConnection{}
	transport := newTestTransport(
		t,
		&recordingTokenProvider{discoverErr: discoveryErr},
		&recordingPublisher{},
		&recordingPublisher{},
		connection,
		usersmodel.User{},
	)

	// When
	err := transport.Subscribe(context.Background(), testBinding())

	// Then
	if !errors.Is(err, discoveryErr) {
		t.Fatalf("subscribe error = %v, want discovery error", err)
	}
	if connection.created {
		t.Error("connection was created after discovery failure")
	}
}
