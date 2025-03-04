package middlewares

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (c *Middlewares) IsAuthenticated(ctx huma.Context, next func(huma.Context)) {
	_, apiKeyErr := c.auth.GetAuthenticatedUserByApiKey(ctx.Context())
	_, sessionErr := c.auth.GetAuthenticatedUser(ctx.Context())

	if apiKeyErr != nil && sessionErr != nil {
		huma.WriteErr(
			c.huma,
			ctx,
			http.StatusUnauthorized,
			"Not authorized",
			apiKeyErr,
			sessionErr,
		)

		return
	}

	next(ctx)
}
