package helpers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/twirapp/twir/apps/api/internal/wrappers"
)

func GetHeadersFromCtx(ctx context.Context) (http.Header, error) {
	headers, ok := ctx.Value(wrappers.ContextHeadersKey).(http.Header)
	if !ok {
		return nil, fmt.Errorf("cannot get headers from request")
	}

	return headers, nil
}
