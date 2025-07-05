package channels_files

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_files/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.ChannelFile, error)
	GetMany(ctx context.Context, input GetManyInput) ([]model.ChannelFile, error)
	GetTotalChannelUploadedSizeBytes(ctx context.Context, channelID string) (int64, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelFile, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type GetManyInput struct {
	ChannelID string
}

type CreateInput struct {
	ChannelID string
	FileName  string
	MimeType  string
	Size      int64
}
