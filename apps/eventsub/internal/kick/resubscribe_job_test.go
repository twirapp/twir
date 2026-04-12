package kick

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	entity "github.com/twirapp/twir/libs/entities/user_platform_account"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
)

type mockSubManager struct {
	listResult        []SubscriptionInfo
	listErr           error
	subscribeErr      error
	subscribeAllCalls int
}

func (m *mockSubManager) ListSubscriptions(_ context.Context, _ string) ([]SubscriptionInfo, error) {
	return m.listResult, m.listErr
}

func (m *mockSubManager) SubscribeAll(_ context.Context, _ string, _ string) error {
	m.subscribeAllCalls++
	return m.subscribeErr
}

type mockPlatformAccountsRepo struct {
	accounts []entity.UserPlatformAccount
	err      error
}

func (m *mockPlatformAccountsRepo) GetAllByPlatform(_ context.Context, _ platform.Platform) ([]entity.UserPlatformAccount, error) {
	return m.accounts, m.err
}

func (m *mockPlatformAccountsRepo) GetByUserIDAndPlatform(_ context.Context, _ uuid.UUID, _ platform.Platform) (entity.UserPlatformAccount, error) {
	return entity.Nil, nil
}

func (m *mockPlatformAccountsRepo) GetAllByUserID(_ context.Context, _ uuid.UUID) ([]entity.UserPlatformAccount, error) {
	return nil, nil
}

func (m *mockPlatformAccountsRepo) GetByPlatformUserID(_ context.Context, _ platform.Platform, _ string) (entity.UserPlatformAccount, error) {
	return entity.Nil, nil
}

func (m *mockPlatformAccountsRepo) Upsert(_ context.Context, _ user_platform_accounts.UpsertInput) (entity.UserPlatformAccount, error) {
	return entity.Nil, nil
}

func (m *mockPlatformAccountsRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func TestResubscribeJob_MissingSubscriptions(t *testing.T) {
	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Type: "chat.message.sent"},
			{Type: "channel.follow"},
		},
	}

	accountsRepo := &mockPlatformAccountsRepo{
		accounts: []entity.UserPlatformAccount{
			{
				PlatformUserID: "12345",
				AccessToken:    "token-abc",
			},
		},
	}

	job := &ResubscribeJob{
		subManager:               subMgr,
		userPlatformAccountsRepo: accountsRepo,
		logger:                   slog.Default(),
		interval:                 23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 1 {
		t.Errorf("expected SubscribeAll called 1 time, got %d", subMgr.subscribeAllCalls)
	}
}

func TestResubscribeJob_AllPresent(t *testing.T) {
	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Type: "chat.message.sent"},
			{Type: "channel.follow"},
			{Type: "stream.online"},
			{Type: "stream.offline"},
		},
	}

	accountsRepo := &mockPlatformAccountsRepo{
		accounts: []entity.UserPlatformAccount{
			{
				PlatformUserID: "67890",
				AccessToken:    "token-xyz",
			},
		},
	}

	job := &ResubscribeJob{
		subManager:               subMgr,
		userPlatformAccountsRepo: accountsRepo,
		logger:                   slog.Default(),
		interval:                 23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called, got %d calls", subMgr.subscribeAllCalls)
	}
}

func TestResubscribeJob_ListSubscriptionsError(t *testing.T) {
	subMgr := &mockSubManager{
		listErr: errors.New("network error"),
	}

	accountsRepo := &mockPlatformAccountsRepo{
		accounts: []entity.UserPlatformAccount{
			{
				PlatformUserID: "11111",
				AccessToken:    "token-err",
			},
		},
	}

	job := &ResubscribeJob{
		subManager:               subMgr,
		userPlatformAccountsRepo: accountsRepo,
		logger:                   slog.Default(),
		interval:                 23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called on error, got %d calls", subMgr.subscribeAllCalls)
	}
}
