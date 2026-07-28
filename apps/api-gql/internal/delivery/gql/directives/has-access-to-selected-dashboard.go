package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
)

func (c *Directives) HasAccessToSelectedDashboard(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, err
	}

	dashboardId, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	dashboardUUID, err := uuid.Parse(dashboardId)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel: %w", err)
	}
	if c.dashboardAccess == nil {
		return nil, fmt.Errorf("dashboard access service is not configured")
	}
	hasAccess, err := c.dashboardAccess.CanAccess(ctx, dashboardaccess.Subject{
		ID:         user.ID,
		IsBotAdmin: user.IsBotAdmin,
	}, dashboardUUID, "")
	if err != nil {
		return nil, fmt.Errorf("cannot check dashboard access: %w", err)
	}
	if hasAccess {
		return next(ctx)
	}

	return nil, fmt.Errorf("user does not have access to selected dashboard")
}
