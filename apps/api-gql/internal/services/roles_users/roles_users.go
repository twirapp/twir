package roles_users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/roles_users"
	"github.com/twirapp/twir/libs/repositories/roles_users/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	RolesUsersRepository roles_users.Repository
}

func New(opts Opts) *Service {
	return &Service{
		rolesUsersRepository: opts.RolesUsersRepository,
	}
}

type Service struct {
	rolesUsersRepository roles_users.Repository
}

type CreateInput struct {
	UserID string
	RoleID uuid.UUID
}

func (c *Service) mapToEntity(m model.RoleUser) entity.ChannelRoleUser {
	return entity.ChannelRoleUser{
		ID:     m.ID,
		UserID: m.UserID,
		RoleID: m.RoleID,
	}
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.ChannelRoleUser, error) {
	user, err := c.rolesUsersRepository.Create(
		ctx, roles_users.CreateInput{
			UserID: input.UserID,
			RoleID: input.RoleID,
		},
	)
	if err != nil {
		return entity.ChannelRoleUserNil, fmt.Errorf("cannot create role user: %w", err)
	}

	return c.mapToEntity(user), nil
}

func (c *Service) CreateMany(ctx context.Context, inputs []CreateInput) (
	[]entity.ChannelRoleUser,
	error,
) {
	convertedInputs := make([]roles_users.CreateInput, len(inputs))
	for i, input := range inputs {
		convertedInputs[i] = roles_users.CreateInput{
			UserID: input.UserID,
			RoleID: input.RoleID,
		}
	}

	users, err := c.rolesUsersRepository.CreateMany(ctx, convertedInputs)
	if err != nil {
		return nil, fmt.Errorf("cannot create role users: %w", err)
	}

	result := make([]entity.ChannelRoleUser, len(users))
	for i, u := range users {
		result[i] = c.mapToEntity(u)
	}

	return result, nil
}
