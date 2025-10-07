package middlewares

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (c *Middlewares) IsAdmin(ctx huma.Context, next func(huma.Context)) {
	user, err := c.auth.GetAuthenticatedUserModel(ctx.Context())
	if err != nil {
		huma.WriteErr(
			c.huma,
			ctx,
			http.StatusInternalServerError,
			"Cannot get authenticated user",
			err,
		)

		return
	}

	if !user.IsBotAdmin {
		huma.WriteErr(
			c.huma,
			ctx,
			http.StatusForbidden,
			"Not a bot admin",
			nil,
		)

		return
	}

	next(ctx)
}
