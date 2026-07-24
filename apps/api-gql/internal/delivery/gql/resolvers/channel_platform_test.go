package resolvers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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
		channelentity.Channel{
			ID: dashboardID,
			Bindings: []channelplatformentity.ChannelPlatform{
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
		channelentity.Channel{
			ID: dashboardID,
			Bindings: []channelplatformentity.ChannelPlatform{
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

func TestChannelPlatformOptionsReturnRegisteredPlatformsInDomainOrder(t *testing.T) {
	t.Parallel()

	resolver := newChannelPlatformTestResolver(
		uuid.New(),
		channelentity.Channel{},
		nil,
		testResolverPlatformRegistry(
			platformentity.PlatformVKVideoLive,
			platformentity.PlatformTwitch,
		),
	)

	got, err := (&queryResolver{resolver}).ChannelPlatformOptions(context.Background())
	if err != nil {
		t.Fatalf("ChannelPlatformOptions() error = %v", err)
	}
	if len(got) != 2 || got[0].Platform != gqlmodel.PlatformTwitch || got[1].Platform != gqlmodel.PlatformVkVideoLive {
		t.Fatalf("ChannelPlatformOptions() = %#v, want Twitch then VK", got)
	}
	if len(got[0].Capabilities) != len(platformentity.PlatformTwitch.Capabilities()) || len(got[1].Capabilities) != 0 {
		t.Fatalf("ChannelPlatformOptions() capabilities = %#v", got)
	}
}

func TestChannelPlatformBindingMutationsUseGenericOperations(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	kickUserID := uuid.New()
	kickBinding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: kickUserID, PlatformChannelID: "kick-channel", Enabled: true,
	}
	operations := &resolverChannelPlatformBindingsService{binding: kickBinding, connectURL: "https://provider.example/authorize"}
	resolver := newChannelPlatformTestResolverWithDependencies(
		dashboardID,
		channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{
			kickBinding,
			{ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		}},
		map[uuid.UUID]usersmodel.User{kickUserID: {ID: kickUserID, PlatformID: "kick-user", DisplayName: "Kick Name"}},
		testResolverPlatformRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		operations,
	)

	mutation := &mutationResolver{resolver}
	url, err := mutation.ChannelPlatformConnect(context.Background(), gqlmodel.PlatformKick)
	if err != nil {
		t.Fatalf("ChannelPlatformConnect() error = %v", err)
	}
	if url != operations.connectURL || operations.channelID != dashboardID || operations.platform != platformentity.PlatformKick {
		t.Fatalf("ChannelPlatformConnect() = %q, operations = %#v", url, operations)
	}

	updated, err := mutation.ChannelPlatformSetEnabled(context.Background(), gqlmodel.PlatformKick, false)
	if err != nil {
		t.Fatalf("ChannelPlatformSetEnabled() error = %v", err)
	}
	if updated == nil || updated.Enabled || operations.binding.Enabled {
		t.Fatalf("ChannelPlatformSetEnabled() = %#v, binding = %#v", updated, operations.binding)
	}

	disconnected, err := mutation.ChannelPlatformDisconnect(context.Background(), gqlmodel.PlatformKick)
	if err != nil {
		t.Fatalf("ChannelPlatformDisconnect() error = %v", err)
	}
	if !disconnected || operations.deletedID != kickBinding.ID {
		t.Fatalf("ChannelPlatformDisconnect() = %t, deleted = %s", disconnected, operations.deletedID)
	}
}

func newChannelPlatformTestResolver(
	dashboardID uuid.UUID,
	channel channelentity.Channel,
	users map[uuid.UUID]usersmodel.User,
	registry *appplatform.Registry,
) *Resolver {
	return newChannelPlatformTestResolverWithDependencies(
		dashboardID,
		channel,
		users,
		registry,
		nil,
	)
}

func newChannelPlatformTestResolverWithDependencies(
	dashboardID uuid.UUID,
	channel channelentity.Channel,
	users map[uuid.UUID]usersmodel.User,
	registry *appplatform.Registry,
	operations *resolverChannelPlatformBindingsService,
) *Resolver {
	if operations == nil {
		operations = &resolverChannelPlatformBindingsService{connectURL: "https://provider.example/authorize"}
	}
	operations.channel = channel
	operations.users = users
	operations.registry = registry

	return &Resolver{deps: Deps{
		ChannelPlatformBindingsService: operations,
		ChannelPlatformDashboard:       resolverDashboardGetter{dashboardID: dashboardID.String()},
	}}
}

type resolverChannelPlatformBindingsService struct {
	channel    channelentity.Channel
	users      map[uuid.UUID]usersmodel.User
	registry   *appplatform.Registry
	binding    channelplatformentity.ChannelPlatform
	connectURL string
	channelID  uuid.UUID
	platform   platformentity.Platform
	deletedID  uuid.UUID
}

func (s *resolverChannelPlatformBindingsService) List(_ context.Context, channelID uuid.UUID) ([]channelplatformservice.Binding, error) {
	if s.channel.ID != channelID {
		return nil, nil
	}

	bindings := make([]channelplatformservice.Binding, 0, len(s.channel.Bindings))
	for _, binding := range s.channel.Bindings {
		provider, ok := s.registry.Get(binding.Platform)
		if !ok || provider == nil {
			continue
		}
		bindings = append(bindings, channelplatformservice.Binding{Binding: binding, Profile: s.users[binding.UserID], Capabilities: binding.Platform.Capabilities()})
	}

	return bindings, nil
}

func (s *resolverChannelPlatformBindingsService) Options() []channelplatformservice.Option {
	options := make([]channelplatformservice.Option, 0)
	for _, platform := range platformentity.All() {
		provider, ok := s.registry.Get(platform)
		if ok && provider != nil {
			options = append(options, channelplatformservice.Option{Platform: platform, Capabilities: platform.Capabilities()})
		}
	}

	return options
}

func (s *resolverChannelPlatformBindingsService) Connect(_ context.Context, channelID uuid.UUID, platform platformentity.Platform) (string, error) {
	s.channelID = channelID
	s.platform = platform
	return s.connectURL, nil
}

func (s *resolverChannelPlatformBindingsService) Disconnect(_ context.Context, _ uuid.UUID, _ platformentity.Platform) error {
	s.deletedID = s.binding.ID
	return nil
}

func (s *resolverChannelPlatformBindingsService) SetEnabled(_ context.Context, _ uuid.UUID, _ platformentity.Platform, enabled bool) (channelplatformservice.Binding, error) {
	s.binding.Enabled = enabled
	return channelplatformservice.Binding{Binding: s.binding, Profile: s.users[s.binding.UserID], Capabilities: s.binding.Platform.Capabilities()}, nil
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

func (p resolverPlatformProvider) Platform() platformentity.Platform {
	return p.platform
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
