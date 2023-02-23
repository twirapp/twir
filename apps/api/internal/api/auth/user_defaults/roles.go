package user_defaults

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateRoles(userId string) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	var roles []model.ChannelRole
	db.Where(`"channelId" = ?`, userId).Find(&roles)

	if roles != nil && len(roles) > 0 {
		return
	}

	roles = []model.ChannelRole{
		{
			ID:        uuid.NewV4().String(),
			ChannelID: userId,
			Name:      "Broadcaster",
			Type:      model.ChannelRoleTypeBroadcaster,
			System:    true,
			Permissions: []*model.ChannelRolePermission{
				{
					ID:   uuid.NewV4().String(),
					Flag: &model.RoleFlag{Flag: model.RolePermissionAdministrator},
				},
			},
		},
		{
			ID:          uuid.NewV4().String(),
			ChannelID:   userId,
			Name:        "Moderator",
			Type:        model.ChannelRoleTypeModerator,
			System:      true,
			Permissions: nil,
		},
		{
			ID:          uuid.NewV4().String(),
			ChannelID:   userId,
			Name:        "Subscriber",
			Type:        model.ChannelRoleTypeSubscriber,
			System:      true,
			Permissions: nil,
		},
		{
			ID:          uuid.NewV4().String(),
			ChannelID:   userId,
			Name:        "VIP",
			Type:        model.ChannelRoleTypeVip,
			System:      true,
			Permissions: nil,
		},
	}

	err := db.Create(&roles).Error

	if err != nil {
		logger.Error(err)
	}
}
