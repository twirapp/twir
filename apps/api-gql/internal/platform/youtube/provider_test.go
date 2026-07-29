package youtube

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	cfg "github.com/twirapp/twir/libs/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestProviderBuildsOfflineConsentPKCEURLs(t *testing.T) {
	// Given
	provider := New(Opts{Config: cfg.Config{SiteBaseUrl: "https://twir.example.test", YouTubeClientID: "client-id", YouTubeClientSecret: "client-secret"}})

	// When
	userURL, err := url.Parse(provider.GetAuthURL("user-state", "user-challenge"))
	if err != nil {
		t.Fatalf("parse user authorization URL: %v", err)
	}
	botURL, err := url.Parse(provider.GetBotSetupAuthURL("bot-state", "bot-challenge"))
	if err != nil {
		t.Fatalf("parse bot authorization URL: %v", err)
	}

	// Then
	for _, testCase := range []struct {
		name      string
		url       *url.URL
		state     string
		challenge string
		redirect  string
	}{
		{name: "user", url: userURL, state: "user-state", challenge: "user-challenge", redirect: "https://twir.example.test/login/youtube"},
		{name: "bot", url: botURL, state: "bot-state", challenge: "bot-challenge", redirect: "https://twir.example.test/api/auth/youtube/bot-callback"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			query := testCase.url.Query()
			if got := query.Get("state"); got != testCase.state {
				t.Fatalf("state = %q, want %q", got, testCase.state)
			}
			if got := query.Get("redirect_uri"); got != testCase.redirect {
				t.Fatalf("redirect_uri = %q, want %q", got, testCase.redirect)
			}
			if got := query.Get("access_type"); got != "offline" {
				t.Fatalf("access_type = %q, want offline", got)
			}
			if got := query.Get("prompt"); got != "consent" {
				t.Fatalf("prompt = %q, want consent", got)
			}
			if got := query.Get("code_challenge"); got != testCase.challenge || query.Get("code_challenge_method") != "S256" {
				t.Fatalf("PKCE query = %q / %q", got, query.Get("code_challenge_method"))
			}
		})
	}
}

func TestProviderGetUserMapsOwnChannel(t *testing.T) {
	// Given
	provider := &Provider{httpClient: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.String() != youtubeChannelsURL {
			t.Fatalf("channels URL = %s", request.URL)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer access-token" {
			t.Fatalf("authorization = %q", authorization)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"items":[{"id":"channel-id","snippet":{"title":"Channel Title","customUrl":"@channel","thumbnails":{"high":{"url":"https://example.test/high.jpg"}}}}]}`)), Header: make(http.Header)}, nil
	})}}

	// When
	user, err := provider.GetUser(context.Background(), "access-token")

	// Then
	if err != nil {
		t.Fatalf("get YouTube user: %v", err)
	}
	if user.ID != "channel-id" || user.Login != "@channel" || user.DisplayName != "Channel Title" || user.Avatar != "https://example.test/high.jpg" {
		t.Fatalf("user = %#v", user)
	}
}
