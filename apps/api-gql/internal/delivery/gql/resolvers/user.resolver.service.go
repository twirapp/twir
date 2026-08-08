package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	"gorm.io/gorm"
)

type dashboardPlatformReader interface {
	GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error)
}

func (r *authenticatedUserResolver) getDashboardPlatform(
	ctx context.Context,
	channelID string,
	userID string,
) string {
	return resolveDashboardPlatform(ctx, r.deps.ChannelService, channelID, userID)
}

func resolveDashboardPlatform(
	ctx context.Context,
	reader dashboardPlatformReader,
	channelID string,
	userID string,
) string {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return ""
	}

	normalizedChannel, err := reader.GetChannelByID(ctx, parsedChannelID)
	if err != nil {
		return ""
	}

	for _, binding := range normalizedChannel.Bindings {
		if binding.UserID.String() == userID {
			return binding.Platform.String()
		}
	}
	if len(normalizedChannel.Bindings) == 0 {
		return ""
	}

	return normalizedChannel.Bindings[0].Platform.String()
}

func ownedDashboardsQuery(db *gorm.DB, ctx context.Context, userID string) *gorm.DB {
	return db.WithContext(ctx).Where(
		`EXISTS (
			SELECT 1
			FROM channel_platforms AS cp_owner
			WHERE cp_owner.channel_id = channels.id
				AND cp_owner.user_id = ?::uuid
		)`,
		userID,
	)
}

func (r *authenticatedUserResolver) getAuthenticatedUserChannel(ctx context.Context) (channelentity.Channel, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return channelentity.Nil, fmt.Errorf("get selected dashboard: %w", err)
	}

	if dashboardID == "" {
		return channelentity.Nil, nil
	}

	parsedDashboardID, err := uuid.Parse(dashboardID)
	if err != nil {
		// Legacy sessions may hold a non-UUID dashboard id (e.g. a twitch channel id),
		// in that case there is no channel to resolve.
		return channelentity.Nil, nil
	}

	channel, err := r.deps.ChannelService.GetChannelByID(ctx, parsedDashboardID)
	if err == nil {
		return channel, nil
	}

	if !errors.Is(err, channelsrepo.ErrNotFound) {
		return channelentity.Nil, fmt.Errorf("get selected dashboard channel: %w", err)
	}

	// When the request is authenticated by an API key, the selected dashboard
	// holds an internal user ID instead of a channel ID, so resolve the channel
	// through the user's platform binding.
	user, userErr := r.deps.UsersService.GetByID(ctx, dashboardID)
	if userErr != nil {
		return channelentity.Nil, nil
	}

	channel, err = r.deps.ChannelService.GetChannelByBindingUserID(ctx, user.Platform, parsedDashboardID)
	if err != nil {
		return channelentity.Nil, nil
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
				Platform: r.getDashboardPlatform(ctx, channel.ID, obj.ID),
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

		ownedChannelIDs := make([]uuid.UUID, 0, len(ownChannels))
		for _, channel := range ownChannels {
			channelID, err := uuid.Parse(channel.ID)
			if err != nil {
				return nil, fmt.Errorf("parse owned channel id: %w", err)
			}
			ownedChannelIDs = append(ownedChannelIDs, channelID)
		}

		ownedChannelEntities, err := r.deps.ChannelsRepository.GetByIDs(ctx, ownedChannelIDs)
		if err != nil {
			return nil, fmt.Errorf("get owned channels: %w", err)
		}

		channelAPIKeys := make(map[string]string, len(ownedChannelEntities))
		for _, channel := range ownedChannelEntities {
			if channel.ApiKey != nil {
				channelAPIKeys[channel.ID.String()] = *channel.ApiKey
			}
		}

		for _, channel := range ownChannels {
			dashboard := gqlmodel.Dashboard{
				ID:            channel.ID,
				Platform:      r.getDashboardPlatform(ctx, channel.ID, obj.ID),
				Flags:         []gqlmodel.ChannelRolePermissionEnum{gqlmodel.ChannelRolePermissionEnumCanAccessDashboard},
				APIKey:        obj.APIKey,
				PlanID:        channel.PlanID,
				ChannelAPIKey: channelAPIKeys[channel.ID],
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
				platform = r.getDashboardPlatform(ctx, role.Role.Channel.ID, obj.ID)
			}

			dashboard := gqlmodel.Dashboard{
				ID:            role.Role.Channel.ID,
				Platform:      platform,
				Flags:         append(existing.Flags, flags...),
				APIKey:        existing.APIKey,
				PlanID:        role.Role.Channel.PlanID,
				ChannelAPIKey: existing.ChannelAPIKey,
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
				platform = r.getDashboardPlatform(ctx, role.Channel.ID, obj.ID)
			}

			dashboard := gqlmodel.Dashboard{
				ID:            role.ChannelID,
				Platform:      platform,
				Flags:         append(existing.Flags, flags...),
				APIKey:        existing.APIKey,
				PlanID:        existing.PlanID,
				ChannelAPIKey: existing.ChannelAPIKey,
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
