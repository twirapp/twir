package channels_commands_prefix

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.ChannelsCommandsPrefix, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelsCommandsPrefix, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
		model.ChannelsCommandsPrefix,
		error,
	)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID string
	Prefix    string
}

type UpdateInput struct {
	Prefix string
}
