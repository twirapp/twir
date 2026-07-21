package manager

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

type fakeChannelsRepo struct {
	channel channelsmodel.Channel
	err     error
}

func (f *fakeChannelsRepo) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (f *fakeChannelsRepo) GetByID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	if f.err != nil {
		return channelsmodel.Nil, f.err
	}

	return f.channel, nil
}

func (f *fakeChannelsRepo) GetByBindingUserID(context.Context, platform.Platform, uuid.UUID) (channelsmodel.Channel, error) {
	return f.channel, f.err
}

func (f *fakeChannelsRepo) GetByPlatformChannelID(context.Context, platform.Platform, string) (channelsmodel.Channel, error) {
	return f.channel, f.err
}

func (f *fakeChannelsRepo) GetBySlug(context.Context, channelsrepo.GetBySlugInput) (channelsmodel.Channel, error) {
	return f.channel, f.err
}

func (f *fakeChannelsRepo) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (f *fakeChannelsRepo) Update(context.Context, uuid.UUID, channelsrepo.UpdateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) Create(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func newTestManager(repo channelsrepo.Repository) *Manager {
	return &Manager{
		channelsRepo: repo,
		channelService: channelservice.NewChannelService(
			repo,
			&buscore.Bus{},
			cfg.Config{},
			kvinmemory.New(),
			nil,
		),
	}
}

func TestResolveTwitchSubscriptionIdentitiesSelectsTwitchBinding(t *testing.T) {
	twitchPlatformID := "tw-123"
	botID := "bot-456"
	m := newTestManager(
		&fakeChannelsRepo{
			channel: channelsmodel.Channel{
				Bindings: []channelplatformsmodel.ChannelPlatform{
					{
						Platform:          platform.PlatformKick,
						UserID:            uuid.New(),
						PlatformChannelID: "kick-123",
						Enabled:           true,
					},
					{
						Platform:          platform.PlatformTwitch,
						UserID:            uuid.New(),
						PlatformChannelID: twitchPlatformID,
						Enabled:           true,
						BotConfig:         json.RawMessage(`{"bot_id":"bot-456"}`),
					},
				},
			},
		},
	)

	channelID := uuid.New().String()
	broadcasterID, resolvedBotID, err := m.resolveTwitchSubscriptionIdentities(context.Background(), channelID)
	if err != nil {
		t.Fatalf("resolveTwitchSubscriptionIdentities returned error: %v", err)
	}

	if broadcasterID != twitchPlatformID {
		t.Fatalf("expected broadcaster ID %q, got %q", twitchPlatformID, broadcasterID)
	}

	if resolvedBotID != botID {
		t.Fatalf("expected bot ID %q, got %q", botID, resolvedBotID)
	}
}

func TestResolveTwitchSubscriptionIdentitiesRequiresPlatformID(t *testing.T) {
	m := newTestManager(&fakeChannelsRepo{})

	_, _, err := m.resolveTwitchSubscriptionIdentities(context.Background(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error when Twitch platform ID is missing")
	}
}

func TestResolveTwitchSubscriptionIdentitiesReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("boom")
	m := newTestManager(&fakeChannelsRepo{err: wantErr})

	_, _, err := m.resolveTwitchSubscriptionIdentities(context.Background(), uuid.New().String())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}
