package scheduledvips

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/entities/platform"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type scheduledVipChannelLookupStub struct {
	channel channelsmodel.Channel
}

func (s scheduledVipChannelLookupStub) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return s.channel, nil
}

type scheduledVipRepositoryStub struct {
	scheduledvipsrepository.Repository
	vip     scheduledvipsentity.ScheduledVip
	deleted []uuid.UUID
	created []scheduledvipsrepository.CreateInput
}

func (s *scheduledVipRepositoryStub) GetByID(context.Context, uuid.UUID) (scheduledvipsentity.ScheduledVip, error) {
	return s.vip, nil
}

func (s *scheduledVipRepositoryStub) Delete(_ context.Context, id uuid.UUID) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *scheduledVipRepositoryStub) Create(_ context.Context, input scheduledvipsrepository.CreateInput) error {
	s.created = append(s.created, input)
	return nil
}

type scheduledVipUsersRepositoryStub struct {
	usersrepository.Repository
	user usersmodel.User
}

func (s scheduledVipUsersRepositoryStub) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return s.user, nil
}

func (s scheduledVipUsersRepositoryStub) GetByPlatformID(
	context.Context,
	platform.Platform,
	string,
) (usersmodel.User, error) {
	return s.user, nil
}

type scheduledVipCaptureTransport struct {
	calls  int
	method string
	path   string
	query  url.Values
}

func (t *scheduledVipCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.calls++
	t.method = req.Method
	t.path = req.URL.Path
	t.query = req.URL.Query()

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		Request:    req,
	}, nil
}

func newScheduledVipTestClient(t *testing.T, transport http.RoundTripper) *helix.Client {
	t.Helper()

	client, err := helix.NewClient(&helix.Options{
		ClientID: "test-client",
		HTTPClient: &http.Client{
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("new Helix client: %v", err)
	}

	return client
}

func scheduledVipBindings(ownerID uuid.UUID, twitchChannelID string) []channelplatformsmodel.ChannelPlatform {
	return []channelplatformsmodel.ChannelPlatform{
		{
			Platform:          platform.PlatformVKVideoLive,
			UserID:            uuid.New(),
			PlatformChannelID: "vk-channel",
			BotConfig:         json.RawMessage(`{`),
		},
		{
			Platform:          platform.PlatformKick,
			UserID:            uuid.New(),
			PlatformChannelID: "kick-channel",
			BotConfig:         json.RawMessage(`{`),
		},
		{
			Platform:          platform.PlatformTwitch,
			UserID:            ownerID,
			PlatformChannelID: twitchChannelID,
		},
	}
}

func TestRemoveUsesSelectedTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	vipID := uuid.New()
	vipUserID := uuid.New()
	ownerID := uuid.New()
	repository := &scheduledVipRepositoryStub{vip: scheduledvipsentity.ScheduledVip{
		ID:        vipID,
		ChannelID: channelID.String(),
		UserID:    vipUserID.String(),
	}}
	transport := &scheduledVipCaptureTransport{}
	var clientOwnerID uuid.UUID
	service := &Service{
		repo: repository,
		channelService: scheduledVipChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: scheduledVipBindings(ownerID, "twitch-channel"),
		}},
		usersRepo: scheduledVipUsersRepositoryStub{user: usersmodel.User{
			ID:         vipUserID,
			PlatformID: "vip-platform-user",
		}},
		newUserClient: func(userID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = userID
			return newScheduledVipTestClient(t, transport), nil
		},
	}

	if err := service.Remove(context.Background(), RemoveInput{
		ID:        vipID.String(),
		ChannelID: channelID.String(),
	}); err != nil {
		t.Fatalf("remove scheduled VIP: %v", err)
	}
	if len(repository.deleted) != 1 || repository.deleted[0] != vipID {
		t.Fatalf("deleted IDs = %#v, want [%s]", repository.deleted, vipID)
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.method != http.MethodDelete || transport.path != "/helix/channels/vips" {
		t.Fatalf("request = %s %s, want DELETE /helix/channels/vips", transport.method, transport.path)
	}
	if got := transport.query.Get("broadcaster_id"); got != "twitch-channel" {
		t.Fatalf("broadcaster ID = %q, want selected Twitch channel", got)
	}
	if got := transport.query.Get("user_id"); got != "vip-platform-user" {
		t.Fatalf("VIP user ID = %q, want platform user ID", got)
	}
}

func TestRemoveDeletesWhenSelectedTwitchBindingHasEmptyProviderID(t *testing.T) {
	channelID := uuid.New()
	vipID := uuid.New()
	vipUserID := uuid.New()
	ownerID := uuid.New()
	repository := &scheduledVipRepositoryStub{vip: scheduledvipsentity.ScheduledVip{
		ID:        vipID,
		ChannelID: channelID.String(),
		UserID:    vipUserID.String(),
	}}
	transport := &scheduledVipCaptureTransport{}
	service := &Service{
		repo: repository,
		channelService: scheduledVipChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: scheduledVipBindings(ownerID, ""),
		}},
		usersRepo: scheduledVipUsersRepositoryStub{user: usersmodel.User{
			ID:         vipUserID,
			PlatformID: "vip-platform-user",
		}},
		newUserClient: func(uuid.UUID) (*helix.Client, error) {
			return newScheduledVipTestClient(t, transport), nil
		},
	}

	if err := service.Remove(context.Background(), RemoveInput{
		ID:        vipID.String(),
		ChannelID: channelID.String(),
	}); err != nil {
		t.Fatalf("remove scheduled VIP: %v", err)
	}
	if len(repository.deleted) != 1 || repository.deleted[0] != vipID {
		t.Fatalf("deleted IDs = %#v, want [%s]", repository.deleted, vipID)
	}
	if transport.calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", transport.calls)
	}
}

func TestRemoveRejectsMissingTwitchBindingWithoutDeleting(t *testing.T) {
	channelID := uuid.New()
	vipID := uuid.New()
	vipUserID := uuid.New()
	repository := &scheduledVipRepositoryStub{vip: scheduledvipsentity.ScheduledVip{
		ID:        vipID,
		ChannelID: channelID.String(),
		UserID:    vipUserID.String(),
	}}
	clientCalls := 0
	service := &Service{
		repo: repository,
		channelService: scheduledVipChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: scheduledVipBindings(uuid.New(), "twitch-channel")[:2],
		}},
		newUserClient: func(uuid.UUID) (*helix.Client, error) {
			clientCalls++
			return nil, nil
		},
	}

	err := service.Remove(context.Background(), RemoveInput{
		ID:        vipID.String(),
		ChannelID: channelID.String(),
	})
	if err == nil {
		t.Fatal("expected missing Twitch binding error")
	}
	if len(repository.deleted) != 0 {
		t.Fatalf("deleted IDs = %#v, want none", repository.deleted)
	}
	if clientCalls != 0 {
		t.Fatalf("user client calls = %d, want 0", clientCalls)
	}
}

func TestCreateWithTwitchVipKeepsLegacyEmptyProviderIDBehavior(t *testing.T) {
	channelID := uuid.New()
	ownerID := uuid.New()
	targetUserID := uuid.New()
	repository := &scheduledVipRepositoryStub{}
	transport := &scheduledVipCaptureTransport{}
	var clientOwnerID uuid.UUID
	service := &Service{
		repo: repository,
		channelService: scheduledVipChannelLookupStub{channel: channelsmodel.Channel{
			Bindings: scheduledVipBindings(ownerID, ""),
		}},
		usersRepo: scheduledVipUsersRepositoryStub{user: usersmodel.User{ID: targetUserID}},
		newUserClient: func(userID uuid.UUID) (*helix.Client, error) {
			clientOwnerID = userID
			return newScheduledVipTestClient(t, transport), nil
		},
	}

	err := service.CreateWithTwitchVip(context.Background(), CreateWithTwitchVipInput{
		ChannelID: channelID.String(),
		UserID:    "target-platform-user",
	})
	if err != nil {
		t.Fatalf("create scheduled Twitch VIP: %v", err)
	}
	if clientOwnerID != ownerID {
		t.Fatalf("user client owner ID = %s, want %s", clientOwnerID, ownerID)
	}
	if transport.calls != 1 {
		t.Fatalf("HTTP calls = %d, want 1", transport.calls)
	}
	if len(repository.created) != 1 || repository.created[0].UserID != targetUserID.String() {
		t.Fatalf("created scheduled VIPs = %#v, want target database user", repository.created)
	}
}
