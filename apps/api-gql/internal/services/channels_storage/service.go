package channels_storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/twirapp/twir/libs/repositories/channels_storage"
	"github.com/twirapp/twir/libs/repositories/channels_storage/model"
	"go.uber.org/fx"
)

const maxStorageSize = 30 * 1024 * 1024 // 30MB

type Opts struct {
	fx.In

	StorageRepository channels_storage.Repository
}

type Service struct {
	storageRepository channels_storage.Repository
}

var ErrNotFound = errors.New("storage entry not found")
var ErrStorageLimitExceeded = errors.New("storage limit exceeded")

func New(opts Opts) *Service {
	return &Service{
		storageRepository: opts.StorageRepository,
	}
}

func (s *Service) GetAllByChannelID(ctx context.Context, channelID string) (
	[]model.ChannelStorage,
	error,
) {
	return s.storageRepository.GetAllByChannelID(ctx, channelID)
}

func (s *Service) GetByKey(ctx context.Context, channelID string, key string) (
	model.ChannelStorage,
	error,
) {
	entry, err := s.storageRepository.GetByKey(ctx, channelID, key)
	if err != nil {
		if errors.Is(err, channels_storage.ErrNotFound) {
			return model.Nil, ErrNotFound
		}
		return model.Nil, err
	}

	return entry, nil
}

func (s *Service) Set(ctx context.Context, channelID string, key string, value json.RawMessage) (
	model.ChannelStorage,
	error,
) {
	currentSize, err := s.storageRepository.GetTotalSizeByChannelID(ctx, channelID)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to check storage size: %w", err)
	}

	newValueSize := int64(len(value))

	existing, err := s.storageRepository.GetByKey(ctx, channelID, key)
	if err != nil && !errors.Is(err, channels_storage.ErrNotFound) {
		return model.Nil, err
	}

	var existingSize int64
	if !existing.IsNil() {
		existingSize = int64(len(existing.Value))
	}

	if currentSize-existingSize+newValueSize > maxStorageSize {
		return model.Nil, fmt.Errorf(
			"%w: current %d bytes + new %d bytes exceeds limit of %d bytes",
			ErrStorageLimitExceeded,
			currentSize-existingSize,
			newValueSize,
			maxStorageSize,
		)
	}

	entry, err := s.storageRepository.Set(ctx, channels_storage.SetInput{
		ChannelID: channelID,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		return model.Nil, err
	}

	return entry, nil
}

func (s *Service) Delete(ctx context.Context, channelID string, key string) error {
	return s.storageRepository.Delete(ctx, channelID, key)
}

func (s *Service) DeleteAllByChannelID(ctx context.Context, channelID string) error {
	return s.storageRepository.DeleteAllByChannelID(ctx, channelID)
}

func (s *Service) GetTotalSizeByChannelID(ctx context.Context, channelID string) (int64, error) {
	return s.storageRepository.GetTotalSizeByChannelID(ctx, channelID)
}
