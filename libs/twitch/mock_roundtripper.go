package twitch

import (
	"net/http"
	"net/url"
	"strings"

	cfg "github.com/twirapp/twir/libs/config"
)

// MockRoundTripper rewrites Twitch requests to the mock server.
type MockRoundTripper struct {
	Base    http.RoundTripper
	ApiUrl  string
	AuthUrl string
}

func (t *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())

	if strings.Contains(req.URL.Host, "id.twitch.tv") {
		if u, err := url.Parse(t.AuthUrl); err == nil {
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
		}
	} else if strings.Contains(req.URL.Host, "api.twitch.tv") {
		if u, err := url.Parse(t.ApiUrl); err == nil {
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
		}
	}

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(req)
}

// NewMockRoundTripper wraps the given base transport with mock routing.
func NewMockRoundTripper(base http.RoundTripper, config cfg.Config) http.RoundTripper {
	return &MockRoundTripper{
		Base:    base,
		ApiUrl:  config.TwitchMockApiUrl,
		AuthUrl: config.TwitchMockAuthUrl,
	}
}
