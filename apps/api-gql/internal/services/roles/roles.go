package roles

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/roles"
	"github.com/twirapp/twir/libs/repositories/roles/model"
	"github.com/twirapp/twir/libs/repositories/roles_users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	RolesRepository      roles.Repository
	RolesUsersRepository roles_users.Repository
	Logger               logger.Logger
}

func New(opts Opts) *Service {
	return &Service{
		rolesRepository:      opts.RolesRepository,
		rolesUsersRepository: opts.RolesUsersRepository,
		logger:               opts.Logger,
	}
}

type Service struct {
	rolesRepository      roles.Repository
	rolesUsersRepository roles_users.Repository
	logger               logger.Logger
}

var maxRoles = 20

func (c *Service) modelToEntity(m model.Role) entity.ChannelRole {
	return entity.ChannelRole{
		ID:                        m.ID,
		ChannelID:                 m.ChannelID,
		Name:                      m.Name,
		Type:                      entity.ChannelRoleEnum(m.Type.String()),
		Permissions:               m.Permissions,
		RequiredWatchTime:         m.RequiredWatchTime,
		RequiredMessages:          m.RequiredMessages,
		RequiredUsedChannelPoints: m.RequiredUsedChannelPoints,
	}
}

func (c *Service) GetManyByIDS(ctx context.Context, ids []uuid.UUID) ([]entity.ChannelRole, error) {
	dbRoles, err := c.rolesRepository.GetManyByIDS(ctx, ids)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.ChannelRole, 0, len(dbRoles))
	for _, dbRole := range dbRoles {
		entities = append(entities, c.modelToEntity(dbRole))
	}

	return entities, nil
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.ChannelRole,
	error,
) {
	dbRoles, err := c.rolesRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.ChannelRole, 0, len(dbRoles))
	for _, dbRole := range dbRoles {
		entities = append(entities, c.modelToEntity(dbRole))
	}

	slices.SortFunc(
		entities,
		func(a, b entity.ChannelRole) int {
			typeIdx := lo.IndexOf(entity.AllChannelRoleTypeEnum, a.Type)

			return typeIdx - lo.IndexOf(entity.AllChannelRoleTypeEnum, b.Type)
		},
	)

	return entities, nil
}

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name                      string
	Type                      entity.ChannelRoleEnum
	Permissions               []string
	RequiredWatchTime         int64
	RequiredMessages          int32
	RequiredUsedChannelPoints int64
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.ChannelRole, error) {
	dbRoles, err := c.rolesRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.ChannelRoleNil, err
	}

	if len(dbRoles) >= maxRoles {
		return entity.ChannelRoleNil, fmt.Errorf("maximum number of roles reached")
	}

	dbRole, err := c.rolesRepository.Create(
		ctx, roles.CreateInput{
			ChannelID:                 input.ChannelID,
			Name:                      input.Name,
			Type:                      model.ChannelRoleEnum(input.Type.String()),
			Permissions:               input.Permissions,
			RequiredWatchTime:         input.RequiredWatchTime,
			RequiredMessages:          input.RequiredMessages,
			RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
		},
	)
	if err != nil {
		return entity.ChannelRole{}, err
	}

	c.logger.Audit(
		"Role create",
		audit.Fields{
			NewValue:      dbRole,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(dbRole.ID.String()),
		},
	)

	return c.modelToEntity(dbRole), nil
}

type UpdateInput struct {
	ChannelID string
	ActorID   string

	Name                      *string
	Permissions               []string
	RequiredWatchTime         *int64
	RequiredMessages          *int32
	RequiredUsedChannelPoints *int64
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	entity.ChannelRole,
	error,
) {
	dbRole, err := c.rolesRepository.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelRoleNil, err
	}

	if dbRole.ChannelID != input.ChannelID {
		return entity.ChannelRoleNil, fmt.Errorf("role doesn't belong to the channel")
	}

	updateInput := roles.UpdateInput{
		Name:                      input.Name,
		Permissions:               input.Permissions,
		RequiredWatchTime:         input.RequiredWatchTime,
		RequiredMessages:          input.RequiredMessages,
		RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
	}

	newRole, err := c.rolesRepository.Update(ctx, id, updateInput)
	if err != nil {
		return entity.ChannelRole{}, err
	}

	c.logger.Audit(
		"Role update",
		audit.Fields{
			OldValue:      dbRole,
			NewValue:      newRole,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(newRole.ID.String()),
		},
	)

	return c.modelToEntity(newRole), nil
}

type DeleteInput struct {
	ChannelID string
	ActorID   string
	ID        uuid.UUID
}

func (c *Service) Delete(ctx context.Context, input DeleteInput) error {
	dbRole, err := c.rolesRepository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if dbRole.ChannelID != input.ChannelID {
		return fmt.Errorf("role doesn't belong to the channel")
	}

	if dbRole.Type != model.ChannelRoleTypeCustom {
		return fmt.Errorf("cannot remove default roles")
	}

	if err := c.rolesRepository.Delete(ctx, input.ID); err != nil {
		return err
	}

	c.logger.Audit(
		"Role remove",
		audit.Fields{
			OldValue:      dbRole,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(dbRole.ID.String()),
		},
	)

	return nil
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.ChannelRole, error) {
	dbRole, err := c.rolesRepository.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelRoleNil, err
	}

	return c.modelToEntity(dbRole), nil
}
