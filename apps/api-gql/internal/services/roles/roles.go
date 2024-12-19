package roles

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/roles"
	"github.com/twirapp/twir/libs/repositories/roles/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	RolesRepository roles.Repository
}

func New(opts Opts) *Service {
	return &Service{
		rolesRepository: opts.RolesRepository,
	}
}

type Service struct {
	rolesRepository roles.Repository
}

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

	return entities, nil
}
