package roles_with_roles_users

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TrmManager        trm.Manager
	RolesService      *roles.Service
	RolesUsersService *roles_users.Service
	Logger            logger.Logger
}

func New(opts Opts) *Service {
	return &Service{
		trmManager:        opts.TrmManager,
		rolesService:      opts.RolesService,
		rolesUsersService: opts.RolesUsersService,
		logger:            opts.Logger,
	}
}

type Service struct {
	trmManager        trm.Manager
	rolesService      *roles.Service
	rolesUsersService *roles_users.Service
	logger            logger.Logger
}

type CreateInput struct {
	Role  roles.CreateInput
	Users []CreateInputUser
}

type CreateInputUser struct {
	UserID string
}

func (c *Service) Create(ctx context.Context, input CreateInput) error {
	var newRole entity.ChannelRole

	err := c.trmManager.Do(
		ctx, func(txCtx context.Context) error {
			role, err := c.rolesService.Create(txCtx, input.Role)
			if err != nil {
				return err
			}

			newRole = role

			usersInputs := make([]roles_users.CreateInput, 0, len(input.Users))
			for _, user := range input.Users {
				usersInputs = append(
					usersInputs, roles_users.CreateInput{
						UserID: user.UserID,
						RoleID: role.ID,
					},
				)
			}

			_, err = c.rolesUsersService.CreateMany(ctx, usersInputs)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create role with users: %w", err)
	}

	c.logger.Audit(
		"Role create",
		audit.Fields{
			NewValue:      newRole,
			ActorID:       &input.Role.ActorID,
			ChannelID:     &input.Role.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(newRole.ID.String()),
		},
	)

	return nil
}
