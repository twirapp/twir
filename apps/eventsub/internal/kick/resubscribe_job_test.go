package kick

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	entity "github.com/twirapp/twir/libs/entities/kick_bot"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	kick_bots "github.com/twirapp/twir/libs/repositories/kick_bots"
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

type mockKickBotsRepo struct {
	bot entity.KickBot
	err error
}

func (m *mockKickBotsRepo) GetDefault(_ context.Context) (entity.KickBot, error) {
	return entity.Nil, nil
}

func (m *mockKickBotsRepo) GetByID(_ context.Context, _ uuid.UUID) (entity.KickBot, error) {
	return m.bot, m.err
}

func (m *mockKickBotsRepo) GetByKickUserID(_ context.Context, _ uuid.UUID) (entity.KickBot, error) {
	return m.bot, m.err
}

func (m *mockKickBotsRepo) Create(_ context.Context, _ kick_bots.CreateInput) (entity.KickBot, error) {
	return entity.Nil, nil
}

func (m *mockKickBotsRepo) Upsert(_ context.Context, _ kick_bots.UpsertInput) (entity.KickBot, error) {
	return entity.Nil, nil
}

func (m *mockKickBotsRepo) UpdateToken(_ context.Context, _ uuid.UUID, _ kick_bots.UpdateTokenInput) (entity.KickBot, error) {
	return entity.Nil, nil
}

const testCipherKey = "pnyfwfiulmnqlhkvixaeligpprcnlyke"

func mustEncrypt(t *testing.T, plaintext string) string {
	t.Helper()
	enc, err := crypto.Encrypt(plaintext, testCipherKey)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}
	return enc
}

func TestResubscribeJob_MissingSubscriptions(t *testing.T) {
	kickUserID := uuid.New()
	kickBotID := uuid.New()

	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Type: "chat.message.sent"},
			{Type: "channel.follow"},
		},
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:         uuid.New(),
				KickUserID: &kickUserID,
				KickBotID:  &kickBotID,
			},
		},
	}

	botsRepo := &mockKickBotsRepo{
		bot: entity.KickBot{
			ID:          kickBotID.String(),
			KickUserID:  kickUserID,
			AccessToken: mustEncrypt(t, "token-abc"),
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		kickBotsRepo: botsRepo,
		logger:       slog.Default(),
		config:       cfg.Config{TokensCipherKey: testCipherKey},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 1 {
		t.Errorf("expected SubscribeAll called 1 time, got %d", subMgr.subscribeAllCalls)
	}
}

func TestResubscribeJob_AllPresent(t *testing.T) {
	kickUserID := uuid.New()
	kickBotID := uuid.New()

	subMgr := &mockSubManager{
		listResult: []SubscriptionInfo{
			{Type: "chat.message.sent"},
			{Type: "channel.follow"},
			{Type: "stream.online"},
			{Type: "stream.offline"},
		},
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:         uuid.New(),
				KickUserID: &kickUserID,
				KickBotID:  &kickBotID,
			},
		},
	}

	botsRepo := &mockKickBotsRepo{
		bot: entity.KickBot{
			ID:          kickBotID.String(),
			KickUserID:  kickUserID,
			AccessToken: mustEncrypt(t, "token-xyz"),
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		kickBotsRepo: botsRepo,
		logger:       slog.Default(),
		config:       cfg.Config{TokensCipherKey: testCipherKey},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called, got %d calls", subMgr.subscribeAllCalls)
	}
}

func TestResubscribeJob_ListSubscriptionsError(t *testing.T) {
	kickUserID := uuid.New()
	kickBotID := uuid.New()

	subMgr := &mockSubManager{
		listErr: errors.New("network error"),
	}

	chRepo := &mockChannelsRepo{
		channels: []channelsmodel.Channel{
			{
				ID:         uuid.New(),
				KickUserID: &kickUserID,
				KickBotID:  &kickBotID,
			},
		},
	}

	botsRepo := &mockKickBotsRepo{
		bot: entity.KickBot{
			ID:          kickBotID.String(),
			KickUserID:  kickUserID,
			AccessToken: mustEncrypt(t, "token-err"),
		},
	}

	job := &ResubscribeJob{
		subManager:   subMgr,
		channelsRepo: chRepo,
		kickBotsRepo: botsRepo,
		logger:       slog.Default(),
		config:       cfg.Config{TokensCipherKey: testCipherKey},
		interval:     23 * time.Hour,
	}

	job.run(context.Background())

	if subMgr.subscribeAllCalls != 0 {
		t.Errorf("expected SubscribeAll not called on error, got %d calls", subMgr.subscribeAllCalls)
	}
}
