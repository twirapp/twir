package channelservice

import (
	"context"
	"testing"

	"github.com/google/uuid"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	config "github.com/twirapp/twir/libs/config"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/repositories/streams/model"
)

type fakeStreamsRepo struct {
	streamsrepository.Repository
	streams []model.Stream
}

func (f *fakeStreamsRepo) GetListByChannelID(_ context.Context, _ uuid.UUID) ([]model.Stream, error) {
	return f.streams, nil
}

func TestIsChannelOnlineFallsBackToRepositoryOnCacheMiss(t *testing.T) {
	channelID := uuid.New()
	service := NewChannelService(
		nil,
		nil,
		config.Config{},
		kvinmemory.New(),
		&fakeStreamsRepo{streams: []model.Stream{{ID: "stream-1", ChannelID: channelID}}},
	)

	online, err := service.IsChannelOnline(context.Background(), channelID)
	if err != nil {
		t.Fatalf("IsChannelOnline: %v", err)
	}
	if !online {
		t.Fatal("expected online=true when cache misses but repository has a live stream")
	}
}

func TestIsChannelOnlineReturnsFalseWhenNoStreams(t *testing.T) {
	channelID := uuid.New()
	service := NewChannelService(
		nil,
		nil,
		config.Config{},
		kvinmemory.New(),
		&fakeStreamsRepo{},
	)

	online, err := service.IsChannelOnline(context.Background(), channelID)
	if err != nil {
		t.Fatalf("IsChannelOnline: %v", err)
	}
	if online {
		t.Fatal("expected online=false when neither cache nor repository have a stream")
	}
}
