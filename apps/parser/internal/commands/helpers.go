package commands

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

func (c *Commands) shouldCheckCooldown(
	msg twitch.TwitchChatMessage,
	command *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userRoles []model.ChannelRole,
) bool {
	if command.Cooldown == nil || *command.Cooldown == 0 {
		return false
	}

	if msg.IsChatterBroadcaster() {
		return false
	}

	for _, role := range command.RoleCooldowns {
		hasRoleForCheck := lo.SomeBy(
			userRoles,
			func(userRole model.ChannelRole) bool {
				return userRole.ID == role.ID.String()
			},
		)
		if !hasRoleForCheck {
			continue
		}

		return true
	}

	return false
}

// getRoleCooldown returns the cooldown for a specific role if set, otherwise returns nil
func (c *Commands) getRoleCooldown(
	command *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userRoles []model.ChannelRole,
) (*string, *int) {
	if len(command.RoleCooldowns) == 0 {
		return nil, command.Cooldown
	}

	for _, userRole := range userRoles {
		for _, roleCooldown := range command.RoleCooldowns {
			if roleCooldown.RoleID.String() == userRole.ID {
				return lo.ToPtr(roleCooldown.RoleID.String()), &roleCooldown.Cooldown
			}
		}
	}

	// If user has no roles with custom cooldowns, use default
	return nil, command.Cooldown
}

func (c *Commands) prepareCooldownAndPermissionsCheck(
	ctx context.Context,
	userId,
	channelId string,
	msg twitch.TwitchChatMessage,
	command *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
) (
	channelRoles []model.ChannelRole,
	userRoles []model.ChannelRole,
	commandRoles []model.ChannelRole,
	err error,
) {
	if err = c.services.Gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, channelId).
		Preload("Users", `"userId" = ?`, userId).
		Find(&channelRoles).Error; err != nil {
		return
	}

	for _, role := range channelRoles {
		userHasDbRole := lo.SomeBy(
			role.Users,
			func(user *model.ChannelRoleUser) bool {
				return user.UserID == userId
			},
		)

		hasBadge := msg.HasRoleFromDbByType(role.Type.String())

		if userHasDbRole || hasBadge {
			userRoles = append(userRoles, role)
		}

		isCommandRole := lo.SomeBy(
			command.RolesIDS, func(roleId uuid.UUID) bool {
				return roleId.String() == role.ID
			},
		)
		if isCommandRole {
			commandRoles = append(commandRoles, role)
		}
	}

	return
}

func (c *Commands) isUserHasPermissionToCommand(
	userId,
	channelId string,
	command *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	dbUser *model.Users,
	userRoles []model.ChannelRole,
	commandRoles []model.ChannelRole,
) bool {
	if userId == channelId {
		return true
	}

	if dbUser.IsBotAdmin {
		return true
	}

	if slices.Contains(command.DeniedUsersIDS, userId) {
		return false
	}

	if lo.SomeBy(
		command.AllowedUsersIDS, func(id string) bool {
			return id == userId
		},
	) {
		// allowed user
		return true
	}

	if len(command.RolesIDS) == 0 {
		return true
	}

	for _, role := range commandRoles {
		userHasRole := lo.SomeBy(
			userRoles,
			func(userRole model.ChannelRole) bool {
				return userRole.ID == role.ID
			},
		)

		if userHasRole {
			return true
		}
	}

	for _, commandRole := range command.RolesIDS {
		for _, role := range userRoles {
			if role.ID != commandRole.String() {
				continue
			}

			for _, user := range role.Users {
				if user.UserID == userId {
					// user in role
					return true
				}
			}
		}
	}

	if dbUser.Stats != nil {
		watched := time.Duration(dbUser.Stats.Watched) * time.Millisecond
		hoursWatched := int64(watched.Hours())

		// check command restriction by stats
		if (command.RequiredWatchTime > 0 || command.RequiredMessages > 0 || command.RequiredUsedChannelPoints > 0) &&
			dbUser.Stats.UsedChannelPoints >= int64(command.RequiredUsedChannelPoints) &&
			dbUser.Stats.Messages >= int32(command.RequiredMessages) &&
			hoursWatched >= int64(command.RequiredWatchTime) {
			return true
		}

		// check role restriction by stats
		for _, role := range commandRoles {
			if role.RequiredWatchTime == 0 &&
				role.RequiredUsedChannelPoints == 0 &&
				role.RequiredMessages == 0 {
				continue
			}

			if dbUser.Stats.UsedChannelPoints >= role.RequiredUsedChannelPoints &&
				dbUser.Stats.Messages >= role.RequiredMessages &&
				hoursWatched >= role.RequiredWatchTime {
				return true
			}
		}
	}

	return false
}
