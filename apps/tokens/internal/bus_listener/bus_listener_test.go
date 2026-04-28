package bus_listener

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/scorfly/gokick"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	kickbotentity "github.com/twirapp/twir/libs/entities/kick_bot"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

func TestRequestBotToken_DefaultsToTwitchWhenPlatformEmpty(t *testing.T) {
	t.Parallel()

	const cipherKey = "pnyfwfiulmnqlhkvixaeligpprcnlyke"
	accessToken, err := crypto.Encrypt("twitch-access", cipherKey)
	if err != nil {
		t.Fatal(err)
	}
	refreshToken, err := crypto.Encrypt("twitch-refresh", cipherKey)
	if err != nil {
		t.Fatal(err)
	}

	repo := &fakeTokensRepository{
		botToken: &tokenmodel.Token{
			ID:                  uuid.New(),
			AccessToken:         accessToken,
			RefreshToken:        refreshToken,
			ExpiresIn:           3600,
			ObtainmentTimestamp: time.Now().UTC(),
			Scopes:              []string{"chat:edit"},
		},
	}

	impl := &tokensImpl{
		config:           cfg.Config{TokensCipherKey: cipherKey},
		tokensRepository: repo,
		kickBotsRepo:     &fakeKickBotsRepository{},
		newMutex: func(name string) lockableMutex {
			return fakeMutex{}
		},
	}

	resp, err := impl.RequestBotToken(context.Background(), buscoretokens.GetBotTokenRequest{BotId: "bot-1"})
	if err != nil {
		t.Fatalf("RequestBotToken returned error: %v", err)
	}

	if repo.getByBotIDCalls != 1 {
		t.Fatalf("expected GetByBotID to be called once, got %d", repo.getByBotIDCalls)
	}
	if resp.AccessToken != "twitch-access" {
		t.Fatalf("expected decrypted access token, got %q", resp.AccessToken)
	}
	if !reflect.DeepEqual(resp.Scopes, []string{"chat:edit"}) {
		t.Fatalf("unexpected scopes: %#v", resp.Scopes)
	}
	if resp.ExpiresIn != 3600 {
		t.Fatalf("unexpected expires_in: %d", resp.ExpiresIn)
	}
	if repo.updateCalls != 0 {
		t.Fatalf("expected no token update, got %d", repo.updateCalls)
	}
	if impl.kickBotsRepo.(*fakeKickBotsRepository).getDefaultCalls != 0 {
		t.Fatalf("expected no kick bot lookup")
	}
}

func TestRequestBotToken_KickRefreshesDefaultBot(t *testing.T) {
	t.Parallel()

	const cipherKey = "pnyfwfiulmnqlhkvixaeligpprcnlyke"
	accessToken, err := crypto.Encrypt("old-access", cipherKey)
	if err != nil {
		t.Fatal(err)
	}
	refreshToken, err := crypto.Encrypt("old-refresh", cipherKey)
	if err != nil {
		t.Fatal(err)
	}

	kickRepo := &fakeKickBotsRepository{
		defaultBot: kickbotentity.KickBot{
			ID:                  uuid.NewString(),
			AccessToken:         accessToken,
			RefreshToken:        refreshToken,
			Scopes:              []string{"chat:write"},
			ExpiresIn:           1,
			ObtainmentTimestamp: time.Now().UTC().Add(-time.Hour),
			KickUserLogin:       "kick-bot",
		},
	}

	refresher := &fakeKickTokenRefresher{
		resp: gokick.TokenResponse{
			AccessToken:  "new-access",
			RefreshToken: "new-refresh",
			ExpiresIn:    7200,
			Scope:        "chat:write channel:read",
		},
	}

	impl := &tokensImpl{
		config:       cfg.Config{TokensCipherKey: cipherKey},
		log:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		kickBotsRepo: kickRepo,
		newMutex: func(name string) lockableMutex {
			return fakeMutex{}
		},
		newKickTokenRefresher: func() (kickTokenRefresher, error) {
			return refresher, nil
		},
	}

	resp, err := impl.RequestBotToken(context.Background(), buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformKick})
	if err != nil {
		t.Fatalf("RequestBotToken returned error: %v", err)
	}

	if kickRepo.getDefaultCalls != 1 {
		t.Fatalf("expected GetDefault once, got %d", kickRepo.getDefaultCalls)
	}
	if refresher.calls != 1 {
		t.Fatalf("expected one kick refresh, got %d", refresher.calls)
	}
	if refresher.refreshToken != "old-refresh" {
		t.Fatalf("expected decrypted refresh token, got %q", refresher.refreshToken)
	}
	if kickRepo.updateCalls != 1 {
		t.Fatalf("expected UpdateToken once, got %d", kickRepo.updateCalls)
	}
	if resp.AccessToken != "new-access" {
		t.Fatalf("expected refreshed access token, got %q", resp.AccessToken)
	}
	if !reflect.DeepEqual(resp.Scopes, []string{"chat:write", "channel:read"}) {
		t.Fatalf("unexpected scopes: %#v", resp.Scopes)
	}
	if resp.ExpiresIn != 7200 {
		t.Fatalf("unexpected expires_in: %d", resp.ExpiresIn)
	}

	updatedAccessToken, err := crypto.Decrypt(kickRepo.updated.AccessToken, cipherKey)
	if err != nil {
		t.Fatal(err)
	}
	if updatedAccessToken != "new-access" {
		t.Fatalf("unexpected persisted access token: %q", updatedAccessToken)
	}
}

type fakeMutex struct{}

func (fakeMutex) Lock() error            { return nil }
func (fakeMutex) Unlock() (bool, error)  { return true, nil }

type fakeKickTokenRefresher struct {
	calls        int
	refreshToken string
	resp         gokick.TokenResponse
	err          error
}

func (f *fakeKickTokenRefresher) RefreshToken(ctx context.Context, refreshToken string) (gokick.TokenResponse, error) {
	f.calls++
	f.refreshToken = refreshToken
	return f.resp, f.err
}

type fakeTokensRepository struct {
	botToken         *tokenmodel.Token
	getByBotIDCalls  int
	updateCalls      int
}

func (f *fakeTokensRepository) GetByID(ctx context.Context, id uuid.UUID) (*tokenmodel.Token, error) {
	panic("unexpected call")
}

func (f *fakeTokensRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*tokenmodel.Token, error) {
	panic("unexpected call")
}

func (f *fakeTokensRepository) GetByBotID(ctx context.Context, botID string) (*tokenmodel.Token, error) {
	f.getByBotIDCalls++
	return f.botToken, nil
}

func (f *fakeTokensRepository) CreateUserToken(ctx context.Context, input tokensrepository.CreateInput) (*tokenmodel.Token, error) {
	panic("unexpected call")
}

func (f *fakeTokensRepository) UpdateTokenByID(ctx context.Context, id uuid.UUID, input tokensrepository.UpdateTokenInput) (*tokenmodel.Token, error) {
	f.updateCalls++
	return f.botToken, nil
}

type fakeKickBotsRepository struct {
	defaultBot       kickbotentity.KickBot
	getDefaultCalls  int
	updateCalls      int
	updated          kickbotsrepository.UpdateTokenInput
}

func (f *fakeKickBotsRepository) GetDefault(ctx context.Context) (kickbotentity.KickBot, error) {
	f.getDefaultCalls++
	return f.defaultBot, nil
}

func (f *fakeKickBotsRepository) GetByID(ctx context.Context, id uuid.UUID) (kickbotentity.KickBot, error) {
	panic("unexpected call")
}

func (f *fakeKickBotsRepository) GetByKickUserID(ctx context.Context, kickUserID uuid.UUID) (kickbotentity.KickBot, error) {
	panic("unexpected call")
}

func (f *fakeKickBotsRepository) Create(ctx context.Context, input kickbotsrepository.CreateInput) (kickbotentity.KickBot, error) {
	panic("unexpected call")
}

func (f *fakeKickBotsRepository) Upsert(ctx context.Context, input kickbotsrepository.UpsertInput) (kickbotentity.KickBot, error) {
	panic("unexpected call")
}

func (f *fakeKickBotsRepository) UpdateToken(ctx context.Context, id uuid.UUID, input kickbotsrepository.UpdateTokenInput) (kickbotentity.KickBot, error) {
	f.updateCalls++
	f.updated = input
	f.defaultBot.AccessToken = input.AccessToken
	f.defaultBot.RefreshToken = input.RefreshToken
	f.defaultBot.ExpiresIn = input.ExpiresIn
	f.defaultBot.ObtainmentTimestamp = input.ObtainmentTimestamp
	f.defaultBot.Scopes = input.Scopes
	return f.defaultBot, nil
}

var _ tokensrepository.Repository = (*fakeTokensRepository)(nil)
var _ kickbotsrepository.Repository = (*fakeKickBotsRepository)(nil)
