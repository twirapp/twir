package app

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
)

func NewAuthMiddleware(api huma.API, sessions *auth.Auth) func(
	ctx huma.Context,
	next func(huma.Context),
) {
	return func(ctx huma.Context, next func(huma.Context)) {
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			if _, ok := opScheme["api-key"]; ok {
				isAuthorizationRequired = true
				break
			}
		}

		if !isAuthorizationRequired {
			next(ctx)
			return
		}

		user, err := sessions.GetAuthenticatedUser(ctx.Context())
		if err != nil {
			huma.WriteErr(
				api,
				ctx,
				401,
				"not authenticated",
				err,
			)

			return
		}

		if user == nil {
			huma.WriteErr(
				api,
				ctx,
				401,
				"not authenticated",
				nil,
			)

			return
		}

		next(ctx)
	}
}
