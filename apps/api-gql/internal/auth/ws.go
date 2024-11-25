package auth

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler/transport"
)

type WsApiKeyContextKey struct{}

func (s *Auth) getWsAuthenticatedApiKey(ctx context.Context) (string, error) {
	key := ctx.Value(WsApiKeyContextKey{})
	if key == nil {
		return "", fmt.Errorf("api key not presented in context")
	}

	apiKey, ok := key.(string)
	if !ok {
		return "", fmt.Errorf("api key not a string")
	}

	return apiKey, nil
}

func WsGqlInitFunc(ctx context.Context, initPayload transport.InitPayload) (
	context.Context,
	*transport.InitPayload,
	error,
) {
	key := initPayload.GetString("api-key")
	if key == "" {
		return ctx, &initPayload, nil
	}

	ctx = context.WithValue(ctx, WsApiKeyContextKey{}, key)
	return ctx, &initPayload, nil
}
