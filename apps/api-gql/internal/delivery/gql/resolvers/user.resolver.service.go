package resolvers

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (r *authenticatedUserResolver) getAvailableDashboards(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) ([]gqlmodel.Dashboard, error) {
	dashboardsEntities := make(map[string]gqlmodel.Dashboard)

	kickAccounts, _ := r.deps.UserPlatformAccountsRepository.GetAllByPlatform(ctx, platformentity.PlatformKick)
	kickProfileByUserID := make(map[string]gqlmodel.KickProfile)
	for _, acc := range kickAccounts {
		var pic *string
		if acc.PlatformAvatar != "" {
			pic = &acc.PlatformAvatar
		}
		kickProfileByUserID[acc.UserID.String()] = gqlmodel.KickProfile{
			ID:             acc.PlatformUserID,
			Slug:           acc.PlatformLogin,
			DisplayName:    acc.PlatformDisplayName,
			ProfilePicture: pic,
			IsLive:         false,
			FollowersCount: 0,
		}
	}

	if obj.IsBotAdmin {
		var channels []model.Channels
		if err := r.deps.Gorm.WithContext(ctx).Preload("User").Find(&channels).Error; err != nil {
			return nil, err
		}

		for _, channel := range channels {
			dashboard := gqlmodel.Dashboard{
				ID:       channel.ID,
				Platform: channel.Platform,
				Flags: []gqlmodel.ChannelRolePermissionEnum{
					gqlmodel.ChannelRolePermissionEnumCanAccessDashboard,
				},
				PlanID: channel.PlanID,
			}

			if channel.User != nil {
				dashboard.APIKey = channel.User.ApiKey
			}

			if channel.Platform == string(platformentity.PlatformKick) && channel.User != nil {
				if kp, ok := kickProfileByUserID[channel.User.ID]; ok {
					dashboard.KickProfile = &kp
				}
			}

			dashboardsEntities[channel.ID] = dashboard
		}
	} else {
		var ownChannels []model.Channels
		if err := r.deps.Gorm.WithContext(ctx).Where("user_id = ?", obj.ID).Find(&ownChannels).Error; err != nil {
			return nil, err
		}

		for _, channel := range ownChannels {
			dashboard := gqlmodel.Dashboard{
				ID:       channel.ID,
				Platform: channel.Platform,
				Flags:    []gqlmodel.ChannelRolePermissionEnum{gqlmodel.ChannelRolePermissionEnumCanAccessDashboard},
				APIKey:   obj.APIKey,
				PlanID:   channel.PlanID,
			}

			if channel.Platform == string(platformentity.PlatformKick) {
				if kp, ok := kickProfileByUserID[obj.ID]; ok {
					dashboard.KickProfile = &kp
				}
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
			Preload("Role.Channel.User").
			Find(&roles).
			Error; err != nil {
			return nil, err
		}

		for _, role := range roles {
			if role.Role == nil || role.Role.Channel == nil || role.Role.Channel.User == nil || len(role.Role.Permissions) == 0 {
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
				APIKey:      role.Role.Channel.User.ApiKey,
				KickProfile: existing.KickProfile,
			}

			if role.Role.Channel.User != nil {
				dashboard.APIKey = role.Role.Channel.User.ApiKey
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

			if role.Channel != nil && role.Channel.User != nil {
				dashboard.APIKey = role.Channel.User.ApiKey
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
