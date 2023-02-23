package permissions

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func IsUserHasPermissionToCommand(userId, channelId string, badges []string, commandRoles []string) bool {
	if userId == channelId {
		return true
	}

	if len(commandRoles) == 0 {
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

	for _, role := range roles {
		for _, users := range role.Users {
			if users.UserID == userId {
				return true
			}
		}
	}

	return false
}
