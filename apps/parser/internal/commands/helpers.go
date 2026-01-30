package commands

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/entities/commandrolecooldownentity"
	model "github.com/twirapp/twir/libs/gomodels"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

func (c *Commands) shouldCheckCooldown(
	msg twitch.TwitchChatMessage,
	command *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userRoles []model.ChannelRole,
) bool {
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

	if command.Cooldown != nil && *command.Cooldown > 0 {
		return true
	}

	return false
}

func (c *Commands) isCooldown(
	ctx context.Context,
	command commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userId string,
	userRoles []model.ChannelRole,
) (bool, error) {
	globalRedisKey := fmt.Sprintf("commands:%s:cooldowns:global", command.ID)

	var isShouldSetGlobalRedisCooldown bool
	defer func() {
		if isShouldSetGlobalRedisCooldown {
			c.services.Redis.Set(ctx, globalRedisKey, "", time.Duration(*command.Cooldown)*time.Second)
		}
	}()

	var userChannelRoleUser *model.ChannelRole
	var cooldownRole *commandrolecooldownentity.CommandRoleCooldown

	// roles with the lowest cooldown should be first in userRoles
	slices.SortFunc(
		userRoles, func(a, b model.ChannelRole) int {
			var q, w *commandrolecooldownentity.CommandRoleCooldown

			for _, role := range command.RoleCooldowns {
				if role.ID.String() == a.ID {
					q = &role
				}

				if role.ID.String() == b.ID {
					w = &role
				}
			}

			if q == nil && w == nil {
				return 0
			}
			if q == nil {
				return 1
			}
			if w == nil {
				return -1
			}

			return cmp.Compare(q.Cooldown, w.Cooldown)
		},
	)

	// roles with the lowest cooldown should be first in command.RoleCooldowns
	slices.SortFunc(
		command.RoleCooldowns, func(a, b commandrolecooldownentity.CommandRoleCooldown) int {
			return cmp.Compare(a.Cooldown, b.Cooldown)
		},
	)

	for _, role := range command.RoleCooldowns {
		for _, ur := range userRoles {
			if role.ID.String() == ur.ID {
				userChannelRoleUser = &ur
				break
			}
		}
	}

	if userChannelRoleUser != nil {
		for _, cr := range command.RoleCooldowns {
			if cr.RoleID.String() == userChannelRoleUser.ID {
				cooldownRole = &cr
				break
			}
		}
	}

	if command.CooldownType == "PER_USER" {
		redisKey := fmt.Sprintf("commands:%d:cooldowns:user:%s", command.ID, userId)
		exists, err := c.services.Redis.Exists(ctx, redisKey).Result()
		if err != nil {
			return false, err
		}
		if exists == 1 {
			return true, nil
		}
		var ttl int64
		if cooldownRole != nil {
			ttl = int64(cooldownRole.Cooldown)
		} else {
			ttl = int64(*command.Cooldown)
		}

		if err := c.services.Redis.Set(
			ctx,
			redisKey,
			"1",
			time.Duration(ttl)*time.Second,
		).Err(); err != nil {
			return false, err
		}

		return false, nil
	}

	globalCooldown, _ := c.services.Redis.TTL(ctx, globalRedisKey).Result()
	rolesCooldownsTTLs := make(map[string]time.Duration, len(command.RoleCooldowns))

	var cdwg sync.WaitGroup

	for _, role := range command.RoleCooldowns {
		cdwg.Go(
			func() {
				ttl, err := c.services.Redis.TTL(
					ctx,
					fmt.Sprintf("commands:%d:cooldowns:role:%s", command.ID, role.RoleID),
				).Result()
				if err != nil {
					return
				}
				// key does not exists
				if ttl == -2 {
					rolesCooldownsTTLs[role.ID.String()] = 0
					return
				}

				rolesCooldownsTTLs[role.ID.String()] = ttl
			},
		)
	}

	cdwg.Wait()

	// probably this logic incorrect
	passedFromGlobalCooldown := int(globalCooldown.Seconds()) - *command.Cooldown
	if userChannelRoleUser != nil && cooldownRole != nil {
		passedFromRoleCooldown := int(rolesCooldownsTTLs[userChannelRoleUser.ID].Seconds()) - cooldownRole.Cooldown
		isShouldSetGlobalRedisCooldown = true

		c.services.Redis.Set(
			ctx,
			fmt.Sprintf("commands:%d:cooldowns:role:%s", command.ID, userChannelRoleUser.ID),
			"",
			time.Duration(cooldownRole.Cooldown)*time.Second,
		)

		return passedFromGlobalCooldown > 0 || passedFromRoleCooldown > 0, nil
	}

	return false, nil
}

type CooldownCheckResult struct {
	OnCooldown    bool
	RemainingTime int64
}

func (c *Commands) checkRoleBasedCooldown(
	ctx context.Context,
	command commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userId string,
	channelId string,
	userRoles []model.ChannelRole,
) (*CooldownCheckResult, error) {
	now := time.Now().Unix()

	cooldownDuration := c.getCooldownForUser(command, userRoles)

	if cooldownDuration == 0 {
		return &CooldownCheckResult{OnCooldown: false, RemainingTime: 0}, nil
	}

	redisKey := fmt.Sprintf("cd:%s:%s:last_used", channelId, command.ID)

	lastUsedStr, err := c.services.Redis.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get last_used from redis: %w", err)
	}

	var lastUsed int64
	if lastUsedStr == "" {
		lastUsed = 0
	} else {
		lastUsed, err = strconv.ParseInt(lastUsedStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_used: %w", err)
		}
	}

	elapsed := now - lastUsed

	if elapsed < cooldownDuration {
		remaining := cooldownDuration - elapsed
		return &CooldownCheckResult{
			OnCooldown:    true,
			RemainingTime: remaining,
		}, nil
	}

	if err := c.services.Redis.Set(
		ctx,
		redisKey,
		strconv.FormatInt(now, 10),
		time.Duration(cooldownDuration)*time.Second,
	).Err(); err != nil {
		return nil, fmt.Errorf("failed to set last_used in redis: %w", err)
	}

	return &CooldownCheckResult{OnCooldown: false, RemainingTime: 0}, nil
}

func (c *Commands) getCooldownForUser(
	command commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
	userRoles []model.ChannelRole,
) int64 {
	if len(command.RoleCooldowns) == 0 {
		if command.Cooldown != nil {
			return int64(*command.Cooldown)
		}
		return 0
	}

	minCooldown := int64(60)

	for _, role := range userRoles {
		for _, roleCooldown := range command.RoleCooldowns {
			if role.ID == roleCooldown.RoleID.String() {
				roleCooldownInt := int64(roleCooldown.Cooldown)
				if roleCooldownInt < minCooldown {
					minCooldown = roleCooldownInt
				}
			}
		}
	}

	return minCooldown
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
