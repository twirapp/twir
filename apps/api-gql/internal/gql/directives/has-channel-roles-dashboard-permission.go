package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

func (c *Directives) HasChannelRolesDashboardPermission(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	permission *gqlmodel.ChannelRolePermissionEnum,
) (res interface{}, err error) {
	user, err := c.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	dashboardId, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	if user.ID == dashboardId || user.IsBotAdmin {
		return next(ctx)
	}

	var userRoles []model.ChannelRoleUser
	if err := c.gorm.
		WithContext(ctx).
		Where(`channels_roles_users."userId"`, user.ID).
		Joins("Role", `"channelId = ?"`, dashboardId).
		Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("cannot get user userRoles, probably have no access: %w", err)
	}

	for _, role := range userRoles {
		for _, perm := range role.Role.Permissions {
			if perm == gqlmodel.ChannelRolePermissionEnumCanAccessDashboard.String() {
				return next(ctx)
			}

			if permission.String() == perm {
				return next(ctx)
			}
		}
	}

	return nil, fmt.Errorf("user has no permission to access this resource")
}
