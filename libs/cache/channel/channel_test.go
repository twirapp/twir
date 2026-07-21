package channel

import (
	"context"
	"testing"

	"github.com/google/uuid"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestNewByTwitchUserIDUsesTwitchPlatformChannelLookup(t *testing.T) {
	channelID := uuid.New()
	repository := &fakeChannelsRepository{
		channel: channelsmodel.Channel{ID: channelID},
	}
	cacher := NewByTwitchUserID(repository, kvinmemory.New())

	got, err := cacher.Get(context.Background(), "twitch-channel")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != channelID {
		t.Fatalf("Get() channel ID = %s, want %s", got.ID, channelID)
	}
	if len(repository.platformChannelLookups) != 1 {
		t.Fatalf("GetByPlatformChannelID calls = %d, want 1", len(repository.platformChannelLookups))
	}

	lookup := repository.platformChannelLookups[0]
	if lookup.platform != platform.PlatformTwitch || lookup.platformChannelID != "twitch-channel" {
		t.Fatalf(
			"GetByPlatformChannelID(%s, %q), want (%s, %q)",
			lookup.platform,
			lookup.platformChannelID,
			platform.PlatformTwitch,
			"twitch-channel",
		)
	}
}

type fakeChannelsRepository struct {
	channel                channelsmodel.Channel
	platformChannelLookups []platformChannelLookup
}

type platformChannelLookup struct {
	platform          platform.Platform
	platformChannelID string
}

func (f *fakeChannelsRepository) GetByPlatformChannelID(
	_ context.Context,
	p platform.Platform,
	platformChannelID string,
) (channelsmodel.Channel, error) {
	f.platformChannelLookups = append(
		f.platformChannelLookups,
		platformChannelLookup{platform: p, platformChannelID: platformChannelID},
	)

	return f.channel, nil
}

func (f *fakeChannelsRepository) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (f *fakeChannelsRepository) GetAllByBindingPlatform(context.Context, platform.Platform) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (f *fakeChannelsRepository) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) GetByBindingUserID(
	context.Context,
	platform.Platform,
	uuid.UUID,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) GetBySlug(context.Context, channelsrepo.GetBySlugInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (f *fakeChannelsRepository) Update(
	context.Context,
	uuid.UUID,
	channelsrepo.UpdateInput,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) Create(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}
