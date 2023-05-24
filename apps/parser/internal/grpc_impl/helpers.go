package grpc_impl

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
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

	commandRoles := []*model.ChannelRole{}

	mappedCommandsRoles := lo.Map(command.RolesIDS, func(id string, _ int) string {
		return id
	})
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

	if lo.SomeBy(command.DeniedUsersIDS, func(id string) bool {
		return id == userId
	}) {
		return false
	}

	if lo.SomeBy(command.AllowedUsersIDS, func(id string) bool {
		return id == userId
	}) {
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

	return false
}
