package resolvers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestLinkedAccountsOmitsUnavailablePlatformBindings(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchUserID := uuid.New()
	vkUserID := uuid.New()
	channel := channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: twitchUserID, PlatformChannelID: "twitch-channel", Enabled: true},
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true},
	}}
	service := channelplatformservice.New(
		resolverChannelReader{channel: channel},
		resolverUserLookup{users: map[uuid.UUID]usersmodel.User{
			twitchUserID: {ID: twitchUserID, PlatformID: "twitch-user", Login: "twitch-login"},
			vkUserID:     {ID: vkUserID, PlatformID: "vk-user", Login: "vk-login"},
		}},
		nil,
		nil,
		appplatform.NewRegistry([]appplatform.PlatformProvider{resolverPlatformProvider{platform: platformentity.PlatformTwitch}}),
		nil,
	)
	resolver := &authenticatedUserResolver{Resolver: &Resolver{deps: Deps{ChannelPlatformBindingsService: service}}}

	accounts, err := resolver.linkedAccountsForChannel(context.Background(), channel)
	if err != nil {
		t.Fatalf("linkedAccountsForChannel() error = %v", err)
	}
	if len(accounts) != 1 || accounts[0].Platform != platformentity.PlatformTwitch.String() || accounts[0].PlatformUserID != "twitch-user" {
		t.Fatalf("linkedAccountsForChannel() = %#v, want only registered Twitch account", accounts)
	}
}
