package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (c *Directives) IsAdmin(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, err
	}

	if !user.IsBotAdmin {
		return nil, fmt.Errorf("not a bot admin")
	}

	return next(ctx)
}
