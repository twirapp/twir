package seventv

import (
	"net/http"
)

type headerTransport struct {
	base  http.RoundTripper
	token string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)

	return t.base.RoundTrip(req)
}
