package vk

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebSocketTokenClientDiscoverChatChannelUsesDevAPIContract(t *testing.T) {
	// Given
	const (
		accessToken       = "oauth-access-token"
		currentChannelURL = "https://live.vkvideo.ru/main-channel"
		chatChannel       = "channel-chat:opaque-main-channel"
	)
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v1/current_user":
			requestCount++
			if requestCount != 1 {
				t.Fatalf("current user request order = %d, want 1", requestCount)
			}
			if request.Method != http.MethodGet {
				t.Errorf("current user request method = %s, want GET", request.Method)
			}
			if request.URL.RawQuery != "" {
				t.Errorf("current user request query = %q, want empty", request.URL.RawQuery)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer "+accessToken {
				t.Errorf("current user authorization = %q", got)
			}
			body, err := io.ReadAll(request.Body)
			if err != nil {
				t.Fatalf("read current user request body: %v", err)
			}
			if len(body) != 0 {
				t.Errorf("current user request body = %q, want empty", body)
			}

			_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/main-channel"}}}`))
		case "/v1/channels":
			requestCount++
			if requestCount != 2 {
				t.Fatalf("channels request order = %d, want 2", requestCount)
			}
			if request.Method != http.MethodPost {
				t.Errorf("channels request method = %s, want POST", request.Method)
			}
			if request.URL.RawQuery != "" {
				t.Errorf("channels request query = %q, want empty", request.URL.RawQuery)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer "+accessToken {
				t.Errorf("channels authorization = %q", got)
			}
			if got := request.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("channels content type = %q, want application/json", got)
			}
			body, err := io.ReadAll(request.Body)
			if err != nil {
				t.Fatalf("read channels request body: %v", err)
			}
			if got := string(body); got != `{"channels":[{"url":"https://live.vkvideo.ru/main-channel"}]}` {
				t.Errorf("channels request body = %q", got)
			}

			_, _ = writer.Write([]byte(`{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/another-channel","web_socket_channels":{"chat":"channel-chat:other"}}},{"channel":{"url":"https://live.vkvideo.ru/main-channel","web_socket_channels":{"chat":"channel-chat:opaque-main-channel"}}}]}}`))
		default:
			t.Fatalf("unexpected request = %s %s", request.Method, request.URL)
		}
	}))
	defer server.Close()

	client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video WebSocket token client: %v", err)
	}

	// When
	channel, err := client.DiscoverChatChannel(context.Background(), OAuthAccessToken(accessToken))

	// Then
	if err != nil {
		t.Fatalf("discover VK Video chat channel: %v", err)
	}
	if channel != WebSocketChannel(chatChannel) {
		t.Errorf("chat channel = %q", channel)
	}
	if requestCount != 2 {
		t.Errorf("request count = %d, want 2", requestCount)
	}
}

func TestWebSocketTokenClientDiscoverChatChannelRejectsMissingOrBlankValuesWithoutLeakingSecrets(t *testing.T) {
	const (
		accessToken       = "oauth-access-token"
		currentChannelURL = "https://live.vkvideo.ru/main-channel-secret"
		chatChannel       = "channel-chat:opaque-main-channel-secret"
	)
	tests := []struct {
		name                string
		currentUserResponse string
		channelsResponse    string
		wantRequests        int
	}{
		{
			name:                "current channel URL is missing",
			currentUserResponse: `{"data":{"channel":{}}}`,
			wantRequests:        1,
		},
		{
			name:                "current channel URL is blank",
			currentUserResponse: `{"data":{"channel":{"url":" "}}}`,
			wantRequests:        1,
		},
		{
			name:                "current channel is missing from channels response",
			currentUserResponse: `{"data":{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret"}}}`,
			channelsResponse:    `{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/other-channel-secret","web_socket_channels":{"chat":"channel-chat:other-secret"}}}]}}`,
			wantRequests:        2,
		},
		{
			name:                "current channel chat value is missing",
			currentUserResponse: `{"data":{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret"}}}`,
			channelsResponse:    `{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret","web_socket_channels":{}}}]}}`,
			wantRequests:        2,
		},
		{
			name:                "current channel chat value is blank",
			currentUserResponse: `{"data":{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret"}}}`,
			channelsResponse:    `{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret","web_socket_channels":{"chat":" "}}}]}}`,
			wantRequests:        2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				requestCount++
				switch request.URL.Path {
				case "/v1/current_user":
					_, _ = writer.Write([]byte(test.currentUserResponse))
				case "/v1/channels":
					_, _ = writer.Write([]byte(test.channelsResponse))
				default:
					t.Fatalf("unexpected request = %s %s", request.Method, request.URL)
				}
			}))
			defer server.Close()

			client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
			if err != nil {
				t.Fatalf("create VK Video WebSocket token client: %v", err)
			}

			// When
			channel, err := client.DiscoverChatChannel(context.Background(), OAuthAccessToken(accessToken))

			// Then
			if err == nil {
				t.Fatal("expected invalid discovery response to fail")
			}
			if channel != "" {
				t.Errorf("chat channel = %q, want empty", channel)
			}
			if requestCount != test.wantRequests {
				t.Errorf("request count = %d, want %d", requestCount, test.wantRequests)
			}
			secrets := []string{accessToken, currentChannelURL, chatChannel, test.currentUserResponse}
			if test.channelsResponse != "" {
				secrets = append(secrets, test.channelsResponse)
			}
			assertWebSocketTokenErrorIsSafe(t, err, secrets...)
		})
	}
}

func TestWebSocketTokenClientDiscoverChatChannelReturnsSanitizedProviderError(t *testing.T) {
	const (
		accessToken       = "oauth-access-token"
		currentChannelURL = "https://live.vkvideo.ru/main-channel-secret"
		chatChannel       = "channel-chat:opaque-main-channel-secret"
		response          = `{"error_description":"raw provider response secret"}`
	)
	tests := []struct {
		name        string
		failingPath string
		statusCode  int
	}{
		{name: "current user rejects request", failingPath: "/v1/current_user", statusCode: http.StatusUnauthorized},
		{name: "channels rejects request", failingPath: "/v1/channels", statusCode: http.StatusBadGateway},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Path == test.failingPath {
					writer.WriteHeader(test.statusCode)
					_, _ = writer.Write([]byte(response))
					return
				}

				switch request.URL.Path {
				case "/v1/current_user":
					_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret"}}}`))
				case "/v1/channels":
					_, _ = writer.Write([]byte(`{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/main-channel-secret","web_socket_channels":{"chat":"channel-chat:opaque-main-channel-secret"}}}]}}`))
				default:
					t.Fatalf("unexpected request = %s %s", request.Method, request.URL)
				}
			}))
			defer server.Close()

			client, err := NewWebSocketTokenClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
			if err != nil {
				t.Fatalf("create VK Video WebSocket token client: %v", err)
			}

			// When
			_, err = client.DiscoverChatChannel(context.Background(), OAuthAccessToken(accessToken))

			// Then
			if err == nil {
				t.Fatal("expected provider rejection")
			}
			var providerError *ProviderError
			if !errors.As(err, &providerError) {
				t.Fatalf("expected ProviderError, got %T (%v)", err, err)
			}
			if providerError.StatusCode != test.statusCode {
				t.Errorf("provider error status = %d, want %d", providerError.StatusCode, test.statusCode)
			}
			assertWebSocketTokenErrorIsSafe(t, err, accessToken, currentChannelURL, chatChannel, response)
		})
	}
}
