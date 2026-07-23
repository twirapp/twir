package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
)

func (c *Directives) HasChannelRolesDashboardPermission(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	permission *gqlmodel.ChannelRolePermissionEnum,
) (res interface{}, err error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, err
	}
	dashboardID, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	dashboardUUID, err := uuid.Parse(dashboardID)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel: %w", err)
	}
	if c.dashboardAccess == nil {
		return nil, fmt.Errorf("dashboard access service is not configured")
	}

	permissionValue := ""
	if permission != nil {
		permissionValue = permission.String()
	}
	hasAccess, err := c.dashboardAccess.CanAccess(ctx, dashboardaccess.Subject{
		ID:         user.ID,
		IsBotAdmin: user.IsBotAdmin,
	}, dashboardUUID, permissionValue)
	if err != nil {
		return nil, fmt.Errorf("cannot check dashboard access: %w", err)
	}
	if hasAccess {
		return next(ctx)
	}

	return nil, fmt.Errorf("user has no permission to access this resource")
}
