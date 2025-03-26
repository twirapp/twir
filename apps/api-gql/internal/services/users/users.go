package users

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/repositories/users/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	UsersRepository users.Repository
}

func New(opts Opts) *Service {
	return &Service{
		usersRepository: opts.UsersRepository,
	}
}

type Service struct {
	usersRepository users.Repository
}

type UpdateInput struct {
	IsBotAdmin        *bool
	ApiKey            *string
	IsBanned          *bool
	HideOnLandingPage *bool
	TokenID           *string
}

func (c *Service) modelToEntity(m model.User) entity.User {
	return entity.User{
		ID:                m.ID,
		TokenID:           m.TokenID.Ptr(),
		IsBotAdmin:        m.IsBotAdmin,
		ApiKey:            m.ApiKey,
		IsBanned:          m.IsBanned,
		HideOnLandingPage: m.HideOnLandingPage,
	}
}

func (c *Service) Update(ctx context.Context, id string, input UpdateInput) (entity.User, error) {
	newUser, err := c.usersRepository.Update(
		ctx,
		id,
		users.UpdateInput{
			IsBanned:          input.IsBanned,
			IsBotAdmin:        input.IsBotAdmin,
			ApiKey:            input.ApiKey,
			HideOnLandingPage: input.HideOnLandingPage,
			TokenID:           input.TokenID,
		},
	)
	if err != nil {
		return entity.UserNil, err
	}

	return c.modelToEntity(newUser), nil
}

func (c *Service) GetByID(ctx context.Context, id string) (entity.User, error) {
	user, err := c.usersRepository.GetByID(ctx, id)
	if err != nil {
		return entity.UserNil, err
	}

	return c.modelToEntity(user), nil
}

type GetManyInput struct {
	Page       int
	PerPage    int
	IDs        []string
	IsBotAdmin *bool
	IsBanned   *bool
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.User, error) {
	dbUsers, err := c.usersRepository.GetManyByIDS(
		ctx,
		users.GetManyInput{
			Page:       input.Page,
			PerPage:    input.PerPage,
			IDs:        input.IDs,
			IsBotAdmin: input.IsBotAdmin,
			IsBanned:   input.IsBanned,
		},
	)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.User, 0, len(dbUsers))
	for _, user := range dbUsers {
		entities = append(entities, c.modelToEntity(user))
	}

	return entities, nil
}
