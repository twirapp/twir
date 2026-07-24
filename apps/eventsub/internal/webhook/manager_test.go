package webhook

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	eventplatforms "github.com/twirapp/twir/apps/eventsub/internal/platforms"
	channelsmodel "github.com/twirapp/twir/libs/entities/channel"
	channelplatformsmodel "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	platformsregistry "github.com/twirapp/twir/libs/platforms"
	channels "github.com/twirapp/twir/libs/repositories/channels"
)

type bulkChannelsRepo struct {
	allByPlatformErr   error
	channelsByPlatform map[platform.Platform][]channelsmodel.Channel
	getManyErr         error
	getManyCalls       int
	platforms          []platform.Platform
}

func (r *bulkChannelsRepo) GetAllByBindingPlatform(
	_ context.Context,
	p platform.Platform,
) ([]channelsmodel.Channel, error) {
	r.platforms = append(r.platforms, p)
	return r.channelsByPlatform[p], r.allByPlatformErr
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
				transports: newTransportRegistry(
					&recordingTransport{platform: platform.PlatformKick},
				),
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

func newTransportRegistry(
	transports ...eventplatforms.EventTransport,
) *platformsregistry.Registry[eventplatforms.EventTransport] {
	registry := platformsregistry.New[eventplatforms.EventTransport]()
	for _, transport := range transports {
		registry.Register(transport)
	}

	return registry
}

type recordingTransport struct {
	platform     platform.Platform
	subscribed   []channelplatformsmodel.ChannelPlatform
	unsubscribed []channelplatformsmodel.ChannelPlatform
}

func (p *recordingTransport) Platform() platform.Platform {
	return p.platform
}

func (*recordingTransport) Capabilities() platform.Capabilities {
	return platform.Capabilities{platform.CapabilityEventsFollow}
}

func (p *recordingTransport) Subscribe(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
) error {
	p.subscribed = append(p.subscribed, binding)
	return nil
}

func (p *recordingTransport) Unsubscribe(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
) error {
	p.unsubscribed = append(p.unsubscribed, binding)
	return nil
}

func (*recordingTransport) SetCallbackBaseURL(string) {}

func TestSubscribeAllPlatformsUsesRegisteredTransportPlatform(t *testing.T) {
	twitchUserID := uuid.New()
	repo := &bulkChannelsRepo{
		channelsByPlatform: map[platform.Platform][]channelsmodel.Channel{
			platform.PlatformTwitch: {
				{
					ID: uuid.New(),
					Bindings: []channelplatformsmodel.ChannelPlatform{
						{
							Platform: platform.PlatformTwitch,
							UserID:   twitchUserID,
							Enabled:  true,
						},
					},
				},
			},
		},
	}
	transport := &recordingTransport{platform: platform.PlatformTwitch}
	manager := &Manager{
		logger:       slog.Default(),
		channelsRepo: repo,
		transports:   newTransportRegistry(transport),
	}

	if err := manager.subscribeAllPlatforms(context.Background()); err != nil {
		t.Fatalf("subscribeAllPlatforms returned error: %v", err)
	}

	if len(transport.subscribed) != 1 || transport.subscribed[0].UserID != twitchUserID {
		t.Errorf("transport subscriptions = %v, want [%s]", transport.subscribed, twitchUserID)
	}
	if len(repo.platforms) != 1 || repo.platforms[0] != platform.PlatformTwitch {
		t.Errorf("binding platform lookups = %v, want [%s]", repo.platforms, platform.PlatformTwitch)
	}
}
