package spotify

import (
	"io"
	"net/http"
	"strings"
)

func spotifyTestResponse(req *http.Request, statusCode int, body string, headers http.Header) *http.Response {
	if headers == nil {
		headers = make(http.Header)
	}

	return &http.Response{
		StatusCode: statusCode,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}
