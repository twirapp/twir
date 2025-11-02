package humahelpers

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
)

func GetClientIpFromCtx(ctx context.Context) (string, error) {
	gCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return "", err
	}

	return gCtx.ClientIP(), nil
}

func GetClientUserAgentFromCtx(ctx context.Context) (string, error) {
	gCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return "", err
	}
	
	return gCtx.GetHeader("user-agent"), nil
}
