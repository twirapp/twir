package bus_listener

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/scorfly/gokick"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	kickbotentity "github.com/twirapp/twir/libs/entities/kick_bot"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
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
			ID:                  uuid.New(),
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

func TestRequestUserToken_VKRefreshesAndPersistsRotatedTokens(t *testing.T) {
	impl, repo, refresher := newVKTokenTestImplementation(t, "new-refresh")

	response, err := impl.RequestUserToken(context.Background(), buscoretokens.GetUserTokenRequest{UserId: uuid.New()})
	if err != nil {
		t.Fatalf("request VK user token: %v", err)
	}
	if response.AccessToken != "new-access" {
		t.Fatalf("access token = %q, want new-access", response.AccessToken)
	}
	if refresher.calls != 1 || refresher.refreshToken != "old-refresh" || refresher.deviceID != "device-id" {
		t.Fatalf("unexpected VK refresh invocation: %#v", refresher)
	}

	persistedRefreshToken, err := crypto.Decrypt(*repo.lastUpdate.RefreshToken, impl.config.TokensCipherKey)
	if err != nil {
		t.Fatalf("decrypt persisted refresh token: %v", err)
	}
	if persistedRefreshToken != "new-refresh" {
		t.Fatalf("persisted refresh token = %q, want new-refresh", persistedRefreshToken)
	}
}

func TestRequestUserToken_VKPreservesRefreshTokenWhenProviderOmitsIt(t *testing.T) {
	impl, repo, _ := newVKTokenTestImplementation(t, "")

	if _, err := impl.RequestUserToken(context.Background(), buscoretokens.GetUserTokenRequest{UserId: uuid.New()}); err != nil {
		t.Fatalf("request VK user token: %v", err)
	}

	persistedRefreshToken, err := crypto.Decrypt(*repo.lastUpdate.RefreshToken, impl.config.TokensCipherKey)
	if err != nil {
		t.Fatalf("decrypt persisted refresh token: %v", err)
	}
	if persistedRefreshToken != "old-refresh" {
		t.Fatalf("persisted refresh token = %q, want old-refresh", persistedRefreshToken)
	}
}

func TestRequestUserToken_VKRejectsMissingDeviceIDBeforeRefresh(t *testing.T) {
	impl, repo, refresher := newVKTokenTestImplementation(t, "new-refresh")
	repo.userToken.DeviceID = nil

	_, err := impl.RequestUserToken(context.Background(), buscoretokens.GetUserTokenRequest{UserId: uuid.New()})
	if err == nil || !strings.Contains(err.Error(), "device ID") {
		t.Fatalf("expected missing device ID error, got %v", err)
	}
	if refresher.calls != 0 {
		t.Fatalf("VK refresh must not be called without a stored device ID, got %d calls", refresher.calls)
	}
	if repo.updateCalls != 0 {
		t.Fatalf("token must not be updated without a stored device ID, got %d updates", repo.updateCalls)
	}
}

func newVKTokenTestImplementation(t *testing.T, refreshToken string) (*tokensImpl, *fakeTokensRepository, *fakeVKTokenRefresher) {
	t.Helper()

	const cipherKey = "pnyfwfiulmnqlhkvixaeligpprcnlyke"
	oldAccessToken, err := crypto.Encrypt("old-access", cipherKey)
	if err != nil {
		t.Fatal(err)
	}
	oldRefreshToken, err := crypto.Encrypt("old-refresh", cipherKey)
	if err != nil {
		t.Fatal(err)
	}
	encryptedDeviceID, err := crypto.Encrypt("device-id", cipherKey)
	if err != nil {
		t.Fatal(err)
	}

	repo := &fakeTokensRepository{
		userToken: &tokenmodel.Token{
			ID:                  uuid.New(),
			AccessToken:         oldAccessToken,
			RefreshToken:        oldRefreshToken,
			DeviceID:            &encryptedDeviceID,
			ExpiresIn:           1,
			ObtainmentTimestamp: time.Now().UTC().Add(-time.Hour),
			Scopes:              []string{"vkid.personal_info"},
		},
	}
	refresher := &fakeVKTokenRefresher{
		response: &vk.IDToken{
			AccessToken:  "new-access",
			RefreshToken: refreshToken,
			ExpiresIn:    7200,
			Scopes:       []string{"vkid.personal_info"},
		},
	}

	return &tokensImpl{
		config:           cfg.Config{TokensCipherKey: cipherKey},
		log:              slog.New(slog.NewTextHandler(io.Discard, nil)),
		tokensRepository: repo,
		usersRepository: &fakeUsersRepository{
			user: usersmodel.User{Platform: platformentity.PlatformVKVideoLive},
		},
		newMutex: func(string) lockableMutex {
			return fakeMutex{}
		},
		newVKTokenRefresher: func() (vkTokenRefresher, error) {
			return refresher, nil
		},
	}, repo, refresher
}

type fakeMutex struct{}

func (fakeMutex) Lock() error           { return nil }
func (fakeMutex) Unlock() (bool, error) { return true, nil }

type fakeKickTokenRefresher struct {
	calls        int
	refreshToken string
	resp         gokick.TokenResponse
	err          error
}

type fakeVKTokenRefresher struct {
	calls        int
	refreshToken string
	deviceID     string
	response     *vk.IDToken
	err          error
}

func (f *fakeVKTokenRefresher) RefreshToken(ctx context.Context, input vk.IDRefreshTokenInput) (*vk.IDToken, error) {
	f.calls++
	f.refreshToken = input.RefreshToken
	f.deviceID = input.DeviceID
	return f.response, f.err
}

func (f *fakeKickTokenRefresher) RefreshToken(ctx context.Context, refreshToken string) (gokick.TokenResponse, error) {
	f.calls++
	f.refreshToken = refreshToken
	return f.resp, f.err
}

type fakeTokensRepository struct {
	botToken         *tokenmodel.Token
	userToken        *tokenmodel.Token
	getByBotIDCalls  int
	getByUserIDCalls int
	updateCalls      int
	lastUpdate       tokensrepository.UpdateTokenInput
}

func (f *fakeTokensRepository) GetByID(ctx context.Context, id uuid.UUID) (*tokenmodel.Token, error) {
	panic("unexpected call")
}

func (f *fakeTokensRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*tokenmodel.Token, error) {
	f.getByUserIDCalls++
	return f.userToken, nil
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
	f.lastUpdate = input
	if f.userToken != nil {
		if input.AccessToken != nil {
			f.userToken.AccessToken = *input.AccessToken
		}
		if input.RefreshToken != nil {
			f.userToken.RefreshToken = *input.RefreshToken
		}
		if input.ExpiresIn != nil {
			f.userToken.ExpiresIn = *input.ExpiresIn
		}
		if input.ObtainmentTimestamp != nil {
			f.userToken.ObtainmentTimestamp = *input.ObtainmentTimestamp
		}
		if len(input.Scopes) > 0 {
			f.userToken.Scopes = input.Scopes
		}
		return f.userToken, nil
	}

	return f.botToken, nil
}

type fakeKickBotsRepository struct {
	defaultBot      kickbotentity.KickBot
	getDefaultCalls int
	updateCalls     int
	updated         kickbotsrepository.UpdateTokenInput
}

type fakeUsersRepository struct {
	user usersmodel.User
}

func (f *fakeUsersRepository) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return f.user, nil
}

func (f *fakeUsersRepository) GetByPlatformID(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
	panic("unexpected call")
}

func (f *fakeUsersRepository) GetManyByIDS(context.Context, usersrepository.GetManyInput) ([]usersmodel.User, error) {
	panic("unexpected call")
}

func (f *fakeUsersRepository) Update(context.Context, uuid.UUID, usersrepository.UpdateInput) (usersmodel.User, error) {
	panic("unexpected call")
}

func (f *fakeUsersRepository) GetRandomOnlineUser(context.Context, usersrepository.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	panic("unexpected call")
}

func (f *fakeUsersRepository) GetOnlineUsersWithFilters(context.Context, usersrepository.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	panic("unexpected call")
}

func (f *fakeUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	panic("unexpected call")
}

func (f *fakeUsersRepository) Create(context.Context, usersrepository.CreateInput) (usersmodel.User, error) {
	panic("unexpected call")
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
var _ usersrepository.Repository = (*fakeUsersRepository)(nil)
