package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/twirapp/twir/libs/config"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
	"go.uber.org/fx"
)

func TestAppHasCompleteDependencyGraph(t *testing.T) {
	if err := fx.ValidateApp(App); err != nil {
		t.Fatalf("validate bots dependency graph: %v", err)
	}
}

func TestNewVKVideoChatClientUsesConfiguredDevAPIBaseURL(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v1/current_user":
			_, _ = writer.Write([]byte(`{"data":{"channel":{"url":"https://live.vkvideo.ru/channel"}}}`))
		case "/v1/channel":
			_, _ = writer.Write([]byte(`{"data":{"stream":{"id":"stream-42"}}}`))
		case "/v1/chat/message/send":
			writer.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := newVKVideoChatClient(cfg.Config{VKVideoDevAPIBaseURL: server.URL})
	if err != nil {
		t.Fatalf("create VK Video chat client: %v", err)
	}

	// When
	err = client.SendTextMessage(context.Background(), vkintegrations.SendTextMessageInput{
		OwnerAccessToken: "owner-token",
		BotAccessToken:   "bot-token",
		Content:          "hello from bot",
	})

	// Then
	if err != nil {
		t.Fatalf("send VK Video chat message through configured base URL: %v", err)
	}
}
