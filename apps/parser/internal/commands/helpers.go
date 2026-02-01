package commands

import (
	"context"
	"fmt"
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
	//if msg.IsChatterBroadcaster() {
	//	return false
	//}

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

	if command.Cooldown != nil && *command.Cooldown > 0 {
		return true
	}

	return false
}

type CooldownCheckResult struct {
	OnCooldown    bool
	RemainingTime int64
}

func (c *Commands) checkRoleBasedCooldown(
	ctx context.Context,
	command commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	channelId string,
	userRoles []model.ChannelRole,
) (*CooldownCheckResult, error) {
	now := time.Now().Unix()
	cooldownDuration := c.getCooldownForUser(command, userRoles)

	redisKey := fmt.Sprintf("cd:%s:%s:last_used", channelId, command.ID)

	if cooldownDuration == 0 {
		err := c.services.Redis.Set(
			ctx,
			redisKey,
			now,
			time.Hour*24,
		).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to set last_used in redis: %w", err)
		}
		return &CooldownCheckResult{OnCooldown: false, RemainingTime: 0}, nil
	}

	// Atomic check and update using Lua script to prevent race conditions
	luaScript := `
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local cooldown = tonumber(ARGV[2])

		local last_used = redis.call('GET', key)
		if not last_used then
			last_used = 0
		else
			last_used = tonumber(last_used)
		end

		local elapsed = now - last_used

		if elapsed < cooldown then
			-- On cooldown
			return {1, cooldown - elapsed}
		else
			-- Can execute, update timestamp
			redis.call('SET', key, tostring(now), 'EX', 86400)
			return {0, 0}
		end
	`

	result, err := c.services.Redis.Eval(
		ctx,
		luaScript,
		[]string{redisKey},
		now,
		cooldownDuration,
	).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check cooldown: %w", err)
	}

	resultArr := result.([]interface{})
	onCooldown := resultArr[0].(int64) == 1
	remaining := resultArr[1].(int64)

	return &CooldownCheckResult{
		OnCooldown:    onCooldown,
		RemainingTime: remaining,
	}, nil
}

func (c *Commands) getCooldownForUser(
	command commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userRoles []model.ChannelRole,
) int64 {
	if len(command.RoleCooldowns) > 0 {
		minCooldown := int64(-1)

		for _, role := range userRoles {
			for _, roleCooldown := range command.RoleCooldowns {
				if role.ID == roleCooldown.RoleID.String() {
					roleCooldownInt := int64(roleCooldown.Cooldown)
					if minCooldown == -1 || roleCooldownInt < minCooldown {
						minCooldown = roleCooldownInt
					}
				}
			}
		}

		if minCooldown != -1 {
			return minCooldown
		}
	}

	if command.Cooldown != nil {
		return int64(*command.Cooldown)
	}

	return 0
}

// todo: move to eventsub
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
