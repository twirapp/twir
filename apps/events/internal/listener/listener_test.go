package listener

import (
	"context"
	"testing"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

type channelRepositoryFake struct {
	channelsrepository.Repository

	channel                 channelentity.Channel
	lookupPlatform          platform.Platform
	lookupPlatformChannelID string
}

func (f *channelRepositoryFake) GetByPlatformChannelID(
	_ context.Context,
	p platform.Platform,
	platformChannelID string,
) (channelentity.Channel, error) {
	f.lookupPlatform = p
	f.lookupPlatformChannelID = platformChannelID
	return f.channel, nil
}

func TestResolveInternalChannelIDUsesPlatformChannelBinding(t *testing.T) {
	channelID := uuid.New()
	repo := &channelRepositoryFake{
		channel: channelentity.Channel{ID: channelID},
	}
	implementation := EventsGrpcImplementation{
		channelService: channelservice.NewChannelService(
			repo,
			&buscore.Bus{},
			cfg.Config{},
			nil,
			nil,
		),
	}

	resolvedChannelID, err := implementation.resolveInternalChannelID(
		context.Background(),
		platform.PlatformKick,
		"kick-channel",
	)
	if err != nil {
		t.Fatalf("resolveInternalChannelID returned error: %v", err)
	}
	if resolvedChannelID != channelID {
		t.Errorf("channel ID = %s, want %s", resolvedChannelID, channelID)
	}
	if repo.lookupPlatform != platform.PlatformKick {
		t.Errorf("lookup platform = %q, want %q", repo.lookupPlatform, platform.PlatformKick)
	}
	if repo.lookupPlatformChannelID != "kick-channel" {
		t.Errorf("lookup platform channel ID = %q, want %q", repo.lookupPlatformChannelID, "kick-channel")
	}
}
