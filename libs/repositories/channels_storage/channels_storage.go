package channels_storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_storage/model"
)

var ErrNotFound = errors.New("storage entry not found")

type Repository interface {
	GetAllByChannelID(ctx context.Context, channelID string) ([]model.ChannelStorage, error)
	GetByKey(ctx context.Context, channelID string, key string) (model.ChannelStorage, error)
	Set(ctx context.Context, input SetInput) (model.ChannelStorage, error)
	Delete(ctx context.Context, channelID string, key string) error
	DeleteAllByChannelID(ctx context.Context, channelID string) error
	GetTotalSizeByChannelID(ctx context.Context, channelID string) (int64, error)
}

type SetInput struct {
	ChannelID string
	Key       string
	Value     json.RawMessage
}

type UpdateInput struct {
	ID        uuid.UUID
	ChannelID string
	Key       string
	Value     json.RawMessage
}
