package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func (c *Directives) HasAccessToSelectedDashboard(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	_, err := c.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	_, err = c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: perform access check here
	// if dashboardID != user.ID {
	// 	return nil, fmt.Errorf("user does not have access to dashboard")
	// }

	return next(ctx)
}
