package grpc_impl

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *parserGrpcServer) shouldCheckCooldown(badges []string) bool {
	return !lo.Contains(badges, "BROADCASTER") &&
		!lo.Contains(badges, "MODERATOR") &&
		!lo.Contains(badges, "SUBSCRIBER")
}

func (c *parserGrpcServer) isUserHasPermissionToCommand(
	ctx context.Context,
	userId,
	channelId string,
	badges []string,
	command *model.ChannelsCommands,
) bool {
	if userId == channelId {
		return true
	}

	if len(command.RolesIDS) == 0 {
		return true
	}

	dbUser := &model.Users{}
	err := c.services.Gorm.
		WithContext(ctx).
		Where(`"id" = ?`, userId).
		Preload("Stats", `"channelId" = ? AND "userId" = ?`, channelId, userId).
		Find(dbUser).Error
	if err != nil {
		c.services.Logger.Sugar().Error(err)
	}

	if dbUser.IsBotAdmin {
		return true
	}

	var userRoles []*model.ChannelRole

	err = c.services.Gorm.
		WithContext(ctx).Model(&model.ChannelRole{}).
		Where(`"channelId" = ?`, channelId).
		Preload("Users", `"userId" = ?`, userId).
		Find(&userRoles).
		Error
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return false
	}

	var commandRoles []*model.ChannelRole

	mappedCommandsRoles := lo.Map(
		command.RolesIDS, func(id string, _ int) string {
			return id
		},
	)
	err = c.services.Gorm.
		WithContext(ctx).Where(`"id" IN ?`, mappedCommandsRoles).Find(&commandRoles).Error
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return false
	}

	for _, role := range commandRoles {
		if role.Type == model.ChannelRoleTypeCustom {
			continue
		}

		for _, badge := range badges {
			if strings.EqualFold(role.Type.String(), badge) {
				return true
			}
		}
	}

	if lo.SomeBy(
		command.DeniedUsersIDS, func(id string) bool {
			return id == userId
		},
	) {
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
			err = json.Unmarshal(role.Settings, settings)
			if err != nil {
				c.services.Logger.Sugar().Error(err)
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
