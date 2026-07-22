package middlewares

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/enums/dashboard_permissions"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func (c *Middlewares) HasChannelRolesDashboardPermission(permission dashboard_permissions.ChannelRolePermissionEnum) func(
	hc huma.Context,
	next func(huma.Context),
) {
	return func(hc huma.Context, next func(huma.Context)) {
		ctx := hc.Context()

		user, err := c.auth.GetAuthenticatedUserModel(ctx)
		if err != nil {
			huma.WriteErr(
				c.huma,
				hc,
				http.StatusInternalServerError,
				"Cannot get authenticated user",
				err,
			)

			return
		}
		dashboardId, err := c.auth.GetSelectedDashboard(ctx)
		if err != nil {
			huma.WriteErr(
				c.huma,
				hc,
				http.StatusInternalServerError,
				"Cannot get selected dashboard",
				err,
			)

			return
		}

		isOwner, err := c.isSelectedDashboardOwner(ctx, dashboardId, user.ID)
		if err != nil {
			huma.WriteErr(
				c.huma,
				hc,
				http.StatusInternalServerError,
				"Cannot get channel",
				err,
			)
			return
		}

		if isOwner || user.IsBotAdmin {
			next(hc)
			return
		}

		var channelRoles []model.ChannelRole
		if err := c.gorm.
			WithContext(ctx).
			Where(`"channelId" = ?::uuid`, dashboardId).
			Preload("Users", `user_id = ?`, user.ID).
			Find(&channelRoles).
			Error; err != nil {
			huma.WriteErr(
				c.huma,
				hc,
				http.StatusInternalServerError,
				"Cannot get channel roles",
				err,
			)

			return
		}

		var userStat model.UsersStats
		if err := c.gorm.
			WithContext(ctx).
			Where(`user_id = ? AND channel_id = ?::uuid`, user.ID, dashboardId).
			First(&userStat).
			Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			huma.WriteErr(
				c.huma,
				hc,
				http.StatusInternalServerError,
				"Cannot get user stats",
				err,
			)

			return
		}

		if hasChannelRolesDashboardAccess(channelRoles, user.ID, userStat, &permission) {
			next(hc)
			return
		}

		huma.WriteErr(
			c.huma,
			hc,
			http.StatusForbidden,
			"user does not have access to this permission",
			nil,
		)

		return
	}
}
