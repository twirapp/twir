package vk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestVideoChatClientSendsTextMessageUsingOwnerAndBotTokens(t *testing.T) {
	// Given
	channelURL := "https://live.vkvideo.ru/channel?tag=alpha&tag=beta"
	streamID := "stream / id&1"
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests[request.URL.Path]++

		switch request.URL.Path {
		case "/configured-api/v1/current_user":
			if request.Method != http.MethodGet {
				t.Errorf("current user request method = %s", request.Method)
			}
			if request.URL.RawQuery != "" {
				t.Errorf("current user query = %q", request.URL.RawQuery)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer owner-token" {
				t.Errorf("current user authorization = %q", got)
			}
			_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/channel?tag=alpha&tag=beta"}}}`))
		case "/configured-api/v1/channel":
			if request.Method != http.MethodGet {
				t.Errorf("channel request method = %s", request.Method)
			}
			if want := (url.Values{"channel_url": []string{channelURL}}).Encode(); request.URL.RawQuery != want {
				t.Errorf("channel query = %q, want %q", request.URL.RawQuery, want)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer bot-token" {
				t.Errorf("channel authorization = %q", got)
			}
			_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/channel?tag=alpha&tag=beta"},"stream":{"id":"stream / id&1"}}}`))
		case "/configured-api/v1/chat/message/send":
			if request.Method != http.MethodPost {
				t.Errorf("send request method = %s", request.Method)
			}
			if want := (url.Values{"channel_url": []string{channelURL}, "stream_id": []string{streamID}}).Encode(); request.URL.RawQuery != want {
				t.Errorf("send query = %q, want %q", request.URL.RawQuery, want)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer bot-token" {
				t.Errorf("send authorization = %q", got)
			}
			if got := request.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("send content type = %q", got)
			}

			var body struct {
				Parts []struct {
					Text struct {
						Content string `json:"content"`
					} `json:"text"`
				} `json:"parts"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatalf("decode send body: %v", err)
			}
			if len(body.Parts) != 1 || body.Parts[0].Text.Content != "hello from bot" {
				t.Errorf("send body = %#v", body)
			}
		default:
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL)
		}
	}))
	defer server.Close()

	client, err := NewVideoChatClient(VideoChatClientOpts{APIBaseURL: server.URL + "/configured-api", HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video chat client: %v", err)
	}

	// When
	err = client.SendTextMessage(context.Background(), SendTextMessageInput{
		OwnerAccessToken: "owner-token",
		BotAccessToken:   "bot-token",
		Content:          "hello from bot",
	})

	// Then
	if err != nil {
		t.Fatalf("send VK Video chat message: %v", err)
	}
	for path, want := range map[string]int{
		"/configured-api/v1/current_user":      1,
		"/configured-api/v1/channel":           1,
		"/configured-api/v1/chat/message/send": 1,
	} {
		if got := requests[path]; got != want {
			t.Errorf("requests for %s = %d, want %d", path, got, want)
		}
	}
}

func TestVideoChatClientSkipsSendWhenChannelIsOffline(t *testing.T) {
	tests := []struct {
		name            string
		channelResponse string
	}{
		{name: "stream is absent", channelResponse: `{"data":{}}`},
		{name: "stream is null", channelResponse: `{"data":{"stream":null}}`},
		{name: "stream ID is blank", channelResponse: `{"data":{"stream":{"id":"  "}}}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			sendRequests := 0
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				switch request.URL.Path {
				case "/v1/current_user":
					if got := request.Header.Get("Authorization"); got != "Bearer owner-token" {
						t.Errorf("current user authorization = %q", got)
					}
					_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/offline"}}}`))
				case "/v1/channel":
					if got := request.Header.Get("Authorization"); got != "Bearer bot-token" {
						t.Errorf("channel authorization = %q", got)
					}
					_, _ = writer.Write([]byte(test.channelResponse))
				case "/v1/chat/message/send":
					sendRequests++
					t.Fatalf("send request made for offline channel")
				default:
					t.Fatalf("unexpected request: %s %s", request.Method, request.URL)
				}
			}))
			defer server.Close()

			client, err := NewVideoChatClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
			if err != nil {
				t.Fatalf("create VK Video chat client: %v", err)
			}

			// When
			err = client.SendTextMessage(context.Background(), SendTextMessageInput{
				OwnerAccessToken: "owner-token",
				BotAccessToken:   "bot-token",
				Content:          "hello from bot",
			})

			// Then
			if err != nil {
				t.Fatalf("send VK Video chat message: %v", err)
			}
			if sendRequests != 0 {
				t.Errorf("send requests = %d, want 0", sendRequests)
			}
		})
	}
}

func TestVideoChatClientDoesNotExposeTokenOrMessageWhenProviderRejectsSend(t *testing.T) {
	// Given
	postRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v1/current_user":
			_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/channel"}}}`))
		case "/v1/channel":
			_, _ = writer.Write([]byte(`{"data":{"stream":{"id":"stream-42"}}}`))
		case "/v1/chat/message/send":
			postRequests++
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(`{"error_description":"message secret"}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := NewVideoChatClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video chat client: %v", err)
	}

	// When
	err = client.SendTextMessage(context.Background(), SendTextMessageInput{
		OwnerAccessToken: "owner-token",
		BotAccessToken:   "bot-token",
		Content:          "message secret",
	})

	// Then
	if err == nil {
		t.Fatal("expected provider rejection")
	}
	if postRequests != 1 {
		t.Errorf("send requests = %d, want 1 without retry", postRequests)
	}
	if strings.Contains(err.Error(), "owner-token") || strings.Contains(err.Error(), "bot-token") || strings.Contains(err.Error(), "message secret") {
		t.Errorf("unsafe provider error = %q", err)
	}
}
