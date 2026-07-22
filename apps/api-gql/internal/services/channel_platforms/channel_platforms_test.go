package channel_platforms

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestListIncludesRegisteredBindingsProfilesAndCapabilities(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	vkUserID := uuid.New()
	service := &Service{
		channels: fakeChannelReader{channel: channelsmodel.Channel{
			ID: channelID,
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: twitchUserID, PlatformChannelID: "twitch-channel", Enabled: true},
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: kickUserID, PlatformChannelID: "kick-channel", Enabled: false},
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true},
			},
		}},
		users: fakeUserLookup{users: map[uuid.UUID]usersmodel.User{
			twitchUserID: {ID: twitchUserID, PlatformID: "twitch-user", Login: "twitch-login", DisplayName: "Twitch Name", Avatar: "https://example.com/twitch.png"},
			kickUserID:   {ID: kickUserID, PlatformID: "kick-user", Login: "kick-login", DisplayName: "Kick Name"},
			vkUserID:     {ID: vkUserID, PlatformID: "vk-user", Login: "", DisplayName: "VK Name", Avatar: "https://example.com/vk.png"},
		}},
		registry: testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick, platformentity.PlatformVKVideoLive),
	}

	got, err := service.List(context.Background(), channelID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("List() returned %d bindings, want 3", len(got))
	}
	if got[0].Profile.DisplayName != "Twitch Name" || !got[0].Binding.Enabled {
		t.Fatalf("first binding = %#v, want Twitch profile and enabled state", got[0])
	}
	if got[1].Profile.DisplayName != "Kick Name" || got[1].Binding.Enabled {
		t.Fatalf("second binding = %#v, want Kick profile and disabled state", got[1])
	}
	if got[2].Profile.DisplayName != "VK Name" || !got[2].Binding.Enabled {
		t.Fatalf("third binding = %#v, want VK profile and enabled state", got[2])
	}
	if !reflect.DeepEqual(got[0].Capabilities, platformentity.PlatformTwitch.Capabilities()) {
		t.Fatalf("Twitch capabilities = %#v, want domain capabilities %#v", got[0].Capabilities, platformentity.PlatformTwitch.Capabilities())
	}
}

func TestListOmitsBindingsForUnavailablePlatforms(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	vkUserID := uuid.New()
	service := &Service{
		channels: fakeChannelReader{channel: channelsmodel.Channel{
			ID: channelID,
			Bindings: []channelplatformsmodel.ChannelPlatform{{
				ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true,
			}},
		}},
		users:    fakeUserLookup{users: map[uuid.UUID]usersmodel.User{vkUserID: {ID: vkUserID}}},
		registry: testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
	}

	got, err := service.List(context.Background(), channelID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("List() = %#v, want unavailable VK binding omitted", got)
	}
}

func TestBindingOperationsUseGenericDependencies(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	userID := uuid.New()
	binding := channelplatformsmodel.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: userID, PlatformChannelID: "kick-channel", Enabled: true,
	}
	repository := &fakeBindingRepository{binding: binding}
	oauth := &fakeOAuthStarter{url: "https://provider.example/authorize"}
	service := &Service{
		channels: fakeChannelReader{channel: channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
			binding,
			{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel"},
		}}},
		users:    fakeUserLookup{users: map[uuid.UUID]usersmodel.User{userID: {ID: userID, PlatformID: "kick-user"}}},
		bindings: repository,
		oauth:    oauth,
		registry: testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
	}

	url, err := service.Connect(context.Background(), channelID, platformentity.PlatformKick, "/dashboard")
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	if url != oauth.url || oauth.channelID != channelID || oauth.platform != platformentity.PlatformKick || oauth.redirectTo != "/dashboard" {
		t.Fatalf("Connect() = %q, starter = %#v", url, oauth)
	}

	updated, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, false)
	if err != nil {
		t.Fatalf("SetEnabled() error = %v", err)
	}
	if updated.Binding.Enabled || repository.patch.Enabled == nil || *repository.patch.Enabled {
		t.Fatalf("SetEnabled() = %#v, patch = %#v, want disabled binding", updated, repository.patch)
	}

	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformKick); err != nil {
		t.Fatalf("Disconnect() error = %v", err)
	}
	if repository.deletedID != binding.ID {
		t.Fatalf("Disconnect() deleted %s, want %s", repository.deletedID, binding.ID)
	}
}

func TestDisconnectRejectsLastAvailableBindingWhenDisabledPlatformIsHidden(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchBinding := channelplatformsmodel.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true,
	}
	repository := &fakeBindingRepository{binding: twitchBinding}
	service := &Service{
		channels: fakeChannelReader{channel: channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
			twitchBinding,
			{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: uuid.New(), PlatformChannelID: "vk-channel", Enabled: false},
		}}},
		bindings: repository,
		registry: testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
	}

	err := service.Disconnect(context.Background(), channelID, platformentity.PlatformTwitch)
	if !errors.Is(err, ErrLastBinding) {
		t.Fatalf("Disconnect() error = %v, want ErrLastBinding", err)
	}
	if repository.deletedID != uuid.Nil {
		t.Fatalf("Disconnect() deleted hidden-only fallback binding %s", repository.deletedID)
	}
}

func TestBindingOperationsRejectUnavailablePlatform(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	service := &Service{
		channels: fakeChannelReader{channel: channelsmodel.Channel{ID: channelID}},
		registry: testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
	}

	if _, err := service.Connect(context.Background(), channelID, platformentity.PlatformVKVideoLive, "/dashboard"); !errors.Is(err, ErrPlatformUnavailable) {
		t.Fatalf("Connect() error = %v, want ErrPlatformUnavailable", err)
	}
	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformVKVideoLive); !errors.Is(err, ErrPlatformUnavailable) {
		t.Fatalf("Disconnect() error = %v, want ErrPlatformUnavailable", err)
	}
	if _, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformVKVideoLive, false); !errors.Is(err, ErrPlatformUnavailable) {
		t.Fatalf("SetEnabled() error = %v, want ErrPlatformUnavailable", err)
	}
}

type fakeChannelReader struct {
	channel channelsmodel.Channel
	err     error
}

func (f fakeChannelReader) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return f.channel, f.err
}

type fakeUserLookup struct {
	users map[uuid.UUID]usersmodel.User
}

func (f fakeUserLookup) GetByID(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
	user, ok := f.users[id]
	if !ok {
		return usersmodel.Nil, errors.New("user not found")
	}
	return user, nil
}

type fakeBindingRepository struct {
	binding   channelplatformsmodel.ChannelPlatform
	patch     channelplatformsrepo.PatchInput
	deletedID uuid.UUID
}

func (f *fakeBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	return f.binding, nil
}

func (f *fakeBindingRepository) Patch(_ context.Context, _ uuid.UUID, input channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
	f.patch = input
	updated := f.binding
	if input.Enabled != nil {
		updated.Enabled = *input.Enabled
	}
	return updated, nil
}

func (f *fakeBindingRepository) Delete(_ context.Context, id uuid.UUID) error {
	f.deletedID = id
	return nil
}

type fakeOAuthStarter struct {
	url        string
	channelID  uuid.UUID
	platform   platformentity.Platform
	redirectTo string
}

func (f *fakeOAuthStarter) StartPlatformAuthForChannel(_ context.Context, channelID uuid.UUID, platform platformentity.Platform, redirectTo string) (string, error) {
	f.channelID = channelID
	f.platform = platform
	f.redirectTo = redirectTo
	return f.url, nil
}

type testPlatformProvider struct {
	platform platformentity.Platform
}

func (p testPlatformProvider) Name() string {
	return p.platform.String()
}

func (testPlatformProvider) GetAuthURL(string, string) string {
	return ""
}

func (testPlatformProvider) ExchangeCode(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
	return nil, nil
}

func (testPlatformProvider) RefreshToken(context.Context, appplatform.RefreshTokenInput) (*appplatform.PlatformTokens, error) {
	return nil, nil
}

func (testPlatformProvider) GetUser(context.Context, string) (*appplatform.PlatformUser, error) {
	return nil, nil
}

func testRegistry(platforms ...platformentity.Platform) *appplatform.Registry {
	providers := make([]appplatform.PlatformProvider, 0, len(platforms))
	for _, platform := range platforms {
		providers = append(providers, testPlatformProvider{platform: platform})
	}
	return appplatform.NewRegistry(providers)
}
