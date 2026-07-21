package webhook

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channels "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type bulkChannelsRepo struct {
	allByPlatformErr error
	getManyErr       error
	getManyCalls     int
	platforms        []platform.Platform
}

func (r *bulkChannelsRepo) GetAllByBindingPlatform(
	_ context.Context,
	p platform.Platform,
) ([]channelsmodel.Channel, error) {
	r.platforms = append(r.platforms, p)
	return nil, r.allByPlatformErr
}

func (r *bulkChannelsRepo) GetMany(context.Context, channels.GetManyInput) ([]channelsmodel.Channel, error) {
	r.getManyCalls++
	return nil, r.getManyErr
}

func (r *bulkChannelsRepo) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *bulkChannelsRepo) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *bulkChannelsRepo) GetByBindingUserID(context.Context, platform.Platform, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *bulkChannelsRepo) GetByPlatformChannelID(context.Context, platform.Platform, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *bulkChannelsRepo) GetBySlug(context.Context, channels.GetBySlugInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *bulkChannelsRepo) GetCount(context.Context, channels.GetCountInput) (int, error) {
	return 0, nil
}

func (r *bulkChannelsRepo) Update(context.Context, uuid.UUID, channels.UpdateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (r *bulkChannelsRepo) Create(context.Context, channels.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func TestBulkKickOperationsUseCompleteBindingList(t *testing.T) {
	allByPlatformErr := errors.New("complete binding list")
	getManyErr := errors.New("legacy paged list")

	tests := []struct {
		name string
		run  func(*Manager, context.Context) error
	}{
		{name: "subscribe", run: (*Manager).subscribeAllPlatforms},
		{name: "unsubscribe", run: (*Manager).unsubscribeAllPlatforms},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &bulkChannelsRepo{
				allByPlatformErr: allByPlatformErr,
				getManyErr:       getManyErr,
			}
			manager := &Manager{
				logger:       slog.Default(),
				channelsRepo: repo,
			}

			err := tt.run(manager, context.Background())
			if !errors.Is(err, allByPlatformErr) {
				t.Fatalf("bulk operation error = %v, want %v", err, allByPlatformErr)
			}
			if repo.getManyCalls != 0 {
				t.Fatalf("GetMany calls = %d, want 0", repo.getManyCalls)
			}
			if len(repo.platforms) != 1 || repo.platforms[0] != platform.PlatformKick {
				t.Fatalf("binding platforms = %v, want [%s]", repo.platforms, platform.PlatformKick)
			}
		})
	}
}
