package vk

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebSocketTokenClientConnectionTokenUsesDevAPIContract(t *testing.T) {
	// Given
	const accessToken = "oauth-access-token"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/v1/websocket/token" {
			t.Fatalf("connection token request = %s %s", request.Method, request.URL)
		}
		if request.URL.RawQuery != "" {
			t.Errorf("connection token query = %q, want empty", request.URL.RawQuery)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer "+accessToken {
			t.Errorf("connection token authorization = %q", got)
		}

		_, _ = writer.Write([]byte(`{"data":{"token":"connection-token"}}`))
	}))
	defer server.Close()

	client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video WebSocket token client: %v", err)
	}

	// When
	token, err := client.ConnectionToken(context.Background(), OAuthAccessToken(accessToken))

	// Then
	if err != nil {
		t.Fatalf("get VK Video WebSocket connection token: %v", err)
	}
	if token != WebSocketConnectionToken("connection-token") {
		t.Errorf("connection token = %q", token)
	}
}

func TestWebSocketTokenClientSubscriptionTokenSelectsRequestedChannelToken(t *testing.T) {
	// Given
	const (
		accessToken = "oauth-access-token"
		channel     = "requested-channel"
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/v1/websocket/subscription_token" {
			t.Fatalf("subscription token request = %s %s", request.Method, request.URL)
		}
		if got := request.URL.Query().Get("channels"); got != channel {
			t.Errorf("subscription token channels query = %q", got)
		}
		if len(request.URL.Query()) != 1 {
			t.Errorf("subscription token query = %q", request.URL.RawQuery)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer "+accessToken {
			t.Errorf("subscription token authorization = %q", got)
		}

		_, _ = writer.Write([]byte(`{"data":{"channel_tokens":[{"channel":"other-channel","token":"other-token"},{"channel":"requested-channel","token":"subscription-token"}]}}`))
	}))
	defer server.Close()

	client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video WebSocket token client: %v", err)
	}

	// When
	token, err := client.SubscriptionToken(context.Background(), OAuthAccessToken(accessToken), WebSocketChannel(channel))

	// Then
	if err != nil {
		t.Fatalf("get VK Video WebSocket subscription token: %v", err)
	}
	if token != WebSocketSubscriptionToken("subscription-token") {
		t.Errorf("subscription token = %q", token)
	}
}

func TestWebSocketTokenClientRejectsMissingTokensWithoutLeakingSecrets(t *testing.T) {
	const (
		accessToken = "oauth-access-token"
		channel     = "requested-channel"
	)
	tests := []struct {
		name     string
		response string
		request  func(*WebSocketTokenClient) error
	}{
		{
			name:     "connection token is missing",
			response: `{"data":{}}`,
			request: func(client *WebSocketTokenClient) error {
				_, err := client.ConnectionToken(context.Background(), OAuthAccessToken(accessToken))
				return err
			},
		},
		{
			name:     "connection token is blank",
			response: `{"data":{"token":" "}}`,
			request: func(client *WebSocketTokenClient) error {
				_, err := client.ConnectionToken(context.Background(), OAuthAccessToken(accessToken))
				return err
			},
		},
		{
			name:     "requested subscription token is blank",
			response: `{"data":{"channel_tokens":[{"channel":"requested-channel","token":" "}]}}`,
			request: func(client *WebSocketTokenClient) error {
				_, err := client.SubscriptionToken(context.Background(), OAuthAccessToken(accessToken), WebSocketChannel(channel))
				return err
			},
		},
		{
			name:     "requested subscription channel is missing",
			response: `{"data":{"channel_tokens":[{"channel":"other-channel","token":"other-token"}]}}`,
			request: func(client *WebSocketTokenClient) error {
				_, err := client.SubscriptionToken(context.Background(), OAuthAccessToken(accessToken), WebSocketChannel(channel))
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(test.response))
			}))
			defer server.Close()

			client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
			if err != nil {
				t.Fatalf("create VK Video WebSocket token client: %v", err)
			}

			// When
			err = test.request(client)

			// Then
			if err == nil {
				t.Fatal("expected invalid token response to fail")
			}
			assertWebSocketTokenErrorIsSafe(t, err, accessToken, channel, test.response)
		})
	}
}

func TestWebSocketTokenClientReturnsSanitizedProviderError(t *testing.T) {
	// Given
	const (
		accessToken = "oauth-access-token"
		channel     = "requested-channel"
		response    = `{"error_description":"raw provider response"}`
	)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video WebSocket token client: %v", err)
	}

	// When
	_, err = client.SubscriptionToken(context.Background(), OAuthAccessToken(accessToken), WebSocketChannel(channel))

	// Then
	if err == nil {
		t.Fatal("expected provider rejection")
	}
	var providerError *ProviderError
	if !errors.As(err, &providerError) {
		t.Fatalf("expected ProviderError, got %T (%v)", err, err)
	}
	if providerError.StatusCode != http.StatusBadRequest {
		t.Errorf("provider error status = %d", providerError.StatusCode)
	}
	assertWebSocketTokenErrorIsSafe(t, err, accessToken, channel, response)
}

func assertWebSocketTokenErrorIsSafe(t *testing.T, err error, secrets ...string) {
	t.Helper()

	for _, secret := range secrets {
		if strings.Contains(err.Error(), secret) {
			t.Errorf("unsafe WebSocket token error = %q", err)
		}
	}
}
