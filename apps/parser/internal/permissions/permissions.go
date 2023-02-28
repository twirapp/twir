package permissions

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)


func IsUserHasPermissionToCommand(userId, channelId string, badges []string, command *model.ChannelsCommands) bool {
	if userId == channelId {
		return true
	}

	if len(command.RolesIDS) == 0 {
		return true
	}

	db := do.MustInvoke[gorm.DB](di.Provider)


	var roles []*model.ChannelRole

	err := db.Model(&model.ChannelRole{}).
		Where(`"channelId" = ?`, channelId).
		Preload("Users", `"userId" = ?`, userId).
		Find(&roles).
		Error

	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, badge := range badges {
		for _, role := range roles {
			if role.Type.String() == badge && role.System {
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
		return true
	}


	for _, commandRole := range command.RolesIDS {
		for _, role := range roles {
			if role.ID != commandRole {
				continue
			}

			for _, user := range role.Users {
				if user.UserID == userId {
					return true
				}
			}
		}
	}

	return false
}
