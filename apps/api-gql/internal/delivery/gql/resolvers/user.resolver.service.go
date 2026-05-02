package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

var errCannotUnlinkCurrentPlatform = fmt.Errorf("cannot unlink current platform account")
var errCannotUnlinkLastPlatform = fmt.Errorf("cannot unlink the last linked platform account")

func unlinkPlatformUpdates(
	channel channelsmodel.Channel,
	currentPlatform platformentity.Platform,
	targetPlatform platformentity.Platform,
) (map[string]any, error) {
	if !targetPlatform.IsValid() {
		return nil, fmt.Errorf("unsupported platform: %s", targetPlatform)
	}

	if currentPlatform == targetPlatform {
		return nil, errCannotUnlinkCurrentPlatform
	}

	if len(channel.Platforms()) < 2 {
		return nil, errCannotUnlinkLastPlatform
	}

	switch targetPlatform {
	case platformentity.PlatformTwitch:
		if !channel.TwitchConnected() {
			return nil, fmt.Errorf("platform %s is not linked", targetPlatform)
		}

		return map[string]any{"twitch_user_id": nil}, nil
	case platformentity.PlatformKick:
		if !channel.KickConnected() {
			return nil, fmt.Errorf("platform %s is not linked", targetPlatform)
		}

		return map[string]any{
			"kick_user_id": nil,
			"kick_bot_id":  nil,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", targetPlatform)
	}
}

func getDashboardPlatform(channel model.Channels, userID string) string {
	if channel.KickUserID != nil && *channel.KickUserID == userID {
		return "kick"
	}

	if channel.TwitchUserID != nil && *channel.TwitchUserID == userID {
		return "twitch"
	}

	return channel.Platform()
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

	channel, err := r.deps.ChannelsRepository.GetByID(ctx, parsedDashboardID)
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
			Where("twitch_user_id = (SELECT id FROM users WHERE id = ?::uuid LIMIT 1)", obj.ID).
			Or("kick_user_id = (SELECT id FROM users WHERE id = ?::uuid LIMIT 1)", obj.ID).
			Find(&ownChannels).Error; err != nil {
			return nil, err
		}

		for _, channel := range ownChannels {
			dashboard := gqlmodel.Dashboard{
				ID:       channel.ID,
				Platform: getDashboardPlatform(channel, obj.ID),
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
				`"userId" = ?::uuid`,
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
				platform = role.Role.Channel.Platform()
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
		Where(`"userId" = ?::uuid`, obj.ID).
		Find(&usersStats).Error; err != nil {
		return nil, err
	}

	for _, stat := range usersStats {
		var channelRoles []model.ChannelRole
		if err := r.deps.Gorm.WithContext(ctx).Where(
			`"channelId" = ?::uuid`,
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
			platform := existing.Platform
			if platform == "" {
				platform = role.Channel.Platform()
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
