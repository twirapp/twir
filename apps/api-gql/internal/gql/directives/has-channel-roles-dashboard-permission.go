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

	var channelRoles []model.ChannelRole
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Preload("Users", `"userId" = ?`, user.ID).
		Find(&channelRoles).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get channel roles: %w", err)
	}

	for _, role := range channelRoles {
		// we do not check does role.Users contains request author user
		// because we are doing preload by user id
		if len(role.Users) == 0 || len(role.Permissions) == 0 {
			continue
		}

		for _, perm := range role.Permissions {
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
