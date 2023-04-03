package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/di"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var neededRoles = []model.ChannelRoleEnum{
	model.ChannelRoleTypeBroadcaster,
	model.ChannelRoleTypeModerator,
	model.ChannelRoleTypeVip,
	model.ChannelRoleTypeSubscriber,
}

func CreateRolesAndCommand(db *gorm.DB, userId string) error {
	parserGrpc := do.MustInvoke[parser.ParserClient](di.Provider)

	defaultCommands, err := parserGrpc.GetDefaultCommands(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	var currentRoles []model.ChannelRole
	var currentCommands []model.ChannelsCommands

	err = db.Where(`"channelId" = ?`, userId).Find(&currentCommands).Error
	if err != nil {
		return err
	}

	err = db.Where(`"channelId" = ?`, userId).Find(&currentRoles).Error
	if err != nil {
		return err
	}

	var roles []*model.ChannelRole

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, role := range neededRoles {
			existedRole, ok := lo.Find(currentRoles, func(item model.ChannelRole) bool {
				return item.Type == role
			})

			if ok {
				roles = append(roles, &existedRole)
				continue
			}

			newRole := &model.ChannelRole{
				ID:        uuid.New().String(),
				ChannelID: userId,
				Name:      role.String(),
				Type:      role,
				Permissions: lo.
					If(
						role == model.ChannelRoleTypeBroadcaster,
						pq.StringArray{model.RolePermissionCanAccessDashboard.String()},
					).
					Else(pq.StringArray{}),
			}
			err = tx.Save(newRole).Error
			if err != nil {
				return err
			}
			roles = append(roles, newRole)
		}

		for _, command := range defaultCommands.List {
			_, ok := lo.Find(currentCommands, func(item model.ChannelsCommands) bool {
				return item.DefaultName.Valid && item.DefaultName.String == command.Name
			})

			if ok {
				continue
			}

			roleIds := pq.StringArray{}

			for _, roleName := range command.RolesNames {
				role, ok := lo.Find(roles, func(role *model.ChannelRole) bool {
					return role.Type.String() == roleName
				})
				if !ok {
					continue
				}

				roleIds = append(roleIds, role.ID)
			}

			newCommand := &model.ChannelsCommands{
				ID:                 uuid.New().String(),
				ChannelID:          userId,
				Enabled:            true,
				Name:               command.Name,
				Description:        null.StringFrom(command.Description),
				Visible:            command.Visible,
				RolesIDS:           roleIds,
				Module:             command.Module,
				IsReply:            command.IsReply,
				KeepResponsesOrder: command.KeepResponsesOrder,
				Aliases:            command.Aliases,
				Default:            true,
				DefaultName:        null.StringFrom(command.Name),
				CooldownType:       "GLOBAL",
			}

			err = tx.Save(newCommand).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
