package middlewares

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/enums/dashboard_permissions"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
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

		dashboardUUID, err := uuid.Parse(dashboardId)
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

		if c.dashboardAccess == nil {
			huma.WriteErr(c.huma, hc, http.StatusInternalServerError, "Cannot check dashboard access", nil)
			return
		}
		hasAccess, err := c.dashboardAccess.CanAccess(ctx, dashboardaccess.Subject{
			ID:         user.ID,
			IsBotAdmin: user.IsBotAdmin,
		}, dashboardUUID, permission.String())
		if err != nil {
			huma.WriteErr(c.huma, hc, http.StatusInternalServerError, "Cannot check dashboard access", err)
			return
		}
		if hasAccess {
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
