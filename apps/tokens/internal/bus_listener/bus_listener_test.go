package bus_listener

import (
	"context"
	"errors"
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
	vkvideobotentity "github.com/twirapp/twir/libs/entities/vk_video_bot"
	"github.com/twirapp/twir/libs/integrations/vk"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	vkvideobotsrepository "github.com/twirapp/twir/libs/repositories/vk_video_bots"
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

func TestRequestBotToken_VKVideoLiveReturnsFreshSingletonToken(t *testing.T) {
	t.Parallel()

	impl, repository, refresher, runner := newVKVideoBotTokenTestImplementation(t, 3600, "new-refresh")

	response, err := impl.RequestBotToken(
		context.Background(),
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive, BotId: "ignored"},
	)
	if err != nil {
		t.Fatalf("request VK Video Live bot token: %v", err)
	}
	if response.AccessToken != "old-access" {
		t.Fatalf("access token = %q, want old-access", response.AccessToken)
	}
	if !reflect.DeepEqual(response.Scopes, []string{"chat:write"}) {
		t.Fatalf("scopes = %#v, want chat:write", response.Scopes)
	}
	if response.ExpiresIn != 3600 {
		t.Fatalf("expires in = %d, want 3600", response.ExpiresIn)
	}
	if repository.getCalls != 1 || repository.updateCalls != 0 {
		t.Fatalf("repository calls = get:%d update:%d, want get:1 update:0", repository.getCalls, repository.updateCalls)
	}
	if refresher.calls != 0 {
		t.Fatalf("refresh calls = %d, want 0", refresher.calls)
	}
	if runner.doCalls != 1 || !repository.callsWithinTransaction {
		t.Fatalf("VK singleton access must be performed in one transaction")
	}
}

func TestRequestBotToken_VKVideoLiveRefreshesAndPersistsRotatedRefreshToken(t *testing.T) {
	t.Parallel()

	impl, repository, refresher, _ := newVKVideoBotTokenTestImplementation(t, 1, "rotated-refresh")

	response, err := impl.RequestBotToken(
		context.Background(),
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive},
	)
	if err != nil {
		t.Fatalf("request VK Video Live bot token: %v", err)
	}
	if response.AccessToken != "new-access" {
		t.Fatalf("access token = %q, want new-access", response.AccessToken)
	}
	if refresher.calls != 1 || refresher.refreshToken != "old-refresh" {
		t.Fatalf("unexpected refresh invocation: %#v", refresher)
	}
	if repository.updateCalls != 1 {
		t.Fatalf("update calls = %d, want 1", repository.updateCalls)
	}
	persistedRefreshToken, err := crypto.Decrypt(repository.updated.EncryptedRefreshToken, impl.config.TokensCipherKey)
	if err != nil {
		t.Fatalf("decrypt persisted refresh token: %v", err)
	}
	if persistedRefreshToken != "rotated-refresh" {
		t.Fatalf("persisted refresh token = %q, want rotated-refresh", persistedRefreshToken)
	}
}

func TestRequestBotToken_VKVideoLivePreservesRefreshTokenWhenProviderOmitsIt(t *testing.T) {
	t.Parallel()

	impl, repository, _, _ := newVKVideoBotTokenTestImplementation(t, 1, "")

	if _, err := impl.RequestBotToken(
		context.Background(),
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive},
	); err != nil {
		t.Fatalf("request VK Video Live bot token: %v", err)
	}

	persistedRefreshToken, err := crypto.Decrypt(repository.updated.EncryptedRefreshToken, impl.config.TokensCipherKey)
	if err != nil {
		t.Fatalf("decrypt persisted refresh token: %v", err)
	}
	if persistedRefreshToken != "old-refresh" {
		t.Fatalf("persisted refresh token = %q, want old-refresh", persistedRefreshToken)
	}
}

func TestRequestBotToken_VKVideoLiveReportsMissingSingleton(t *testing.T) {
	t.Parallel()

	impl, repository, _, _ := newVKVideoBotTokenTestImplementation(t, 3600, "new-refresh")
	repository.getErr = vkvideobotsrepository.ErrNotFound

	_, err := impl.RequestBotToken(
		context.Background(),
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive},
	)
	if !errors.Is(err, vkvideobotsrepository.ErrNotFound) {
		t.Fatalf("error = %v, want wrapped ErrNotFound", err)
	}
	if strings.Contains(err.Error(), "old-access") || strings.Contains(err.Error(), "old-refresh") {
		t.Fatalf("missing singleton error must not contain token data: %v", err)
	}
}

func TestRequestBotToken_VKVideoLiveBoundsEncryptionFailures(t *testing.T) {
	t.Parallel()

	impl, _, refresher, _ := newVKVideoBotTokenTestImplementation(t, 1, "new-refresh")
	refresher.onRefresh = func() {
		impl.config.TokensCipherKey = "invalid"
	}

	_, err := impl.RequestBotToken(
		context.Background(),
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive},
	)
	if err == nil {
		t.Fatal("expected encryption error")
	}
	if strings.Contains(err.Error(), "new-access") || strings.Contains(err.Error(), "old-refresh") {
		t.Fatalf("encryption error must not contain token data: %v", err)
	}
}

func TestRequestBotToken_VKVideoLiveSerializesRefreshBeforeSingletonReplacement(t *testing.T) {
	t.Parallel()

	impl, repository, refresher, runner := newVKVideoBotTokenTestImplementation(t, 1, "new-refresh")
	replacement := repository.bot
	replacement.EncryptedAccessToken = mustEncrypt(t, "replacement-access", impl.config.TokensCipherKey)
	replacement.EncryptedRefreshToken = mustEncrypt(t, "replacement-refresh", impl.config.TokensCipherKey)
	replacement.VKUserID = uuid.New()
	refresher.onRefresh = func() {
		runner.afterDo = func() {
			repository.bot = replacement
		}
	}

	response, err := impl.RequestBotToken(
		context.Background(),
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformVKVideoLive},
	)
	if err != nil {
		t.Fatalf("request VK Video Live bot token: %v", err)
	}
	if response.AccessToken != "new-access" {
		t.Fatalf("access token = %q, want new-access", response.AccessToken)
	}
	if !repository.callsWithinTransaction {
		t.Fatal("singleton get and update must remain inside the transaction")
	}
	persistedAccessToken, err := crypto.Decrypt(repository.bot.EncryptedAccessToken, impl.config.TokensCipherKey)
	if err != nil {
		t.Fatalf("decrypt replacement access token: %v", err)
	}
	if persistedAccessToken != "replacement-access" {
		t.Fatalf("stale refresh overwrote replacement with %q", persistedAccessToken)
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
	if refresher.calls != 1 || refresher.refreshToken != "old-refresh" {
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

func TestRequestUserToken_VKRefreshesWithoutStoredDeviceID(t *testing.T) {
	impl, repo, refresher := newVKTokenTestImplementation(t, "new-refresh")
	repo.userToken.DeviceID = nil

	if _, err := impl.RequestUserToken(context.Background(), buscoretokens.GetUserTokenRequest{UserId: uuid.New()}); err != nil {
		t.Fatalf("VK refresh must not require a stored device ID: %v", err)
	}
	if refresher.calls != 1 {
		t.Fatalf("expected one VK refresh, got %d", refresher.calls)
	}
	if repo.updateCalls != 1 {
		t.Fatalf("expected token to be updated once, got %d updates", repo.updateCalls)
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
			Scopes:              []string{"user_info"},
		},
	}
	refresher := &fakeVKTokenRefresher{
		response: &vk.OAuthToken{
			AccessToken:  "new-access",
			RefreshToken: refreshToken,
			ExpiresIn:    7200,
			Scopes:       []string{"user_info"},
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

func newVKVideoBotTokenTestImplementation(
	t *testing.T,
	expiresIn int,
	refreshedRefreshToken string,
) (*tokensImpl, *fakeVKVideoBotsRepository, *fakeVKTokenRefresher, *fakeTransactionRunner) {
	t.Helper()

	const cipherKey = "pnyfwfiulmnqlhkvixaeligpprcnlyke"
	repository := &fakeVKVideoBotsRepository{
		callsWithinTransaction: true,
		bot: vkvideobotentity.VKVideoBot{
			ID:                    uuid.New(),
			EncryptedAccessToken:  mustEncrypt(t, "old-access", cipherKey),
			EncryptedRefreshToken: mustEncrypt(t, "old-refresh", cipherKey),
			Scopes:                []string{"chat:write"},
			ExpiresIn:             expiresIn,
			ObtainmentTimestamp:   time.Now().UTC().Add(-time.Hour),
			VKUserID:              uuid.New(),
		},
	}
	if expiresIn > 1 {
		repository.bot.ObtainmentTimestamp = time.Now().UTC()
	}
	runner := &fakeTransactionRunner{}
	repository.transactionRunner = runner
	refresher := &fakeVKTokenRefresher{
		response: &vk.OAuthToken{
			AccessToken:  "new-access",
			RefreshToken: refreshedRefreshToken,
			ExpiresIn:    7200,
			Scopes:       []string{"chat:write", "channel:read"},
		},
	}

	return &tokensImpl{
		config:          cfg.Config{TokensCipherKey: cipherKey},
		log:             slog.New(slog.NewTextHandler(io.Discard, nil)),
		vkVideoBotsRepo: repository,
		newMutex: func(string) lockableMutex {
			return fakeMutex{}
		},
		newVKBotTokenRefresher: func() (vkTokenRefresher, error) {
			return refresher, nil
		},
		transactionRunner: runner,
	}, repository, refresher, runner
}

func mustEncrypt(t *testing.T, value, key string) string {
	t.Helper()

	encryptedValue, err := crypto.Encrypt(value, key)
	if err != nil {
		t.Fatal(err)
	}

	return encryptedValue
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
	response     *vk.OAuthToken
	err          error
	onRefresh    func()
}

func (f *fakeVKTokenRefresher) RefreshToken(ctx context.Context, refreshToken string) (*vk.OAuthToken, error) {
	f.calls++
	f.refreshToken = refreshToken
	if f.onRefresh != nil {
		f.onRefresh()
	}
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

type fakeTransactionRunner struct {
	active  bool
	doCalls int
	afterDo func()
}

func (f *fakeTransactionRunner) Do(ctx context.Context, callback func(context.Context) error) error {
	f.doCalls++
	f.active = true
	err := callback(ctx)
	f.active = false
	if f.afterDo != nil {
		f.afterDo()
	}

	return err
}

type fakeVKVideoBotsRepository struct {
	bot                    vkvideobotentity.VKVideoBot
	getErr                 error
	updateErr              error
	getCalls               int
	updateCalls            int
	updated                vkvideobotsrepository.UpdateInput
	callsWithinTransaction bool
	transactionRunner      *fakeTransactionRunner
}

func (f *fakeVKVideoBotsRepository) Get(context.Context) (vkvideobotentity.VKVideoBot, error) {
	f.getCalls++
	f.callsWithinTransaction = f.callsWithinTransaction && f.transactionRunner.active
	if f.getErr != nil {
		return vkvideobotentity.Nil, f.getErr
	}

	return f.bot, nil
}

func (f *fakeVKVideoBotsRepository) Lock(context.Context) error {
	return nil
}

func (f *fakeVKVideoBotsRepository) Upsert(
	context.Context,
	vkvideobotsrepository.UpsertInput,
) (vkvideobotentity.VKVideoBot, error) {
	panic("unexpected call")
}

func (f *fakeVKVideoBotsRepository) Update(
	_ context.Context,
	input vkvideobotsrepository.UpdateInput,
) (vkvideobotentity.VKVideoBot, error) {
	f.updateCalls++
	f.callsWithinTransaction = f.callsWithinTransaction && f.transactionRunner.active
	f.updated = input
	if f.updateErr != nil {
		return vkvideobotentity.Nil, f.updateErr
	}

	f.bot.EncryptedAccessToken = input.EncryptedAccessToken
	f.bot.EncryptedRefreshToken = input.EncryptedRefreshToken
	f.bot.Scopes = input.Scopes
	f.bot.ExpiresIn = input.ExpiresIn
	f.bot.ObtainmentTimestamp = input.ObtainmentTimestamp
	f.bot.VKUserID = input.VKUserID
	return f.bot, nil
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
var _ vkvideobotsrepository.Repository = (*fakeVKVideoBotsRepository)(nil)
