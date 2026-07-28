package vkvideo

import (
	"net/url"
	"testing"

	cfg "github.com/twirapp/twir/libs/config"
)

func TestBotSetupAuthURLRequestsChatMessageSendScope(t *testing.T) {
	provider := NewBotSetupProvider(BotSetupProviderOpts{
		Config: cfg.Config{
			SiteBaseUrl:          "https://twir.example.test",
			VKVideoClientID:      "client-id",
			VKVideoClientSecret:  "client-secret",
			VKVideoAPIBaseURL:    "https://api.example.test",
			VKVideoAuthBaseURL:   "https://auth.example.test",
			VKVideoDevAPIBaseURL: "https://devapi.example.test",
		},
	})

	rawURL, err := provider.GetBotSetupAuthURL("state-value")
	if err != nil {
		t.Fatalf("get bot setup auth URL: %v", err)
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse bot setup auth URL: %v", err)
	}
	if got := parsed.Query().Get("scope"); got != vkVideoBotChatSendScope {
		t.Fatalf("scope = %q, want %q", got, vkVideoBotChatSendScope)
	}
}
