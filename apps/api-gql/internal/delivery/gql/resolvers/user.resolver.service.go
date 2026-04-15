package resolvers

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (r *authenticatedUserResolver) getAvailableDashboards(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) ([]gqlmodel.Dashboard, error) {
	dashboardsEntities := make(map[string]gqlmodel.Dashboard)

	if obj.IsBotAdmin {
		var channels []model.Channels
		if err := r.deps.Gorm.WithContext(ctx).Find(&channels).Error; err != nil {
			return nil, err
		}

		for _, channel := range channels {
			dashboard := gqlmodel.Dashboard{
				ID:       channel.ID,
				Platform: channel.Platform(),
				Flags: []gqlmodel.ChannelRolePermissionEnum{
					gqlmodel.ChannelRolePermissionEnumCanAccessDashboard,
				},
				PlanID: channel.PlanID,
			}

			dashboardsEntities[channel.ID] = dashboard
		}
	} else {
		var ownChannels []model.Channels
		if err := r.deps.Gorm.WithContext(ctx).
			Where("twitch_user_id = (SELECT id FROM users WHERE id = ? LIMIT 1)", obj.ID).
			Or("kick_user_id = (SELECT id FROM users WHERE id = ? LIMIT 1)", obj.ID).
			Find(&ownChannels).Error; err != nil {
			return nil, err
		}

		for _, channel := range ownChannels {
			dashboard := gqlmodel.Dashboard{
				ID:       channel.ID,
				Platform: channel.Platform(),
				Flags:    []gqlmodel.ChannelRolePermissionEnum{gqlmodel.ChannelRolePermissionEnumCanAccessDashboard},
				APIKey:   obj.APIKey,
				PlanID:   channel.PlanID,
			}

			dashboardsEntities[channel.ID] = dashboard
		}

		var roles []model.ChannelRoleUser
		if err := r.deps.Gorm.
			WithContext(ctx).
			Where(
				`"userId" = ?`,
				obj.ID,
			).
			Preload("Role").
			Preload("Role.Channel").
			Find(&roles).
			Error; err != nil {
			return nil, err
		}

		for _, role := range roles {
			if role.Role == nil || role.Role.Channel == nil || len(role.Role.Permissions) == 0 {
				continue
			}

			var flags []gqlmodel.ChannelRolePermissionEnum
			for _, flag := range role.Role.Permissions {
				flags = append(flags, gqlmodel.ChannelRolePermissionEnum(flag))
			}

			existing := dashboardsEntities[role.Role.Channel.ID]
			dashboard := gqlmodel.Dashboard{
				ID:          role.Role.Channel.ID,
				Platform:    existing.Platform,
				Flags:       append(existing.Flags, flags...),
				PlanID:      role.Role.Channel.PlanID,
				KickProfile: existing.KickProfile,
			}

			dashboardsEntities[role.Role.Channel.ID] = dashboard
		}
	}

	var usersStats []model.UsersStats
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"userId" = ?`, obj.ID).
		Find(&usersStats).Error; err != nil {
		return nil, err
	}

	for _, stat := range usersStats {
		var channelRoles []model.ChannelRole
		if err := r.deps.Gorm.WithContext(ctx).Where(
			`"channelId" = ?`,
			stat.ChannelID,
		).Find(&channelRoles).
			Error; err != nil {
			return nil, err
		}

		var role model.ChannelRole

		if stat.IsMod {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeModerator
				},
			)
		} else if stat.IsVip {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeVip
				},
			)
		} else if stat.IsSubscriber {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeSubscriber
				},
			)
		}

		var flags []gqlmodel.ChannelRolePermissionEnum
		for _, flag := range role.Permissions {
			flags = append(flags, gqlmodel.ChannelRolePermissionEnum(flag))
		}

		if role.ID != "" && len(flags) > 0 {
			existing := dashboardsEntities[role.ChannelID]
			dashboard := gqlmodel.Dashboard{
				ID:          role.ChannelID,
				Platform:    existing.Platform,
				Flags:       append(existing.Flags, flags...),
				KickProfile: existing.KickProfile,
			}

			dashboardsEntities[role.ChannelID] = dashboard
		}
	}

	return lo.MapToSlice(
		dashboardsEntities,
		func(_ string, value gqlmodel.Dashboard) gqlmodel.Dashboard {
			return value
		},
	), nil
}
