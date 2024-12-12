package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	model "github.com/satont/twir/libs/gomodels"
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

	dashboardId, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	if user.ID == dashboardId || user.IsBotAdmin {
		return next(ctx)
	}

	role := &model.ChannelRoleUser{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`channels_roles_users."userId"`, user.ID).
		Joins("Role", `"channelId = ?"`, dashboardId).
		First(&role).Error; err != nil {
		return nil, fmt.Errorf("cannot get user roles, probably have no access: %w", err)
	}

	return next(ctx)
}
