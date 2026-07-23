package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"gorm.io/gorm"
)

func (r *authenticatedUserResolver) getDashboardPlatform(
	ctx context.Context,
	channel model.Channels,
	userID string,
) string {
	parsedChannelID, err := uuid.Parse(channel.ID)
	if err != nil {
		return channel.Platform()
	}

	normalizedChannel, err := r.deps.ChannelService.GetChannelByID(ctx, parsedChannelID)
	if err != nil {
		return channel.Platform()
	}

	for _, binding := range normalizedChannel.Bindings {
		if binding.UserID.String() == userID {
			return binding.Platform.String()
		}
	}
	for _, binding := range normalizedChannel.Bindings {
		return binding.Platform.String()
	}

	return channel.Platform()
}

func ownedDashboardsQuery(db *gorm.DB, ctx context.Context, userID string) *gorm.DB {
	return db.WithContext(ctx).Where(
		`EXISTS (
			SELECT 1
			FROM channel_platforms AS cp_owner
			WHERE cp_owner.channel_id = channels.id
				AND cp_owner.user_id = ?::uuid
		) OR (
			NOT EXISTS (
				SELECT 1
				FROM channel_platforms AS cp_existing
				WHERE cp_existing.channel_id = channels.id
			) AND (
				channels.twitch_user_id = ?::uuid
				OR channels.kick_user_id = ?::uuid
			)
		)`,
		userID,
		userID,
		userID,
	)
}

func (r *authenticatedUserResolver) getAuthenticatedUserChannel(ctx context.Context) (channelsmodel.Channel, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("get selected dashboard: %w", err)
	}

	if dashboardID == "" {
		return channelsmodel.Nil, nil
	}

	parsedDashboardID, err := uuid.Parse(dashboardID)
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("parse selected dashboard id: %w", err)
	}

	channel, err := r.deps.ChannelService.GetChannelByID(ctx, parsedDashboardID)
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("get selected dashboard channel: %w", err)
	}

	return channel, nil
}

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
				Platform: r.getDashboardPlatform(ctx, channel, obj.ID),
				Flags: []gqlmodel.ChannelRolePermissionEnum{
					gqlmodel.ChannelRolePermissionEnumCanAccessDashboard,
				},
				PlanID: channel.PlanID,
			}

			dashboardsEntities[channel.ID] = dashboard
		}
	} else {
		var ownChannels []model.Channels
		if err := ownedDashboardsQuery(r.deps.Gorm, ctx, obj.ID).
			Find(&ownChannels).Error; err != nil {
			return nil, err
		}

		for _, channel := range ownChannels {
			dashboard := gqlmodel.Dashboard{
				ID:       channel.ID,
				Platform: r.getDashboardPlatform(ctx, channel, obj.ID),
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
				`user_id = ?`,
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
			platform := existing.Platform
			if platform == "" {
				platform = r.getDashboardPlatform(ctx, *role.Role.Channel, obj.ID)
			}

			dashboard := gqlmodel.Dashboard{
				ID:          role.Role.Channel.ID,
				Platform:    platform,
				Flags:       append(existing.Flags, flags...),
				APIKey:      existing.APIKey,
				PlanID:      role.Role.Channel.PlanID,
				KickProfile: existing.KickProfile,
			}

			dashboardsEntities[role.Role.Channel.ID] = dashboard
		}
	}

	var usersStats []model.UsersStats
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`user_id = ?`, obj.ID).
		Find(&usersStats).Error; err != nil {
		return nil, err
	}

	for _, stat := range usersStats {
		var channelRoles []model.ChannelRole
		if err := r.deps.Gorm.WithContext(ctx).
			Where(`"channelId" = ?::uuid`, stat.ChannelID).
			Preload("Channel").
			Find(&channelRoles).
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

		if role.ID != "" && len(flags) > 0 && role.Channel != nil {
			existing := dashboardsEntities[role.ChannelID]
			platform := existing.Platform
			if platform == "" {
				platform = r.getDashboardPlatform(ctx, *role.Channel, obj.ID)
			}

			dashboard := gqlmodel.Dashboard{
				ID:          role.ChannelID,
				Platform:    platform,
				Flags:       append(existing.Flags, flags...),
				APIKey:      existing.APIKey,
				PlanID:      existing.PlanID,
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
