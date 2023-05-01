package permissions

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

func IsUserHasPermissionToCommand(userId, channelId string, badges []string, command *model.ChannelsCommands) bool {
	if userId == channelId {
		return true
	}

	if len(command.RolesIDS) == 0 {
		return true
	}

	db := do.MustInvoke[gorm.DB](di.Provider)

	dbUser := &model.Users{}
	err := db.Where(`"id" = ?`, userId).Find(dbUser).Error
	if err != nil {
		zap.S().Error(err)
	}

	if dbUser.IsBotAdmin {
		return true
	}

	var userRoles []*model.ChannelRole

	err = db.Model(&model.ChannelRole{}).
		Where(`"channelId" = ?`, channelId).
		Preload("Users", `"userId" = ?`, userId).
		Find(&userRoles).
		Error
	if err != nil {
		zap.S().Error(err)
		return false
	}

	commandRoles := []*model.ChannelRole{}

	mappedCommandsRoles := lo.Map(command.RolesIDS, func(id string, _ int) string {
		return id
	})
	err = db.Where(`"id" IN ?`, mappedCommandsRoles).Find(&commandRoles).Error
	if err != nil {
		zap.S().Error(err)
		return false
	}

	for _, role := range commandRoles {
		if role.Type == model.ChannelRoleTypeCustom {
			continue
		}

		for _, badge := range badges {
			if strings.ToLower(role.Type.String()) == strings.ToLower(badge) {
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
