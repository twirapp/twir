package permissions

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

var CommandPerms = []string{"BROADCASTER", "MODERATOR", "VIP", "SUBSCRIBER", "FOLLOWER", "VIEWER"}

func IsUserHasPermissionToCommand(userId, channelId string, badges []string, command *model.ChannelsCommands) bool {
	if userId == channelId {
		return true
	}

	db := do.MustInvoke[gorm.DB](di.Provider)
	dbUser := &model.Users{}

	db.Where(`"id" = ?`, userId).Preload("DashboardAccess").Find(&dbUser)

	if dbUser.ID != "" {
		if dbUser.IsBotAdmin {
			return true
		}

		if dbUser.DashboardAccess != nil && len(dbUser.DashboardAccess) > 0 {
			for _, access := range dbUser.DashboardAccess {
				if access.ChannelID == channelId {
					return true
				}
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

	commandPermIndex := helpers.IndexOf(CommandPerms, command.Permission)

	res := false
	for _, b := range badges {
		idx := helpers.IndexOf(CommandPerms, b)

		if idx == -1 {
			continue
		}
		if idx <= commandPermIndex {
			res = true
			break
		}
	}

	return res
}
