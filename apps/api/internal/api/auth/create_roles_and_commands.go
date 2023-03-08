package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var neededRoles = []model.ChannelRoleEnum{
	model.ChannelRoleTypeBroadcaster,
	model.ChannelRoleTypeModerator,
	model.ChannelRoleTypeVip,
	model.ChannelRoleTypeSubscriber,
}

func createRolesAndCommand(transaction *gorm.DB, services *types.Services, userId string) error {
	defaultCommands, err := services.Grpc.Parser.GetDefaultCommands(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	var createdRoles []*model.ChannelRole
	for _, role := range neededRoles {
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
		err = transaction.Save(newRole).Error
		if err != nil {
			return err
		}
		createdRoles = append(createdRoles, newRole)
	}
	for _, command := range defaultCommands.List {
		roleIds := pq.StringArray{}

		for _, roleName := range command.RolesNames {
			role, ok := lo.Find(createdRoles, func(role *model.ChannelRole) bool {
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

		err = transaction.Save(newCommand).Error
		if err != nil {
			return err
		}
	}

	return nil
}
