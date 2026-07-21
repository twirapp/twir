package kick

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type mockSubManager struct {
	listResult         []SubscriptionInfo
	listErr            error
	subscribeErr       error
	subscribeAllCalls  int
	subscribeAllIDs    []uuid.UUID
	listBroadcasterIDs []int
}

func (m *mockSubManager) ListSubscriptions(_ context.Context, broadcasterUserID int) ([]SubscriptionInfo, error) {
	m.listBroadcasterIDs = append(m.listBroadcasterIDs, broadcasterUserID)
	return m.listResult, m.listErr
}

func (m *mockSubManager) SubscribeAll(_ context.Context, channelID uuid.UUID) error {
	m.subscribeAllCalls++
	m.subscribeAllIDs = append(m.subscribeAllIDs, channelID)
	return m.subscribeErr
}

func testKickBinding(userID uuid.UUID, enabled bool) channelplatformsmodel.ChannelPlatform {
	return testKickBindingWithPlatformChannelID(userID, "12345", enabled)
}

func testKickBindingWithPlatformChannelID(
	userID uuid.UUID,
	platformChannelID string,
	enabled bool,
) channelplatformsmodel.ChannelPlatform {
	return channelplatformsmodel.ChannelPlatform{
		Platform:          platform.PlatformKick,
		UserID:            userID,
		PlatformChannelID: platformChannelID,
		Enabled:           enabled,
	}
}

func TestResubscribeJob_MissingSubscriptions(t *testing.T) {
	kickUserID := uuid.New()

	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Event: "chat.message.sent"},
			{Event: "channel.followed"},
			{Event: "channel.subscription.new"},
			{Event: "channel.subscription.renewal"},
			{Event: "channel.subscription.gifts"},
			{Event: "channel.reward.redemption.updated"},
			{Event: "livestream.status.updated"},
			{Event: "moderation.banned"},
		},
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID: uuid.New(),
				Bindings: []channelplatformsmodel.ChannelPlatform{
					{
						Platform:          platform.PlatformTwitch,
						UserID:            uuid.New(),
						PlatformChannelID: "twitch-channel",
						Enabled:           true,
					},
					testKickBindingWithPlatformChannelID(kickUserID, "98765", true),
				},
			},
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		logger:       slog.Default(),
		config:       cfg.Config{},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 1 {
		t.Errorf("expected SubscribeAll called 1 time, got %d", subMgr.subscribeAllCalls)
	}
	if len(subMgr.subscribeAllIDs) != 1 || subMgr.subscribeAllIDs[0] != kickUserID {
		t.Fatalf("SubscribeAll IDs = %v, want [%s]", subMgr.subscribeAllIDs, kickUserID)
	}
	if len(subMgr.listBroadcasterIDs) != 1 || subMgr.listBroadcasterIDs[0] != 98765 {
		t.Fatalf("ListSubscriptions IDs = %v, want [98765]", subMgr.listBroadcasterIDs)
	}
}

func TestResubscribeJobProcessesAllKickBindingChannels(t *testing.T) {
	const channelCount = 11

	allKickChannels := make([]channelsmodel.Channel, 0, channelCount)
	legacyPage := make([]channelsmodel.Channel, 0, channelCount-1)
	for i := range channelCount {
		channel := channelsmodel.Channel{
			ID: uuid.New(),
			Bindings: []channelplatformsmodel.ChannelPlatform{
				testKickBindingWithPlatformChannelID(
					uuid.New(),
					strconv.Itoa(10000+i),
					true,
				),
			},
		}
		allKickChannels = append(allKickChannels, channel)
		if i < channelCount-1 {
			legacyPage = append(legacyPage, channel)
		}
	}

	subMgr := &mockSubManager{}
	chRepo := &mockChannelsRepo{
		channels:                legacyPage,
		bindingPlatformChannels: allKickChannels,
	}
	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		logger:       slog.Default(),
		config:       cfg.Config{},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != channelCount {
		t.Fatalf("SubscribeAll calls = %d, want %d", subMgr.subscribeAllCalls, channelCount)
	}
	if len(subMgr.listBroadcasterIDs) != channelCount {
		t.Fatalf("ListSubscriptions calls = %d, want %d", len(subMgr.listBroadcasterIDs), channelCount)
	}
	for i, broadcasterID := range subMgr.listBroadcasterIDs {
		want := 10000 + i
		if broadcasterID != want {
			t.Fatalf("ListSubscriptions ID at %d = %d, want %d", i, broadcasterID, want)
		}
	}
	if chRepo.bindingPlatformLookup != platform.PlatformKick {
		t.Fatalf("binding platform lookup = %q, want %q", chRepo.bindingPlatformLookup, platform.PlatformKick)
	}
}

func TestResubscribeJob_AllPresent(t *testing.T) {
	kickUserID := uuid.New()

	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Event: "chat.message.sent"},
			{Event: "channel.followed"},
			{Event: "channel.subscription.new"},
			{Event: "channel.subscription.renewal"},
			{Event: "channel.subscription.gifts"},
			{Event: "channel.reward.redemption.updated"},
			{Event: "livestream.status.updated"},
			{Event: "livestream.metadata.updated"},
			{Event: "moderation.banned"},
		},
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:       uuid.New(),
				Bindings: []channelplatformsmodel.ChannelPlatform{testKickBinding(kickUserID, true)},
			},
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		logger:       slog.Default(),
		config:       cfg.Config{},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called, got %d calls", subMgr.subscribeAllCalls)
	}
}

func TestResubscribeJob_ListSubscriptionsError(t *testing.T) {
	kickUserID := uuid.New()

	subMgr := &mockSubManager{
		listErr: errors.New("network error"),
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:       uuid.New(),
				Bindings: []channelplatformsmodel.ChannelPlatform{testKickBinding(kickUserID, true)},
			},
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		logger:       slog.Default(),
		config:       cfg.Config{},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called on error, got %d calls", subMgr.subscribeAllCalls)
	}
}
