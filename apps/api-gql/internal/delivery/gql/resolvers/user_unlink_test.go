package resolvers

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestUnlinkPlatformAccountRejectsCurrentPlatform(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	kickBinding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true,
	}
	operations := &resolverChannelPlatformBindingsService{binding: kickBinding}
	resolver := newChannelPlatformTestResolverWithDependencies(
		dashboardID,
		channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{
			kickBinding,
			{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		}},
		map[uuid.UUID]usersmodel.User{},
		testResolverPlatformRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		operations,
	)
	resolver.deps.CurrentPlatform = resolverCurrentPlatformGetter{platform: platformentity.PlatformKick.String()}

	_, err := (&mutationResolver{resolver}).UnlinkPlatformAccount(context.Background(), platformentity.PlatformKick.String())
	if !errors.Is(err, errCannotUnlinkCurrentPlatform) {
		t.Fatalf("UnlinkPlatformAccount() error = %v, want errCannotUnlinkCurrentPlatform", err)
	}
	if operations.deletedID != uuid.Nil {
		t.Fatalf("UnlinkPlatformAccount() deleted current platform binding %s", operations.deletedID)
	}
}

type resolverCurrentPlatformGetter struct {
	platform string
	err      error
}

func (r resolverCurrentPlatformGetter) GetCurrentPlatform(context.Context) (string, error) {
	return r.platform, r.err
}
