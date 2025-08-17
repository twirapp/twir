package users

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/users/model"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (model.User, error)
	GetManyByIDS(ctx context.Context, input GetManyInput) ([]model.User, error)
	Update(ctx context.Context, id string, input UpdateInput) (model.User, error)
	GetRandomOnlineUser(ctx context.Context, input GetRandomOnlineUserInput) (model.OnlineUser, error)
	GetByApiKey(ctx context.Context, apiKey string) (model.User, error)
	Create(ctx context.Context, input CreateInput) (model.User, error)
}

type GetManyInput struct {
	Page       int
	PerPage    int
	IDs        []string
	IsBotAdmin *bool
	IsBanned   *bool
}

type UpdateInput struct {
	IsBanned          *bool
	IsBotAdmin        *bool
	ApiKey            *string
	HideOnLandingPage *bool
	TokenID           *string
}

type GetRandomOnlineUserInput struct {
	ChannelID string
}

type CreateInput struct {
	ID                string
	ApiKey            *string
	IsBotAdmin        bool
	IsBanned          bool
	HideOnLandingPage bool
	TokenID           *string
}
