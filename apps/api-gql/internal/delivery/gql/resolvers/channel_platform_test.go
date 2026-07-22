package resolvers

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestChannelPlatformBindingsReturnsProfilesEnabledStateAndCapabilities(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	vkUserID := uuid.New()
	resolver := newChannelPlatformTestResolver(
		dashboardID,
		channelsmodel.Channel{
			ID: dashboardID,
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: twitchUserID, PlatformChannelID: "twitch-channel", Enabled: true},
				{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: kickUserID, PlatformChannelID: "kick-channel", Enabled: false},
				{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true},
			},
		},
		map[uuid.UUID]usersmodel.User{
			twitchUserID: {ID: twitchUserID, PlatformID: "twitch-user", Login: "twitch-login", DisplayName: "Twitch Name", Avatar: "https://example.com/twitch.png"},
			kickUserID:   {ID: kickUserID, PlatformID: "kick-user", Login: "kick-login", DisplayName: "Kick Name"},
			vkUserID:     {ID: vkUserID, PlatformID: "vk-user", DisplayName: "VK Name", Avatar: "https://example.com/vk.png"},
		},
		testResolverPlatformRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick, platformentity.PlatformVKVideoLive),
	)

	got, err := (&queryResolver{resolver}).ChannelPlatformBindings(context.Background())
	if err != nil {
		t.Fatalf("ChannelPlatformBindings() error = %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("ChannelPlatformBindings() returned %d bindings, want 3", len(got))
	}

	if got[0].Platform != gqlmodel.PlatformTwitch || got[0].PlatformDisplayName != "Twitch Name" || !got[0].Enabled || got[0].PlatformAvatar == nil || *got[0].PlatformAvatar != "https://example.com/twitch.png" {
		t.Fatalf("Twitch binding = %#v", got[0])
	}
	if got[1].Platform != gqlmodel.PlatformKick || got[1].PlatformLogin != "kick-login" || got[1].Enabled {
		t.Fatalf("Kick binding = %#v", got[1])
	}
	if got[2].Platform != gqlmodel.PlatformVkVideoLive || got[2].PlatformDisplayName != "VK Name" || !got[2].Enabled {
		t.Fatalf("VK binding = %#v", got[2])
	}
	if len(got[0].Capabilities) != 8 || got[0].Capabilities[0].Name != "chat.read" || got[0].Capabilities[7].Name != "events.reward" {
		t.Fatalf("Twitch capability strings = %#v", got[0].Capabilities)
	}
}

func TestChannelPlatformBindingsOmitsVKWhenFeatureIsDisabled(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	twitchUserID := uuid.New()
	vkUserID := uuid.New()
	resolver := newChannelPlatformTestResolver(
		dashboardID,
		channelsmodel.Channel{
			ID: dashboardID,
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: twitchUserID, PlatformChannelID: "twitch-channel", Enabled: true},
				{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true},
			},
		},
		map[uuid.UUID]usersmodel.User{
			twitchUserID: {ID: twitchUserID, PlatformID: "twitch-user"},
			vkUserID:     {ID: vkUserID, PlatformID: "vk-user"},
		},
		testResolverPlatformRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
	)

	got, err := (&queryResolver{resolver}).ChannelPlatformBindings(context.Background())
	if err != nil {
		t.Fatalf("ChannelPlatformBindings() error = %v", err)
	}
	if len(got) != 1 || got[0].Platform != gqlmodel.PlatformTwitch {
		t.Fatalf("ChannelPlatformBindings() = %#v, want only enabled-platform Twitch binding", got)
	}
}

func TestChannelPlatformBindingMutationsUseGenericOperations(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	kickUserID := uuid.New()
	kickBinding := channelplatformsmodel.ChannelPlatform{
		ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: kickUserID, PlatformChannelID: "kick-channel", Enabled: true,
	}
	repository := &resolverBindingRepository{binding: kickBinding}
	oauth := &resolverOAuthStarter{url: "https://provider.example/authorize"}
	resolver := newChannelPlatformTestResolverWithDependencies(
		dashboardID,
		channelsmodel.Channel{ID: dashboardID, Bindings: []channelplatformsmodel.ChannelPlatform{
			kickBinding,
			{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		}},
		map[uuid.UUID]usersmodel.User{kickUserID: {ID: kickUserID, PlatformID: "kick-user", DisplayName: "Kick Name"}},
		testResolverPlatformRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		repository,
		oauth,
	)

	mutation := &mutationResolver{resolver}
	url, err := mutation.ChannelPlatformConnect(context.Background(), gqlmodel.PlatformKick, "/dashboard/platforms")
	if err != nil {
		t.Fatalf("ChannelPlatformConnect() error = %v", err)
	}
	if url != oauth.url || oauth.channelID != dashboardID || oauth.platform != platformentity.PlatformKick || oauth.redirectTo != "/dashboard/platforms" {
		t.Fatalf("ChannelPlatformConnect() = %q, oauth = %#v", url, oauth)
	}

	updated, err := mutation.ChannelPlatformSetEnabled(context.Background(), gqlmodel.PlatformKick, false)
	if err != nil {
		t.Fatalf("ChannelPlatformSetEnabled() error = %v", err)
	}
	if updated == nil || updated.Enabled || repository.patch.Enabled == nil || *repository.patch.Enabled {
		t.Fatalf("ChannelPlatformSetEnabled() = %#v, patch = %#v", updated, repository.patch)
	}

	disconnected, err := mutation.ChannelPlatformDisconnect(context.Background(), gqlmodel.PlatformKick)
	if err != nil {
		t.Fatalf("ChannelPlatformDisconnect() error = %v", err)
	}
	if !disconnected || repository.deletedID != kickBinding.ID {
		t.Fatalf("ChannelPlatformDisconnect() = %t, deleted = %s", disconnected, repository.deletedID)
	}
}

func newChannelPlatformTestResolver(
	dashboardID uuid.UUID,
	channel channelsmodel.Channel,
	users map[uuid.UUID]usersmodel.User,
	registry *appplatform.Registry,
) *Resolver {
	return newChannelPlatformTestResolverWithDependencies(
		dashboardID,
		channel,
		users,
		registry,
		&resolverBindingRepository{},
		&resolverOAuthStarter{url: "https://provider.example/authorize"},
	)
}

func newChannelPlatformTestResolverWithDependencies(
	dashboardID uuid.UUID,
	channel channelsmodel.Channel,
	users map[uuid.UUID]usersmodel.User,
	registry *appplatform.Registry,
	bindings *resolverBindingRepository,
	oauth *resolverOAuthStarter,
) *Resolver {
	service := channelplatformservice.New(
		resolverChannelReader{channel: channel},
		resolverUserLookup{users: users},
		bindings,
		oauth,
		registry,
		resolverTransactionRunner{},
	)

	return &Resolver{deps: Deps{
		ChannelPlatformBindingsService: service,
		ChannelPlatformDashboard:       resolverDashboardGetter{dashboardID: dashboardID.String()},
	}}
}

type resolverChannelReader struct {
	channel channelsmodel.Channel
}

func (r resolverChannelReader) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return r.channel, nil
}

type resolverUserLookup struct {
	users map[uuid.UUID]usersmodel.User
}

func (r resolverUserLookup) GetByID(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
	user, ok := r.users[id]
	if !ok {
		return usersmodel.Nil, errors.New("user not found")
	}
	return user, nil
}

type resolverBindingRepository struct {
	binding   channelplatformsmodel.ChannelPlatform
	patch     channelplatformsrepo.PatchInput
	deletedID uuid.UUID
}

func (*resolverBindingRepository) LockByChannelID(context.Context, uuid.UUID) error {
	return nil
}

func (r *resolverBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	return r.binding, nil
}

func (r *resolverBindingRepository) Patch(_ context.Context, _ uuid.UUID, input channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
	r.patch = input
	updated := r.binding
	if input.Enabled != nil {
		updated.Enabled = *input.Enabled
	}
	return updated, nil
}

func (r *resolverBindingRepository) Delete(_ context.Context, id uuid.UUID) error {
	r.deletedID = id
	return nil
}

type resolverOAuthStarter struct {
	url        string
	channelID  uuid.UUID
	platform   platformentity.Platform
	redirectTo string
}

func (r *resolverOAuthStarter) StartPlatformAuthForChannel(_ context.Context, channelID uuid.UUID, platform platformentity.Platform, redirectTo string) (string, error) {
	r.channelID = channelID
	r.platform = platform
	r.redirectTo = redirectTo
	return r.url, nil
}

type resolverDashboardGetter struct {
	dashboardID string
}

func (r resolverDashboardGetter) GetSelectedDashboard(context.Context) (string, error) {
	return r.dashboardID, nil
}

type resolverTransactionRunner struct{}

func (resolverTransactionRunner) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type resolverPlatformProvider struct {
	platform platformentity.Platform
}

func (p resolverPlatformProvider) Name() string {
	return p.platform.String()
}

func (resolverPlatformProvider) GetAuthURL(string, string) string {
	return ""
}

func (resolverPlatformProvider) ExchangeCode(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
	return nil, nil
}

func (resolverPlatformProvider) RefreshToken(context.Context, appplatform.RefreshTokenInput) (*appplatform.PlatformTokens, error) {
	return nil, nil
}

func (resolverPlatformProvider) GetUser(context.Context, string) (*appplatform.PlatformUser, error) {
	return nil, nil
}

func testResolverPlatformRegistry(platforms ...platformentity.Platform) *appplatform.Registry {
	providers := make([]appplatform.PlatformProvider, 0, len(platforms))
	for _, platform := range platforms {
		providers = append(providers, resolverPlatformProvider{platform: platform})
	}

	return appplatform.NewRegistry(providers)
}
