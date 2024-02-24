package commands

import (
	"context"
	"slices"
	"time"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Commands) shouldCheckCooldown(
	badges []string,
	command *model.ChannelsCommands,
	userRoles []model.ChannelRole,
) bool {
	if command.Cooldown.Int64 == 0 {
		return false
	}

	if !lo.Contains(badges, "BROADCASTER") &&
		!lo.Contains(badges, "MODERATOR") &&
		!lo.Contains(badges, "SUBSCRIBER") &&
		!lo.Contains(badges, "VIP") {
		return true
	}

	for _, role := range command.CooldownRolesIDs {
		hasRoleForCheck := lo.SomeBy(
			userRoles,
			func(userRole model.ChannelRole) bool {
				return userRole.ID == role
			},
		)
		if !hasRoleForCheck {
			continue
		}

		return true
	}

	return false
}

func (c *Commands) prepareCooldownAndPermissionsCheck(
	ctx context.Context,
	userId,
	channelId string,
	userBadges []string,
	command *model.ChannelsCommands,
) (
	dbUser *model.Users,
	channelRoles []model.ChannelRole,
	userRoles []model.ChannelRole,
	commandRoles []model.ChannelRole,
	err error,
) {
	if err = c.services.Gorm.
		WithContext(ctx).
		Where(`"id" = ?`, userId).
		Preload("Stats", `"channelId" = ? AND "userId" = ?`, channelId, userId).
		First(&dbUser).Error; err != nil {
		return
	}

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
		hasBadge := lo.SomeBy(
			userBadges,
			func(badge string) bool {
				return badge == role.Type.String()
			},
		)

		if userHasDbRole || hasBadge {
			userRoles = append(userRoles, role)
		}

		isCommandRole := lo.SomeBy(
			command.RolesIDS, func(roleId string) bool {
				return roleId == role.ID
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
	command *model.ChannelsCommands,
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
			if role.ID != commandRole {
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
			settings := &model.ChannelRoleSettings{}
			if err := json.Unmarshal(role.Settings, settings); err != nil {
				return false
			}

			if settings.RequiredWatchTime == 0 &&
				settings.RequiredUsedChannelPoints == 0 &&
				settings.RequiredMessages == 0 {
				continue
			}

			if dbUser.Stats.UsedChannelPoints >= settings.RequiredUsedChannelPoints &&
				dbUser.Stats.Messages >= settings.RequiredMessages &&
				hoursWatched >= settings.RequiredWatchTime {
				return true
			}
		}
	}

	return false
}
