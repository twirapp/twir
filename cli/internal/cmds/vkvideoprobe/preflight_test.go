package vkvideoprobe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVKClientPreflight_selectsExactActiveChannel(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer synthetic-access-token" {
			t.Errorf("unexpected authorization header")
		}

		switch request.URL.Path {
		case "/v1/channels/active":
			_, _ = writer.Write([]byte(`{"data":[{"channel":{"url":"https://live.vkvideo.ru/other","web_socket_channels":{"chat":"chat-other"}},"stream":{"id":"stream-other"}},{"channel":{"url":"https://live.vkvideo.ru/exact","web_socket_channels":{"chat":"chat-exact"}},"stream":{"id":"stream-exact"}}]}`))
		case "/v1/websocket/token":
			_, _ = writer.Write([]byte(`{"data":{"token":"synthetic-connection-token"}}`))
		case "/v1/websocket/subscription_token":
			if got := request.URL.Query().Get("channels"); got != "chat-exact" {
				t.Errorf("subscription query channel = %q, want chat-exact", got)
			}
			_, _ = writer.Write([]byte(`{"data":{"channel_tokens":[{"channel":"chat-other","token":"synthetic-other-token"},{"channel":"chat-exact","token":"synthetic-subscription-token"}]}}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := NewVKClient(server.Client(), server.URL)
	if err != nil {
		t.Fatalf("new VK client: %v", err)
	}

	// When
	got, err := client.Preflight(context.Background(), "https://live.vkvideo.ru/exact", "synthetic-access-token")

	// Then
	if err != nil {
		t.Fatalf("preflight: %v", err)
	}
	if got.StreamID != "stream-exact" || got.ChatChannel != "chat-exact" {
		t.Fatalf("selected channel = %#v", got)
	}
	if got.ConnectionToken != "synthetic-connection-token" || got.SubscriptionToken != "synthetic-subscription-token" {
		t.Fatalf("unexpected preflight credentials")
	}
}
