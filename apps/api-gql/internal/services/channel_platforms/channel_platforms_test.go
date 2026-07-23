package channel_platforms

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
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
		users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{userID: {ID: userID, PlatformID: "kick-user"}}},
		bindings:     repository,
		oauth:        oauth,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: &disconnectTransactionRunner{},
	}

	url, err := service.Connect(context.Background(), channelID, platformentity.PlatformKick)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	if url != oauth.url || oauth.channelID != channelID || oauth.platform != platformentity.PlatformKick {
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

func TestSetEnabledUsesTransactionContext(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	userID := uuid.New()
	binding := channelplatformsmodel.ChannelPlatform{
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

func TestServiceExposesEventSubPublisher(t *testing.T) {
	serviceType := reflect.TypeOf(Service{})
	field, found := serviceType.FieldByName("eventSub")
	if !found {
		t.Fatal("channel platform service is missing eventSub")
	}
	if field.Type.Kind() != reflect.Interface {
		t.Fatalf("eventSub type = %s, want interface", field.Type)
	}
	for _, method := range []string{"Subscribe", "Unsubscribe"} {
		if _, found := field.Type.MethodByName(method); !found {
			t.Fatalf("eventSub publisher is missing %s", method)
		}
	}
}

func TestSetEnabledPublishesLifecycleAfterCommit(t *testing.T) {
	t.Parallel()

	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("enabled_%t", enabled), func(t *testing.T) {
			channelID := uuid.New()
			binding := channelplatformsmodel.ChannelPlatform{
				ID:                uuid.New(),
				ChannelID:         channelID,
				Platform:          platformentity.PlatformKick,
				UserID:            uuid.New(),
				PlatformChannelID: "kick-channel",
				Enabled:           !enabled,
			}
			transaction := &lifecycleTransactionRunner{}
			publisher := &recordingEventSubPublisher{committed: &transaction.committed}
			service := &Service{
				bindings:     &fakeBindingRepository{binding: binding},
				users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{binding.UserID: {ID: binding.UserID}}},
				registry:     testRegistry(platformentity.PlatformKick),
				transactions: transaction,
				eventSub:     publisher,
			}

			updated, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, enabled)
			if err != nil {
				t.Fatalf("set enabled: %v", err)
			}
			if updated.Binding.Enabled != enabled || !transaction.committed {
				t.Fatalf("updated binding = %+v, committed = %t", updated.Binding, transaction.committed)
			}

			if enabled {
				want := eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID.String(), Platform: platformentity.PlatformKick}
				if !reflect.DeepEqual(publisher.subscribeRequests, []eventsub.EventsubSubscribeToAllEventsRequest{want}) || len(publisher.unsubscribeRequests) != 0 {
					t.Fatalf("subscribe = %#v, unsubscribe = %#v", publisher.subscribeRequests, publisher.unsubscribeRequests)
				}
				return
			}

			want := eventsub.EventsubUnsubscribeRequest{
				ChannelID: channelID.String(),
				Platform:  platformentity.PlatformKick,
				Binding: &eventsub.EventsubBindingSnapshot{
					ID:                binding.ID.String(),
					UserID:            binding.UserID.String(),
					PlatformChannelID: binding.PlatformChannelID,
				},
			}
			if !reflect.DeepEqual(publisher.unsubscribeRequests, []eventsub.EventsubUnsubscribeRequest{want}) || len(publisher.subscribeRequests) != 0 {
				t.Fatalf("unsubscribe = %#v, subscribe = %#v", publisher.unsubscribeRequests, publisher.subscribeRequests)
			}
		})
	}
}

func TestDisconnectPublishesBindingSnapshotAfterCommit(t *testing.T) {
	channelID := uuid.New()
	binding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platformentity.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "kick-channel",
		Enabled:           true,
	}
	transaction := &lifecycleTransactionRunner{}
	publisher := &recordingEventSubPublisher{committed: &transaction.committed}
	repository := &fakeBindingRepository{binding: binding}
	service := &Service{
		channels: fakeChannelReader{channel: channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
			binding,
			{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		}}},
		bindings:     repository,
		registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		transactions: transaction,
		eventSub:     publisher,
	}

	if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformKick); err != nil {
		t.Fatalf("disconnect: %v", err)
	}
	want := eventsub.EventsubUnsubscribeRequest{
		ChannelID: channelID.String(),
		Platform:  platformentity.PlatformKick,
		Binding: &eventsub.EventsubBindingSnapshot{
			ID:                binding.ID.String(),
			UserID:            binding.UserID.String(),
			PlatformChannelID: binding.PlatformChannelID,
		},
	}
	if repository.deletedID != binding.ID || !transaction.committed || !reflect.DeepEqual(publisher.unsubscribeRequests, []eventsub.EventsubUnsubscribeRequest{want}) {
		t.Fatalf("deleted = %s, committed = %t, unsubscribe = %#v", repository.deletedID, transaction.committed, publisher.unsubscribeRequests)
	}
}

func TestBindingLifecycleDoesNotPublishWhenTransactionFails(t *testing.T) {
	commitErr := errors.New("commit failed")
	channelID := uuid.New()
	binding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platformentity.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "kick-channel",
		Enabled:           true,
	}

	t.Run("set enabled", func(t *testing.T) {
		transaction := &lifecycleTransactionRunner{commitErr: commitErr}
		publisher := &recordingEventSubPublisher{committed: &transaction.committed}
		service := &Service{
			bindings:     &fakeBindingRepository{binding: binding},
			registry:     testRegistry(platformentity.PlatformKick),
			transactions: transaction,
			eventSub:     publisher,
		}
		if _, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, false); !errors.Is(err, commitErr) {
			t.Fatalf("set enabled error = %v, want %v", err, commitErr)
		}
		if len(publisher.subscribeRequests) != 0 || len(publisher.unsubscribeRequests) != 0 {
			t.Fatalf("published requests = subscribe %#v unsubscribe %#v", publisher.subscribeRequests, publisher.unsubscribeRequests)
		}
	})

	t.Run("disconnect", func(t *testing.T) {
		transaction := &lifecycleTransactionRunner{commitErr: commitErr}
		publisher := &recordingEventSubPublisher{committed: &transaction.committed}
		service := &Service{
			channels: fakeChannelReader{channel: channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
				binding,
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
			}}},
			bindings:     &fakeBindingRepository{binding: binding},
			registry:     testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
			transactions: transaction,
			eventSub:     publisher,
		}
		if err := service.Disconnect(context.Background(), channelID, platformentity.PlatformKick); !errors.Is(err, commitErr) {
			t.Fatalf("disconnect error = %v, want %v", err, commitErr)
		}
		if len(publisher.subscribeRequests) != 0 || len(publisher.unsubscribeRequests) != 0 {
			t.Fatalf("published requests = subscribe %#v unsubscribe %#v", publisher.subscribeRequests, publisher.unsubscribeRequests)
		}
	})
}

func TestSetEnabledSucceedsAndLogsPostCommitPublishErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name           string
		enabled        bool
		subscribeErr   error
		unsubscribeErr error
	}{
		{name: "subscribe", enabled: true, subscribeErr: errors.New("subscribe publish failed")},
		{name: "unsubscribe", enabled: false, unsubscribeErr: errors.New("unsubscribe publish failed")},
	} {
		t.Run(test.name, func(t *testing.T) {
			channelID := uuid.New()
			userID := uuid.New()
			binding := channelplatformsmodel.ChannelPlatform{
				ID:                uuid.New(),
				ChannelID:         channelID,
				Platform:          platformentity.PlatformKick,
				UserID:            userID,
				PlatformChannelID: "kick-channel",
				Enabled:           !test.enabled,
			}
			repository := &fakeBindingRepository{binding: binding}
			transaction := &lifecycleTransactionRunner{}
			publisher := &recordingEventSubPublisher{
				committed:      &transaction.committed,
				subscribeErr:   test.subscribeErr,
				unsubscribeErr: test.unsubscribeErr,
			}
			var logs bytes.Buffer
			service := &Service{
				bindings:     repository,
				users:        fakeUserLookup{users: map[uuid.UUID]usersmodel.User{userID: {ID: userID}}},
				registry:     testRegistry(platformentity.PlatformKick),
				transactions: transaction,
				eventSub:     publisher,
				logger:       slog.New(slog.NewJSONHandler(&logs, nil)),
			}

			updated, err := service.SetEnabled(context.Background(), channelID, platformentity.PlatformKick, test.enabled)
			if err != nil {
				t.Fatalf("SetEnabled() error = %v", err)
			}
			if !transaction.committed || updated.Binding.Enabled != test.enabled || repository.patch.Enabled == nil || *repository.patch.Enabled != test.enabled {
				t.Fatalf("updated = %#v, patch = %#v, committed = %t", updated, repository.patch, transaction.committed)
			}

			if test.enabled {
				if len(publisher.subscribeRequests) != 1 || len(publisher.unsubscribeRequests) != 0 {
					t.Fatalf("subscribe = %#v, unsubscribe = %#v", publisher.subscribeRequests, publisher.unsubscribeRequests)
				}
			} else if len(publisher.unsubscribeRequests) != 1 || len(publisher.subscribeRequests) != 0 {
				t.Fatalf("unsubscribe = %#v, subscribe = %#v", publisher.unsubscribeRequests, publisher.subscribeRequests)
			}

			for _, want := range []string{
				`"level":"ERROR"`,
				fmt.Sprintf(`"msg":"cannot publish eventsub %s"`, test.name),
				fmt.Sprintf(`"channel_id":%q`, channelID.String()),
				fmt.Sprintf(`"platform":%q`, binding.Platform.String()),
			} {
				if !strings.Contains(logs.String(), want) {
					t.Fatalf("log = %q, want %q", logs.String(), want)
				}
			}
		})
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
		channels: fakeChannelReader{channel: channelsmodel.Channel{ID: channelID}},
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
	twitchBinding := channelplatformsmodel.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true,
	}
	events := make([]string, 0, 4)
	repository := &orderedDisconnectBindingRepository{binding: twitchBinding, events: &events}
	service := &Service{
		channels: orderedDisconnectChannelReader{
			channel: channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
				twitchBinding,
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
			}},
			events: &events,
		},
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
		channels:     orderedDisconnectChannelReader{events: &events},
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
	twitchBinding := channelplatformsmodel.ChannelPlatform{
		ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true,
	}
	events := make([]string, 0, 4)
	service := New(
		transactionAwareChannelReader{
			channel: channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{
				twitchBinding,
				{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
			}},
			events: &events,
		},
		nil,
		&transactionAwareBindingRepository{binding: twitchBinding, events: &events},
		nil,
		testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		&disconnectTransactionRunner{},
	)

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
	state := &serializedDisconnectState{channelID: channelID, bindings: []channelplatformsmodel.ChannelPlatform{
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: uuid.New(), PlatformChannelID: "twitch-channel", Enabled: true},
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
	}}
	repository := &serializedDisconnectBindingRepository{state: state}
	service := New(
		serializedDisconnectChannelReader{state: state},
		nil,
		repository,
		nil,
		testRegistry(platformentity.PlatformTwitch, platformentity.PlatformKick),
		&disconnectTransactionRunner{},
	)

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

type fakeChannelReader struct {
	channel channelsmodel.Channel
	err     error
}

type orderedDisconnectChannelReader struct {
	channel channelsmodel.Channel
	events  *[]string
}

func (r orderedDisconnectChannelReader) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	*r.events = append(*r.events, "read")
	return r.channel, nil
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

type recordingEventSubPublisher struct {
	committed           *bool
	subscribeErr        error
	unsubscribeErr      error
	subscribeRequests   []eventsub.EventsubSubscribeToAllEventsRequest
	unsubscribeRequests []eventsub.EventsubUnsubscribeRequest
}

func (p *recordingEventSubPublisher) Subscribe(_ context.Context, request eventsub.EventsubSubscribeToAllEventsRequest) error {
	if p.committed != nil && !*p.committed {
		return errors.New("subscribe published before transaction commit")
	}
	p.subscribeRequests = append(p.subscribeRequests, request)
	return p.subscribeErr
}

func (p *recordingEventSubPublisher) Unsubscribe(_ context.Context, request eventsub.EventsubUnsubscribeRequest) error {
	if p.committed != nil && !*p.committed {
		return errors.New("unsubscribe published before transaction commit")
	}
	p.unsubscribeRequests = append(p.unsubscribeRequests, request)
	return p.unsubscribeErr
}

type transactionAwareSetEnabledBindingRepository struct {
	fakeBindingRepository
}

func (r *transactionAwareSetEnabledBindingRepository) GetByChannelAndPlatform(ctx context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	if ctx.Value(lifecycleTransactionContextKey{}) != true {
		return channelplatformsmodel.Nil, errors.New("binding lookup did not receive transaction context")
	}
	return r.binding, nil
}

func (r *transactionAwareSetEnabledBindingRepository) Patch(ctx context.Context, _ uuid.UUID, input channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
	if ctx.Value(lifecycleTransactionContextKey{}) != true {
		return channelplatformsmodel.Nil, errors.New("binding patch did not receive transaction context")
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
	binding channelplatformsmodel.ChannelPlatform
	lockErr error
	events  *[]string
}

func (r *orderedDisconnectBindingRepository) LockByChannelID(context.Context, uuid.UUID) error {
	*r.events = append(*r.events, "lock")
	return r.lockErr
}

func (r *orderedDisconnectBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	*r.events = append(*r.events, "get")
	return r.binding, nil
}

func (r *orderedDisconnectBindingRepository) Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
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

type transactionAwareChannelReader struct {
	channel channelsmodel.Channel
	events  *[]string
}

func (r transactionAwareChannelReader) GetChannelByID(ctx context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return channelsmodel.Nil, errors.New("channel reload did not receive transaction context")
	}
	*r.events = append(*r.events, "read")
	return r.channel, nil
}

type transactionAwareBindingRepository struct {
	binding channelplatformsmodel.ChannelPlatform
	events  *[]string
}

func (r *transactionAwareBindingRepository) LockByChannelID(ctx context.Context, _ uuid.UUID) error {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return errors.New("lock did not receive transaction context")
	}
	*r.events = append(*r.events, "lock")
	return nil
}

func (r *transactionAwareBindingRepository) GetByChannelAndPlatform(ctx context.Context, _ uuid.UUID, _ platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	if ctx.Value(disconnectTransactionContextKey{}) != true {
		return channelplatformsmodel.Nil, errors.New("binding reload did not receive transaction context")
	}
	*r.events = append(*r.events, "get")
	return r.binding, nil
}

func (r *transactionAwareBindingRepository) Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
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
	bindings  []channelplatformsmodel.ChannelPlatform
}

type serializedDisconnectChannelReader struct {
	state *serializedDisconnectState
}

func (r serializedDisconnectChannelReader) GetChannelByID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	bindings := append([]channelplatformsmodel.ChannelPlatform(nil), r.state.bindings...)
	return channelsmodel.Channel{ID: r.state.channelID, Bindings: bindings}, nil
}

type serializedDisconnectBindingRepository struct {
	state       *serializedDisconnectState
	lockCalls   int
	deleteCalls int
}

func (r *serializedDisconnectBindingRepository) LockByChannelID(context.Context, uuid.UUID) error {
	r.lockCalls++
	return nil
}

func (r *serializedDisconnectBindingRepository) GetByChannelAndPlatform(_ context.Context, _ uuid.UUID, platform platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	for _, binding := range r.state.bindings {
		if binding.Platform == platform {
			return binding, nil
		}
	}
	return channelplatformsmodel.Nil, channelplatformsrepo.ErrNotFound
}

func (r *serializedDisconnectBindingRepository) Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, errors.New("unexpected Patch call")
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
