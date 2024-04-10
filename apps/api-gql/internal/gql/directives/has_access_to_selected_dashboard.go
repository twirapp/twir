package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (c *Directives) HasAccessToSelectedDashboard(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	user, err := c.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	dashboardID, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: perform access check here
	if dashboardID != user.ID {
		return nil, fmt.Errorf("user does not have access to dashboard")
	}

	return next(ctx)
}
