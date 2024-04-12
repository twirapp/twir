package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func (c *Directives) IsAuthenticated(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	_, err := c.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	return next(ctx)
}
