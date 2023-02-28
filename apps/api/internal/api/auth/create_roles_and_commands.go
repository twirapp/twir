package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"google.golang.org/protobuf/types/known/emptypb"
)

var neededRoles = []model.ChannelRoleEnum{
	model.ChannelRoleTypeBroadcaster,
	model.ChannelRoleTypeModerator,
	model.ChannelRoleTypeVip,
	model.ChannelRoleTypeSubscriber,
}

func createRolesAndCommand(services types.Services, userId string) error {
	parserGrpc := do.MustInvoke[parser.ParserClient](di.Provider)

	defaultCommands, err := parserGrpc.GetDefaultCommands(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	createdRoles := []*model.ChannelRole{}
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
		err = services.DB.Create(newRole).Error
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
			ID:          uuid.New().String(),
			ChannelID:   userId,
			Name:        command.Name,
			Description: null.StringFrom(command.Description),
			Aliases:     command.Aliases,
			RolesIDS:    roleIds,
			Enabled:     true,
		}
		err = services.DB.Create(newCommand).Error
		if err != nil {
			return err
		}
	}

	return nil
}
