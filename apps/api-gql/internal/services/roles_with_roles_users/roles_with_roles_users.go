package roles_with_roles_users

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/logger"
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
	err := c.trmManager.Do(
		ctx, func(txCtx context.Context) error {
			role, err := c.rolesService.Create(txCtx, input.Role)
			if err != nil {
				return err
			}

			usersInputs := make([]roles_users.CreateInput, 0, len(input.Users))
			for _, user := range input.Users {
				usersInputs = append(
					usersInputs, roles_users.CreateInput{
						UserID: user.UserID,
						RoleID: role.ID,
					},
				)
			}

			_, err = c.rolesUsersService.CreateMany(txCtx, usersInputs)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create role with users: %w", err)
	}

	return nil
}

type UpdateInput struct {
	ID        uuid.UUID
	ChannelID string
	ActorID   string

	Role  roles.UpdateInput
	Users []CreateInputUser
}

func (c *Service) Update(ctx context.Context, input UpdateInput) error {
	dbRole, err := c.rolesService.GetByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	if dbRole.ChannelID != input.ChannelID {
		return fmt.Errorf("role doesn't belong to the channel")
	}

	var newRole entity.ChannelRole
	err = c.trmManager.Do(
		ctx,
		func(txCtx context.Context) error {
			newDbRole, err := c.rolesService.Update(
				txCtx,
				input.ID,
				roles.UpdateInput{
					ChannelID:                 input.ChannelID,
					ActorID:                   input.ActorID,
					Name:                      input.Role.Name,
					Permissions:               input.Role.Permissions,
					RequiredWatchTime:         input.Role.RequiredWatchTime,
					RequiredMessages:          input.Role.RequiredMessages,
					RequiredUsedChannelPoints: input.Role.RequiredUsedChannelPoints,
				},
			)
			if err != nil {
				return err
			}

			newRole = newDbRole

			err = c.rolesUsersService.DeleteManyByRoleID(txCtx, newRole.ID)
			if err != nil {
				return err
			}

			usersInputs := make([]roles_users.CreateInput, 0, len(input.Users))
			for _, user := range input.Users {
				usersInputs = append(
					usersInputs, roles_users.CreateInput{
						UserID: user.UserID,
						RoleID: newRole.ID,
					},
				)
			}

			_, err = c.rolesUsersService.CreateMany(txCtx, usersInputs)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update role with users: %w", err)
	}

	return nil
}
