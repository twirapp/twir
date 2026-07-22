package middlewares

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/enums/dashboard_permissions"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type selectedDashboardChannelGetter interface {
	GetChannelByID(ctx context.Context, channelID uuid.UUID) (channelsmodel.Channel, error)
}

func (c *Middlewares) isSelectedDashboardOwner(
	ctx context.Context,
	dashboardID string,
	userID string,
) (bool, error) {
	channelID, err := uuid.Parse(dashboardID)
	if err != nil {
		return false, err
	}

	channel, err := c.channelGetter.GetChannelByID(ctx, channelID)
	if err != nil {
		return false, err
	}

	for _, binding := range channel.Bindings {
		if binding.UserID.String() == userID {
			return true, nil
		}
	}

	return false, nil
}

func hasChannelRolesDashboardAccess(
	channelRoles []model.ChannelRole,
	userID string,
	userStat model.UsersStats,
	permission *dashboard_permissions.ChannelRolePermissionEnum,
) bool {
	roleToStats := map[model.ChannelRoleEnum]bool{
		model.ChannelRoleTypeModerator:  userStat.IsMod,
		model.ChannelRoleTypeVip:        userStat.IsVip,
		model.ChannelRoleTypeSubscriber: userStat.IsSubscriber,
	}

	for i, role := range channelRoles {
		if roleToStats[role.Type] {
			channelRoles[i].Users = append(
				role.Users,
				&model.ChannelRoleUser{
					ID:     "", // not needed
					UserID: userID,
					RoleID: role.ID,
				},
			)
		}
	}

	for _, role := range channelRoles {
		if len(role.Users) == 0 || len(role.Permissions) == 0 {
			continue
		}

		hasRoleUser := false
		for _, roleUser := range role.Users {
			if roleUser.UserID == userID {
				hasRoleUser = true
				break
			}
		}
		if !hasRoleUser {
			continue
		}

		if permission == nil {
			return true
		}

		for _, rolePermission := range role.Permissions {
			if rolePermission == dashboard_permissions.ChannelRolePermissionEnumCanAccessDashboard.String() ||
				rolePermission == permission.String() {
				return true
			}
		}
	}

	return false
}
