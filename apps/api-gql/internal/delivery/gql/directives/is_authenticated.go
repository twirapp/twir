package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (c *Directives) IsAuthenticated(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	_, apiKeyErr := c.sessions.GetAuthenticatedUserByApiKey(ctx)
	_, sessionErr := c.sessions.GetAuthenticatedUserModel(ctx)

	if apiKeyErr != nil && sessionErr != nil {
		return nil, fmt.Errorf("not authenticated")
	}

	return next(ctx)
}
