package roles

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/roles"
	"github.com/twirapp/twir/libs/repositories/roles/model"
	"github.com/twirapp/twir/libs/repositories/roles_users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	RolesRepository      roles.Repository
	RolesUsersRepository roles_users.Repository
	AuditRecorder        audit.Recorder
	RolesCache           *generic_cacher.GenericCacher[[]model.Role]
}

func New(opts Opts) *Service {
	return &Service{
		rolesRepository:      opts.RolesRepository,
		rolesUsersRepository: opts.RolesUsersRepository,
		auditRecorder:        opts.AuditRecorder,
		rolesCache:           opts.RolesCache,
	}
}

type Service struct {
	rolesRepository      roles.Repository
	rolesUsersRepository roles_users.Repository
	auditRecorder        audit.Recorder
	rolesCache           *generic_cacher.GenericCacher[[]model.Role]
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

	entities := make([]entity.ChannelRole, len(dbRoles))
	for i, dbRole := range dbRoles {
		entities[i] = c.modelToEntity(dbRole)
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

	entities := make([]entity.ChannelRole, len(dbRoles))
	for i, dbRole := range dbRoles {
		entities[i] = c.modelToEntity(dbRole)
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

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(dbRole.ID.String()),
			},
			NewValue: dbRole,
		},
	)

	if err := c.rolesCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.ChannelRoleNil, fmt.Errorf("failed to invalidate roles cache: %w", err)
	}

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
	RolesIDs                  []uuid.UUID
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

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(newRole.ID.String()),
			},
			NewValue: newRole,
			OldValue: dbRole,
		},
	)

	if err := c.rolesCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.ChannelRoleNil, fmt.Errorf("failed to invalidate roles cache: %w", err)
	}

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

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelRoles),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(dbRole.ID.String()),
			},
			OldValue: dbRole,
		},
	)

	if err := c.rolesCache.Invalidate(ctx, input.ChannelID); err != nil {
		return fmt.Errorf("failed to invalidate roles cache: %w", err)
	}

	return nil
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.ChannelRole, error) {
	dbRole, err := c.rolesRepository.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelRoleNil, err
	}

	return c.modelToEntity(dbRole), nil
}
