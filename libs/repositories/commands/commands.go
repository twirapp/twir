package commands

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/commands/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Command, error)
}
