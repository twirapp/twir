package channel_platforms

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func TestListIncludesRegisteredBindingsProfilesAndCapabilities(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	vkUserID := uuid.New()
	service := &Service{
		channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: twitchUserID, PlatformChannelID: "twitch-channel", Enabled: true},
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: kickUserID, PlatformChannelID: "kick-channel", Enabled: false},
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true},
			},
		}}),
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
		channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: vkUserID, PlatformChannelID: "vk-channel", Enabled: true,
			}},
		}}),
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

func TestOptionsReturnsRegisteredPlatformsInDomainOrder(t *testing.T) {
	t.Parallel()

	service := &Service{
		registry: testRegistry(
			platformentity.PlatformVKVideoLive,
			platformentity.PlatformTwitch,
		),
	}

	got := service.Options()
	want := []Option{
		{Platform: platformentity.PlatformTwitch, Capabilities: platformentity.PlatformTwitch.Capabilities()},
		{Platform: platformentity.PlatformVKVideoLive, Capabilities: platformentity.PlatformVKVideoLive.Capabilities()},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Options() = %#v, want %#v", got, want)
	}
}

func TestBindingOperationsUseRepositoryDependencies(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	userID := uuid.New()
	binding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: userID, PlatformChannelID: "kick-channel", Enabled: true,
	}
	repository := &fakeBindingRepository{binding: binding}
	service := &Service{
		channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{
			binding,
			{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel"},
		}}}),
		users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{userID: {ID: userID, PlatformID: "kick-user"}}},
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
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

func TestSetEnabledUsesTransactionContext(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	userID := uuid.New()
	binding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: userID, PlatformChannelID: "kick-channel", Enabled: true,
	}
	service := &Service{
		bindings:     &transactionAwareSetEnabledBindingRepository{fakeBindingRepository: fakeBindingRepository{binding: binding}},
		users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{userID: {ID: userID}}},
		registry:     testRegistry(platformentity.PlatformKick),
		transactions: &lifecycleTransactionRunner{},
	}

	updated, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, false)
	if err != nil {
		t.Fatalf("SetEnabled() error = %v", err)
	}
	if updated.Binding.Enabled {
		t.Fatalf("SetEnabled() = %#v, want disabled binding", updated)
	}
}

func TestServiceUsesRequiredDirectDependencies(t *testing.T) {
	serviceType := reflect.TypeFor[Service]()
	want := map[string]reflect.Type{
		"channels":     reflect.TypeFor[*channelservice.ChannelService](),
		"users":        reflect.TypeFor[usersrepo.Repository](),
		"bindings":     reflect.TypeFor[channelplatformsrepo.Repository](),
		"oauth":        reflect.TypeFor[*authroutes.Auth](),
		"transactions": reflect.TypeFor[trm.Manager](),
		"bus":          reflect.TypeFor[*buscore.Bus](),
	}

	for fieldName, wantType := range want {
		field, found := serviceType.FieldByName(fieldName)
		if !found {
			t.Fatalf("channel platform service is missing %s dependency", fieldName)
		}
		if field.Type != wantType {
			t.Fatalf("%s dependency type = %s, want %s", fieldName, field.Type, wantType)
		}
	}
}

func TestNewFxWiresDirectDependencies(t *testing.T) {
	channels := newTestChannelService(fakeChannelReader{})
	users := fakeUserLookup{}
	bindings := &fakeBindingRepository{}
	auth := &authroutes.Auth{}
	registry := testRegistry(platformentity.PlatformTwitch)
	transactions := &disconnectTransactionRunner{}
	bus := &buscore.Bus{}

	service := NewFx(Opts{
		ChannelService:       channels,
		UsersRepository:      users,
		ChannelPlatformsRepo: bindings,
		Auth:                 auth,
		PlatformRegistry:     registry,
		TrmManager:           transactions,
		TwirBus:              bus,
	})

	if service.channels != channels || !reflect.DeepEqual(service.users, users) || service.bindings != bindings || service.oauth != auth || service.registry != registry || service.transactions != transactions || service.bus != bus {
		t.Fatalf("NewFx() did not retain direct dependencies: %#v", service)
	}
}

func TestSetEnabledCompletesWhenEventSubIsUnavailable(t *testing.T) {
	t.Parallel()

	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("enabled_%t", enabled), func(t *testing.T) {
			channelID := uuid.New()
			binding := channelplatformentity.ChannelPlatform{
				ID:                uuid.New(),
				ChannelID:         channelID,
				Platform:          platformentity.PlatformKick,
				UserID:            uuid.New(),
				PlatformChannelID: "kick-channel",
				Enabled:           !enabled,
			}
			transaction := &lifecycleTransactionRunner{}
			service := &Service{
				bindings:     &fakeBindingRepository{binding: binding},
				users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{binding.UserID: {ID: binding.UserID}}},
				registry:     testRegistry(platformentity.PlatformKick),
				transactions: transaction,
				bus:          &buscore.Bus{},
			}

			updated, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, enabled)
			if err != nil {
				t.Fatalf("set enabled: %v", err)
			}
			if updated.Binding.Enabled != enabled || !transaction.committed {
				t.Fatalf("updated binding = %+v, committed = %t", updated.Binding, transaction.committed)
			}

		})
	}
}

func TestDisconnectCompletesWhenEventSubIsUnavailable(t *testing.T) {
	channelID := uuid.New()
	binding := channelplatformentity.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platformentity.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "kick-channel",
		Enabled:           true,
	}
	transaction := &lifecycleTransactionRunner{}
	repository := &fakeBindingRepository{binding: binding}
	service := &Service{
		channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{
			binding,
			{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		}}}),
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: transaction,
		bus:          &buscore.Bus{},
	}

	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformKick); err != nil {
		t.Fatalf("disconnect: %v", err)
	}
	if repository.deletedID != binding.ID || !transaction.committed {
		t.Fatalf("deleted = %s, committed = %t", repository.deletedID, transaction.committed)
	}
}

func TestBindingLifecycleReturnsTransactionFailure(t *testing.T) {
	commitErr := errors.New("commit failed")
	channelID := uuid.New()
	binding := channelplatformentity.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platformentity.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "kick-channel",
		Enabled:           true,
	}

	t.Run("set enabled", func(t *testing.T) {
		transaction := &lifecycleTransactionRunner{commitErr: commitErr}
		service := &Service{
			bindings:     &fakeBindingRepository{binding: binding},
			registry:     testRegistry(platformentity.PlatformKick),
			transactions: transaction,
			bus:          &buscore.Bus{},
		}
		if _, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, false); !errors.Is(err, commitErr) {
			t.Fatalf("set enabled error = %v, want %v", err, commitErr)
		}
	})

	t.Run("disconnect", func(t *testing.T) {
		transaction := &lifecycleTransactionRunner{commitErr: commitErr}
		service := &Service{
			channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{
				binding,
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
			}}}),
			bindings:     &fakeBindingRepository{binding: binding},
			registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
			transactions: transaction,
			bus:          &buscore.Bus{},
		}
		if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformKick); !errors.Is(err, commitErr) {
			t.Fatalf("disconnect error = %v, want %v", err, commitErr)
		}
	})
}

func TestSetEnabledSucceedsWhenEventSubIsUnavailable(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	userID := uuid.New()
	binding := channelplatformentity.ChannelPlatform{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: userID, Enabled: false}
	repository := &fakeBindingRepository{binding: binding}
	service := &Service{
		bindings:     repository,
		users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{userID: {ID: userID}}},
		registry:     testRegistry(platformentity.PlatformKick),
		transactions: &lifecycleTransactionRunner{},
		bus:          &buscore.Bus{},
	}

	updated, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, true)
	if err != nil {
		t.Fatalf("SetEnabled() error = %v", err)
	}
	if !updated.Binding.Enabled || repository.patch.Enabled == nil || !*repository.patch.Enabled {
		t.Fatalf("updated = %#v, patch = %#v", updated, repository.patch)
	}
}

func TestDisconnectRejectsLastAvailableBindingWhenDisabledPlatformIsHidden(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchBinding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true,
	}
	repository := &fakeBindingRepository{binding: twitchBinding}
	service := &Service{
		channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{
			twitchBinding,
			{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformVKVideoLive, UserID: uuid.New(), PlatformChannelID: "vk-channel", Enabled: false},
		}}}),
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
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
		channels: newTestChannelService(fakeChannelReader{channel: channelentity.Channel{ID: channelID}}),
		registry: testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
	}

	if _, err := service.Connect(context.Background(), channelID, platformentity.PlatformVKVideoLive); !errors.Is(err, ErrPlatformUnavailable) {
		t.Fatalf("Connect() error = %v, want ErrPlatformUnavailable", err)
	}
	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformVKVideoLive); !errors.Is(err, ErrPlatformUnavailable) {
		t.Fatalf("Disconnect() error = %v, want ErrPlatformUnavailable", err)
	}
	if _, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformVKVideoLive, false); !errors.Is(err, ErrPlatformUnavailable) {
		t.Fatalf("SetEnabled() error = %v, want ErrPlatformUnavailable", err)
	}
}

func TestDisconnectLocksBindingsBeforeReloadAndDelete(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchBinding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true,
	}
	events := make([]string, 0, 4)
	repository := &orderedDisconnectBindingRepository{binding: twitchBinding, events: &events}
	service := &Service{
		channels: newTestChannelService(orderedDisconnectChannelReader{
			channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{
				twitchBinding,
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
			}},
			events: &events,
		}),
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
	}

	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformTwitch); err != nil {
		t.Fatalf("Disconnect() error = %v", err)
	}

	if want := []string{"lock", "read", "get", "delete"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("Disconnect() operations = %#v, want %#v", events, want)
	}
}

func TestDisconnectStopsWhenBindingLockFails(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	lockErr := errors.New("binding lock failed")
	events := make([]string, 0, 4)
	repository := &orderedDisconnectBindingRepository{lockErr: lockErr, events: &events}
	service := &Service{
		channels:     newTestChannelService(orderedDisconnectChannelReader{events: &events}),
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
	}

	err := service.Disconnect(context.Background(), channelID, platformentity.PlatformTwitch)
	if !errors.Is(err, lockErr) {
		t.Fatalf("Disconnect() error = %v, want lock error", err)
	}
	if want := []string{"lock"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("Disconnect() operations = %#v, want %#v", events, want)
	}
}

func TestDisconnectUsesTransactionContextForLockReloadAndDelete(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	twitchBinding := channelplatformentity.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true,
	}
	events := make([]string, 0, 4)
	service := &Service{
		channels: newTestChannelService(transactionAwareChannelReader{
			channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{
				twitchBinding,
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
			}},
			events: &events,
		}),
		bindings:     &transactionAwareBindingRepository{binding: twitchBinding, events: &events},
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
	}

	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformTwitch); err != nil {
		t.Fatalf("Disconnect() error = %v", err)
	}
	if want := []string{"lock", "read", "get", "delete"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("Disconnect() operations = %#v, want %#v", events, want)
	}
}

func TestDisconnectSerializedOperationsKeepOneAvailableBinding(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	state := &serializedDisconnectState{channelID: channelID, bindings: []channelplatformentity.ChannelPlatform{
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
	}}
	repository := &serializedDisconnectBindingRepository{state: state}
	service := &Service{
		channels:     newTestChannelService(serializedDisconnectChannelReader{state: state}),
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
	}

	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformTwitch); err != nil {
		t.Fatalf("first Disconnect() error = %v", err)
	}
	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformKick); !errors.Is(err, ErrLastBinding) {
		t.Fatalf("second Disconnect() error = %v, want ErrLastBinding", err)
	}
	if repository.lockCalls != 2 || repository.deleteCalls != 1 || len(state.bindings) != 1 {
		t.Fatalf("serialized disconnect state = locks %d, deletes %d, bindings %#v", repository.lockCalls, repository.deleteCalls, state.bindings)
	}
}

func newTestChannelService(repository channelsrepo.Repository) *channelservice.ChannelService {
	return channelservice.NewChannelService(repository, nil, config.Config{}, nil, nil)
}

type fakeChannelReader struct {
	channelsrepo.Repository
	channel channelentity.Channel
	err     error
}

type orderedDisconnectChannelReader struct {
	channelsrepo.Repository
	channel channelentity.Channel
	events  *[]string
}

func (r orderedDisconnectChannelReader) GetByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	*r.events = append(*r.events, "read")
	return r.channel, nil
}

func (f fakeChannelReader) GetByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	return f.channel, f.err
}

type fakeUserLookup struct {
	usersrepo.Repository
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
	channelplatformsrepo.Repository
	binding   channelplatformentity.ChannelPlatform
	patch     channelplatformsrepo.PatchInput
	deletedID uuid.UUID
}

type lifecycleTransactionContextKey struct{}

type lifecycleTransactionRunner struct {
	committed bool
	commitErr error
}

func (r *lifecycleTransactionRunner) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := fn(context.WithValue(ctx, lifecycleTransactionContextKey{}, true)); err != nil {
		return err
	}
	if r.commitErr != nil {
		return r.commitErr
	}
	r.committed = true
	return nil
}

func (r *lifecycleTransactionRunner) DoWithSettings(
	ctx context.Context,
	_ trm.Settings,
	fn func(context.Context) error,
) error {
	return r.Do(ctx, fn)
}

type transactionAwareSetEnabledBindingRepository struct {
	fakeBindingRepository
}

func (r *transactionAwareSetEnabledBindingRepository) GetByChannelAndPlatform(ctx context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformentity.ChannelPlatform, error) {
	if ctx.Value(lifecycleTransactionContextKey{}) != true {
		return channelplatformentity.Nil, errors.New("binding lookup did not receive transaction context")
	}
	return r.binding, nil
}

func (r *transactionAwareSetEnabledBindingRepository) Patch(ctx context.Context, _ uuid.UUID, input channelplatformsrepo.PatchInput) (channelplatformentity.ChannelPlatform, error) {
	if ctx.Value(lifecycleTransactionContextKey{}) != true {
		return channelplatformentity.Nil, errors.New("binding patch did not receive transaction context")
	}
	r.patch = input
	updated := r.binding
	if input.Enabled != nil {
		updated.Enabled = *input.Enabled
	}
	return updated, nil
}

func (*fakeBindingRepository) LockByChannelID(context.Context, uuid.UUID) error {
	return nil
}

type orderedDisconnectBindingRepository struct {
	channelplatformsrepo.Repository
	binding channelplatformentity.ChannelPlatform
	lockErr error
	events  *[]string
}

func (r *orderedDisconnectBindingRepository) LockByChannelID(context.Context, uuid.UUID) error {
	*r.events = append(*r.events, "lock")
	return r.lockErr
}

func (r *orderedDisconnectBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformentity.ChannelPlatform, error) {
	*r.events = append(*r.events, "get")
	return r.binding, nil
}

func (r *orderedDisconnectBindingRepository) Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformentity.ChannelPlatform, error) {
	return r.binding, nil
}

func (r *orderedDisconnectBindingRepository) Delete(context.Context, uuid.UUID) error {
	*r.events = append(*r.events, "delete")
	return nil
}

type disconnectTransactionContextKey struct{}

type disconnectTransactionRunner struct{}

func (*disconnectTransactionRunner) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(context.WithValue(ctx, disconnectTransactionContextKey{}, true))
}

func (r *disconnectTransactionRunner) DoWithSettings(
	ctx context.Context,
	_ trm.Settings,
	fn func(context.Context) error,
) error {
	return r.Do(ctx, fn)
}

type transactionAwareChannelReader struct {
	channelsrepo.Repository
	channel channelentity.Channel
	events  *[]string
}

func (r transactionAwareChannelReader) GetByID(ctx context.Context, _ uuid.UUID) (channelentity.Channel, error) {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return channelentity.Nil, errors.New("channel reload did not receive transaction context")
	}
	*r.events = append(*r.events, "read")
	return r.channel, nil
}

type transactionAwareBindingRepository struct {
	channelplatformsrepo.Repository
	binding channelplatformentity.ChannelPlatform
	events  *[]string
}

func (r *transactionAwareBindingRepository) LockByChannelID(ctx context.Context, _ uuid.UUID) error {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return errors.New("lock did not receive transaction context")
	}
	*r.events = append(*r.events, "lock")
	return nil
}

func (r *transactionAwareBindingRepository) GetByChannelAndPlatform(ctx context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformentity.ChannelPlatform, error) {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return channelplatformentity.Nil, errors.New("binding reload did not receive transaction context")
	}
	*r.events = append(*r.events, "get")
	return r.binding, nil
}

func (r *transactionAwareBindingRepository) Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformentity.ChannelPlatform, error) {
	return r.binding, nil
}

func (r *transactionAwareBindingRepository) Delete(ctx context.Context, _ uuid.UUID) error {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return errors.New("binding delete did not receive transaction context")
	}
	*r.events = append(*r.events, "delete")
	return nil
}

type serializedDisconnectState struct {
	channelID uuid.UUID
	bindings  []channelplatformentity.ChannelPlatform
}

type serializedDisconnectChannelReader struct {
	channelsrepo.Repository
	state *serializedDisconnectState
}

func (r serializedDisconnectChannelReader) GetByID(_ context.Context, _ uuid.UUID) (channelentity.Channel, error) {
	bindings := append([]channelplatformentity.ChannelPlatform(nil), r.state.bindings...)
	return channelentity.Channel{ID: r.state.channelID, Bindings: bindings}, nil
}

type serializedDisconnectBindingRepository struct {
	channelplatformsrepo.Repository
	state       *serializedDisconnectState
	lockCalls   int
	deleteCalls int
}

func (r *serializedDisconnectBindingRepository) LockByChannelID(context.Context, uuid.UUID) error {
	r.lockCalls++
	return nil
}

func (r *serializedDisconnectBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, platform platformentity.Platform) (channelplatformentity.ChannelPlatform, error) {
	for _, binding := range r.state.bindings {
		if binding.Platform == platform {
			return binding, nil
		}
	}
	return channelplatformentity.Nil, channelplatformsrepo.ErrNotFound
}

func (r *serializedDisconnectBindingRepository) Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformentity.ChannelPlatform, error) {
	return channelplatformentity.Nil, errors.New("unexpected Patch call")
}

func (r *serializedDisconnectBindingRepository) Delete(_ context.Context, id uuid.UUID) error {
	r.deleteCalls++
	for index, binding := range r.state.bindings {
		if binding.ID == id {
			r.state.bindings = append(r.state.bindings[:index], r.state.bindings[index+1:]...)
			return nil
		}
	}
	return channelplatformsrepo.ErrNotFound
}

func (f *fakeBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformentity.ChannelPlatform, error) {
	return f.binding, nil
}

func (f *fakeBindingRepository) Patch(_ context.Context, _ uuid.UUID, input channelplatformsrepo.PatchInput) (channelplatformentity.ChannelPlatform, error) {
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
	url       string
	channelID uuid.UUID
	platform  platformentity.Platform
}

func (f *fakeOAuthStarter) StartPlatformAuthForChannel(_ context.Context, channelID uuid.UUID, platform platformentity.Platform) (string, error) {
	f.channelID = channelID
	f.platform = platform
	return f.url, nil
}

type testPlatformProvider struct {
	platform platformentity.Platform
}

func (p testPlatformProvider) Platform() platformentity.Platform {
	return p.platform
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
