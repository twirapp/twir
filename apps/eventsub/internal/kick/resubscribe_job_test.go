package kick

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type mockSubManager struct {
	listResult        []SubscriptionInfo
	listErr           error
	subscribeErr      error
	subscribeAllCalls int
}

func (m *mockSubManager) ListSubscriptions(_ context.Context, _ int) ([]SubscriptionInfo, error) {
	return m.listResult, m.listErr
}

func (m *mockSubManager) SubscribeAll(_ context.Context, _ string) error {
	m.subscribeAllCalls++
	return m.subscribeErr
}

func TestResubscribeJob_MissingSubscriptions(t *testing.T) {
	kickUserID := uuid.New()

	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Event: "chat.message.sent"},
			{Event: "channel.followed"},
		},
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:         uuid.New(),
				KickUserID: &kickUserID,
				IsEnabled:  true,
			},
		},
	}

	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         kickUserID.String(),
			PlatformID: "12345",
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		usersRepo:    usersRepo,
		logger:       slog.Default(),
		config:       cfg.Config{},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 1 {
		t.Errorf("expected SubscribeAll called 1 time, got %d", subMgr.subscribeAllCalls)
	}
}

func TestResubscribeJob_AllPresent(t *testing.T) {
	kickUserID := uuid.New()

	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Event: "chat.message.sent"},
			{Event: "channel.followed"},
			{Event: "livestream.status.updated"},
		},
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:         uuid.New(),
				KickUserID: &kickUserID,
				IsEnabled:  true,
			},
		},
	}

	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         kickUserID.String(),
			PlatformID: "12345",
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		usersRepo:    usersRepo,
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
				ID:         uuid.New(),
				KickUserID: &kickUserID,
				IsEnabled:  true,
			},
		},
	}

	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         kickUserID.String(),
			PlatformID: "12345",
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		usersRepo:    usersRepo,
		logger:       slog.Default(),
		config:       cfg.Config{},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called on error, got %d calls", subMgr.subscribeAllCalls)
	}
}
