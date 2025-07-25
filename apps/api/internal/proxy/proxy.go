package proxy

import (
	"io"
	"net/http"

	"github.com/twirapp/twir/apps/api/internal/handlers"
)

type Proxy struct {
}

func New() *Proxy {
	return &Proxy{}
}

var _ handlers.IHandler = (*Proxy)(nil)

func (c *Proxy) Pattern() string {
	return "/proxy"
}

func (c *Proxy) Handler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			proxyUrl := r.URL.Query().Get("url")
			if proxyUrl == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("url is required"))
				return
			}

			proxyReq, proxyErr := http.Get(proxyUrl)
			if proxyErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(proxyErr.Error()))
				return
			}

			defer proxyReq.Body.Close()
			body, bodyErr := io.ReadAll(proxyReq.Body)
			if bodyErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(bodyErr.Error()))
				return
			}

			w.Header().Set("Content-Type", proxyReq.Header.Get("Content-Type"))
			w.Write(body)
		},
	)
}
