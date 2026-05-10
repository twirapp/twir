package greetings

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID uuid.UUID, filters GetManyInput) (
		[]model.Greeting,
		error,
	)
	GetByID(ctx context.Context, id uuid.UUID) (model.Greeting, error)
	Create(ctx context.Context, input CreateInput) (model.Greeting, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Greeting, error)
	UpdateManyByChannelID(ctx context.Context, input UpdateManyInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetOneByChannelAndUserID(ctx context.Context, input GetOneInput) (model.Greeting, error)
}

type GetManyInput struct {
	Enabled   *bool
	Processed *bool
}

type CreateInput struct {
	ChannelID    uuid.UUID
	UserID       uuid.UUID
	Enabled      bool
	Text         string
	IsReply      bool
	Processed    bool
	WithShoutOut bool
}

type UpdateInput struct {
	UserID       *uuid.UUID
	Enabled      *bool
	Text         *string
	IsReply      *bool
	Processed    *bool
	WithShoutOut *bool
}

type GetOneInput struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID

	Enabled   *bool
	Processed *bool
}

type UpdateManyInput struct {
	ChannelID uuid.UUID

	Processed *bool
}
