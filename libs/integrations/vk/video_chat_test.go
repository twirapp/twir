package vk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVideoChatClientResolvesStreamForExactOwnerID(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/v1/channels/active" {
			t.Fatalf("active streams request = %s %s", request.Method, request.URL)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer bot-token" {
			t.Errorf("active streams authorization = %q", got)
		}
		_, _ = writer.Write([]byte(`{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/near-match"},"owner":{"id":421},"stream":{"id":"stream-421"}},{"channel":{"url":"https://live.vkvideo.ru/exact-match"},"owner":{"id":"42"},"stream":{"id":"stream-42"}}]}}`))
	}))
	defer server.Close()

	client, err := NewVideoChatClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video chat client: %v", err)
	}

	// When
	stream, err := client.ResolveActiveStream(context.Background(), "bot-token", "42")

	// Then
	if err != nil {
		t.Fatalf("resolve active stream: %v", err)
	}
	if stream.ChannelURL != "https://live.vkvideo.ru/exact-match" || stream.ID != "stream-42" {
		t.Errorf("active stream = %#v", stream)
	}
}

func TestVideoChatClientSendsTextMessageToResolvedStream(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer bot-token" {
			t.Errorf("authorization = %q", got)
		}

		switch request.URL.Path {
		case "/v1/channels/active":
			_, _ = writer.Write([]byte(`{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/channel"},"owner":{"id":42},"stream":{"id":"stream-42"}}]}}`))
		case "/v1/chat/message/send":
			if request.Method != http.MethodPost {
				t.Errorf("send request method = %s", request.Method)
			}
			if got := request.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("send content type = %q", got)
			}
			if got := request.URL.Query().Get("channel_url"); got != "https://live.vkvideo.ru/channel" {
				t.Errorf("channel_url = %q", got)
			}
			if got := request.URL.Query().Get("stream_id"); got != "stream-42" {
				t.Errorf("stream_id = %q", got)
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

	client, err := NewVideoChatClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video chat client: %v", err)
	}

	// When
	err = client.SendTextMessage(context.Background(), SendTextMessageInput{
		AccessToken: "bot-token",
		OwnerID:     "42",
		Content:     "hello from bot",
	})

	// Then
	if err != nil {
		t.Fatalf("send VK Video chat message: %v", err)
	}
}

func TestVideoChatClientDoesNotExposeTokenOrMessageWhenProviderRejectsSend(t *testing.T) {
	// Given
	postRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v1/channels/active":
			_, _ = writer.Write([]byte(`{"data":{"channels":[{"channel":{"url":"https://live.vkvideo.ru/channel"},"owner":{"id":42},"stream":{"id":"stream-42"}}]}}`))
		case "/v1/chat/message/send":
			postRequests++
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(`{"error_description":"message secret"}`))
		}
	}))
	defer server.Close()

	client, err := NewVideoChatClient(VideoChatClientOpts{APIBaseURL: server.URL, HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("create VK Video chat client: %v", err)
	}

	// When
	err = client.SendTextMessage(context.Background(), SendTextMessageInput{
		AccessToken: "bot-token",
		OwnerID:     "42",
		Content:     "message secret",
	})

	// Then
	if err == nil {
		t.Fatal("expected provider rejection")
	}
	if postRequests != 1 {
		t.Errorf("send requests = %d, want 1 without retry", postRequests)
	}
	if strings.Contains(err.Error(), "bot-token") || strings.Contains(err.Error(), "message secret") {
		t.Errorf("unsafe provider error = %q", err)
	}
}
