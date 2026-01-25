package humahelpers

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
)

func GetHostFromCtx(ctx context.Context) (string, error) {
	gCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return "", err
	}

	host := strings.TrimSpace(gCtx.Request.Host)
	if host == "" {
		host = strings.TrimSpace(gCtx.GetHeader("Host"))
	}
	if host == "" {
		return "", errors.New("host header not found")
	}

	host = strings.ToLower(host)
	if strings.Contains(host, ":") {
		if splitHost, _, err := net.SplitHostPort(host); err == nil {
			host = splitHost
		}
	}

	return host, nil
}
