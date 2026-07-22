package resolvers

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestUnlinkPlatformAccountRejectsCurrentPlatform(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	kickBinding := channelplatformsmodel.ChannelPlatform{
		ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true,
	}
	repository := &resolverBindingRepository{binding: kickBinding}
	resolver := newChannelPlatformTestResolverWithDependencies(
		dashboardID,
		channelsmodel.Channel{ID: dashboardID, Bindings: []channelplatformsmodel.ChannelPlatform{
			kickBinding,
			{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		}},
		map[uuid.UUID]usersmodel.User{},
		testResolverPlatformRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		repository,
		&resolverOAuthStarter{url: "https://provider.example/authorize"},
	)
	resolver.deps.CurrentPlatform = resolverCurrentPlatformGetter{platform: platformentity.PlatformKick.String()}

	_, err := (&mutationResolver{resolver}).UnlinkPlatformAccount(context.Background(), platformentity.PlatformKick.String())
	if !errors.Is(err, errCannotUnlinkCurrentPlatform) {
		t.Fatalf("UnlinkPlatformAccount() error = %v, want errCannotUnlinkCurrentPlatform", err)
	}
	if repository.deletedID != uuid.Nil {
		t.Fatalf("UnlinkPlatformAccount() deleted current platform binding %s", repository.deletedID)
	}
}

type resolverCurrentPlatformGetter struct {
	platform string
	err      error
}

func (r resolverCurrentPlatformGetter) GetCurrentPlatform(context.Context) (string, error) {
	return r.platform, r.err
}
