package userswithstats

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/userswithstats/model"
)

type Repository interface {
	GetByUserAndChannelID(ctx context.Context, input GetByUserAndChannelIDInput) (
		model.UserWithStats,
		error,
	)
}

type GetByUserAndChannelIDInput struct {
	UserID    string
	ChannelID string
}
