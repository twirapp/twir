package twitch

import (
	"net/http"
	"testing"

	cfg "github.com/twirapp/twir/libs/config"
)

type captureTransport struct {
	lastRequest *http.Request
}

func (c *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.lastRequest = req
	return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
}

func TestMockRoundTripper(t *testing.T) {
	capture := &captureTransport{}
	config := cfg.Config{
		TwitchMockApiUrl:  "http://mock-api:7777",
		TwitchMockAuthUrl: "http://mock-auth:7777",
	}
	rt := NewMockRoundTripper(capture, config)

	req, _ := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	_, _ = rt.RoundTrip(req)
	if capture.lastRequest.URL.Host != "mock-auth:7777" {
		t.Errorf("expected host mock-auth:7777, got %s", capture.lastRequest.URL.Host)
	}
	if capture.lastRequest.URL.Scheme != "http" {
		t.Errorf("expected scheme http, got %s", capture.lastRequest.URL.Scheme)
	}

	req, _ = http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	_, _ = rt.RoundTrip(req)
	if capture.lastRequest.URL.Host != "mock-api:7777" {
		t.Errorf("expected host mock-api:7777, got %s", capture.lastRequest.URL.Host)
	}

	req, _ = http.NewRequest("GET", "https://example.com/path", nil)
	_, _ = rt.RoundTrip(req)
	if capture.lastRequest.URL.Host != "example.com" {
		t.Errorf("expected host example.com, got %s", capture.lastRequest.URL.Host)
	}
}
