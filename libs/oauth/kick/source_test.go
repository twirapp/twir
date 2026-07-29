package kick

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/libs/crypto"
	entity "github.com/twirapp/twir/libs/entities/kick_bot"
	"github.com/twirapp/twir/libs/oauth"
	kickbots "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokens "github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

const testCipherKey = "0123456789abcdef0123456789abcdef"

func kickTestRedis(t *testing.T) *goredis.Client {
	t.Helper()
	address := os.Getenv("TWIR_KICK_TEST_REDIS_ADDR")
	if address == "" {
		t.Skip("TWIR_KICK_TEST_REDIS_ADDR is not set")
	}
	client := goredis.NewClient(&goredis.Options{Addr: address})
	t.Cleanup(func() { _ = client.Close() })
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("ping Redis: %v", err)
	}
	return client
}

func kickTestOptions(client *goredis.Client, refreshCalls *atomic.Int32, refreshResponse gokick.TokenResponse) SourceOptions {
	return SourceOptions{
		ClientID: "kick-client", ClientSecret: "kick-secret", Redis: client, CipherKey: testCipherKey,
		ClientFactory: func() (Client, error) {
			refreshCalls.Add(1)
			time.Sleep(50 * time.Millisecond)
			return fakeClient{response: refreshResponse}, nil
		},
	}
}

func seedEncryptedToken(t *testing.T, userID uuid.UUID, refreshToken string, obtainedAt time.Time) tokenmodel.Token {
	t.Helper()
	accessToken, err := crypto.Encrypt("kick-old-access", testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	encryptedRefresh, err := crypto.Encrypt(refreshToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	return tokenmodel.Token{
		ID: uuid.New(), AccessToken: accessToken, RefreshToken: encryptedRefresh,
		ExpiresIn: 3600, ObtainmentTimestamp: obtainedAt, Scopes: []string{"user:read"},
	}
}

func TestKickUserSourceRefreshesOnceAndPreservesOmittedRefreshToken(t *testing.T) {
	redisClient := kickTestRedis(t)
	userID := uuid.New()
	repository := newFakeTokensRepository(seedEncryptedToken(t, userID, "kick-original-refresh", time.Now().Add(-2*time.Hour)))
	var refreshCalls atomic.Int32
	options := kickTestOptions(redisClient, &refreshCalls, gokick.TokenResponse{
		AccessToken: "kick-rotated-access", ExpiresIn: 3600, Scope: "channel:read",
	})

	newSource := func() UserTokenSource {
		source, err := NewUserTokenSource(options, repository)
		if err != nil {
			t.Fatal(err)
		}
		return source
	}
	first, second := newSource(), newSource()

	var group sync.WaitGroup
	credentials := make(chan oauth.Credential, 2)
	errors := make(chan error, 2)
	for _, source := range []UserTokenSource{first, second} {
		group.Add(1)
		go func(source UserTokenSource) {
			defer group.Done()
			credential, err := source.Token(context.Background(), userID)
			credentials <- credential
			errors <- err
		}(source)
	}
	group.Wait()
	if err := <-errors; err != nil {
		t.Fatal(err)
	}
	if err := <-errors; err != nil {
		t.Fatal(err)
	}
	if got := refreshCalls.Load(); got != 1 {
		t.Fatalf("provider refreshes = %d, want 1", got)
	}
	firstCredential := <-credentials
	if firstCredential.AccessToken != "kick-rotated-access" {
		t.Fatalf("access token = %q", firstCredential.AccessToken)
	}
	if firstCredential.RefreshToken != "kick-original-refresh" {
		t.Fatalf("omitted refresh token not preserved: %q", firstCredential.RefreshToken)
	}

	updated := repository.updated
	if updated == nil || updated.AccessToken == nil || updated.RefreshToken == nil {
		t.Fatal("repository commit missing encrypted fields")
	}
	committedAccess, err := crypto.Decrypt(*updated.AccessToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	committedRefresh, err := crypto.Decrypt(*updated.RefreshToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	if committedAccess != "kick-rotated-access" || committedRefresh != "kick-original-refresh" {
		t.Fatalf("committed pair = %q/%q", committedAccess, committedRefresh)
	}
}

func TestKickBotSourceRotatesAndPersistsEncryptedPair(t *testing.T) {
	redisClient := kickTestRedis(t)
	botID := uuid.New()
	accessToken, err := crypto.Encrypt("kick-bot-old-access", testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	refreshToken, err := crypto.Encrypt("kick-bot-old-refresh", testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	repository := &fakeKickBotsRepository{bot: entity.KickBot{
		ID: botID, AccessToken: accessToken, RefreshToken: refreshToken,
		ExpiresIn: 3600, ObtainmentTimestamp: time.Now().Add(-2 * time.Hour), Scopes: []string{"chat:write"},
	}}
	var refreshCalls atomic.Int32
	options := kickTestOptions(redisClient, &refreshCalls, gokick.TokenResponse{
		AccessToken: "kick-bot-rotated-access", RefreshToken: "kick-bot-rotated-refresh", ExpiresIn: 3600, Scope: "chat:write moderation:write",
	})
	source, err := NewDefaultBotTokenSource(options, repository)
	if err != nil {
		t.Fatal(err)
	}

	credential, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if credential.AccessToken != "kick-bot-rotated-access" || credential.RefreshToken != "kick-bot-rotated-refresh" {
		t.Fatalf("credential = %q/%q", credential.AccessToken, credential.RefreshToken)
	}
	if len(credential.Scopes) != 2 || credential.Scopes[1] != "moderation:write" {
		t.Fatalf("scopes = %#v", credential.Scopes)
	}
	if repository.updated == nil {
		t.Fatal("bot repository commit missing")
	}
	committedAccess, err := crypto.Decrypt(repository.updated.AccessToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	committedRefresh, err := crypto.Decrypt(repository.updated.RefreshToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	if committedAccess != "kick-bot-rotated-access" || committedRefresh != "kick-bot-rotated-refresh" {
		t.Fatalf("committed bot pair = %q/%q", committedAccess, committedRefresh)
	}
	if got := refreshCalls.Load(); got != 1 {
		t.Fatalf("provider refreshes = %d, want 1", got)
	}
}

func TestKickAppTokenSourceFetchesOnceUnderConcurrency(t *testing.T) {
	redisClient := kickTestRedis(t)
	var fetchCalls atomic.Int32
	source, err := NewAppTokenSource(SourceOptions{
		ClientID: "kick-client", ClientSecret: "kick-secret", Redis: redisClient,
		AppClientFactory: func() (AppClient, error) {
			fetchCalls.Add(1)
			time.Sleep(50 * time.Millisecond)
			return fakeAppClient{response: gokick.AppTokenResponse{AccessToken: "kick-app-access", ExpiresIn: 3600}}, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var group sync.WaitGroup
	errors := make(chan error, 8)
	for range 8 {
		group.Add(1)
		go func() {
			defer group.Done()
			token, err := source.Token(context.Background())
			if err == nil && token.AccessToken != "kick-app-access" {
				err = oauth.ErrInvalidCredential
			}
			errors <- err
		}()
	}
	group.Wait()
	for range 8 {
		if err := <-errors; err != nil {
			t.Fatal(err)
		}
	}
	if got := fetchCalls.Load(); got != 1 {
		t.Fatalf("app token fetches = %d, want 1", got)
	}
}

type fakeAppClient struct{ response gokick.AppTokenResponse }

func (client fakeAppClient) GetAppAccessToken(context.Context) (gokick.AppTokenResponse, error) {
	return client.response, nil
}

type fakeTokensRepository struct {
	token   tokenmodel.Token
	userID  uuid.UUID
	updated *tokens.UpdateTokenInput
	mu      sync.Mutex
}

func newFakeTokensRepository(token tokenmodel.Token) *fakeTokensRepository {
	return &fakeTokensRepository{token: token, userID: uuid.New()}
}

func (repository *fakeTokensRepository) GetByID(context.Context, uuid.UUID) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	return &repository.token, nil
}

func (repository *fakeTokensRepository) GetByUserID(context.Context, uuid.UUID) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	return &repository.token, nil
}

func (repository *fakeTokensRepository) GetByBotID(context.Context, string) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	return &repository.token, nil
}

func (repository *fakeTokensRepository) CreateUserToken(context.Context, tokens.CreateInput) (*tokenmodel.Token, error) {
	return nil, tokens.ErrNotFound
}

func (repository *fakeTokensRepository) UpdateTokenByID(_ context.Context, _ uuid.UUID, input tokens.UpdateTokenInput) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	copied := input
	repository.updated = &copied
	if input.AccessToken != nil {
		repository.token.AccessToken = *input.AccessToken
	}
	if input.RefreshToken != nil {
		repository.token.RefreshToken = *input.RefreshToken
	}
	if input.ExpiresIn != nil {
		repository.token.ExpiresIn = *input.ExpiresIn
	}
	if input.ObtainmentTimestamp != nil {
		repository.token.ObtainmentTimestamp = *input.ObtainmentTimestamp
	}
	if input.Scopes != nil {
		repository.token.Scopes = input.Scopes
	}
	return &repository.token, nil
}

type fakeKickBotsRepository struct {
	mu      sync.Mutex
	bot     entity.KickBot
	updated *kickbots.UpdateTokenInput
}

func (repository *fakeKickBotsRepository) GetDefault(context.Context) (entity.KickBot, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	return repository.bot, nil
}

func (repository *fakeKickBotsRepository) GetByID(_ context.Context, id uuid.UUID) (entity.KickBot, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if repository.bot.ID != id {
		return entity.KickBot{}, kickbots.ErrNotFound
	}
	return repository.bot, nil
}

func (repository *fakeKickBotsRepository) GetByKickUserID(context.Context, uuid.UUID) (entity.KickBot, error) {
	return entity.KickBot{}, kickbots.ErrNotFound
}

func (repository *fakeKickBotsRepository) Create(context.Context, kickbots.CreateInput) (entity.KickBot, error) {
	return entity.KickBot{}, kickbots.ErrNotFound
}

func (repository *fakeKickBotsRepository) Upsert(context.Context, kickbots.UpsertInput) (entity.KickBot, error) {
	return entity.KickBot{}, kickbots.ErrNotFound
}

func (repository *fakeKickBotsRepository) UpdateToken(_ context.Context, id uuid.UUID, input kickbots.UpdateTokenInput) (entity.KickBot, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if repository.bot.ID != id {
		return entity.KickBot{}, kickbots.ErrNotFound
	}
	copied := input
	repository.updated = &copied
	repository.bot.AccessToken = input.AccessToken
	repository.bot.RefreshToken = input.RefreshToken
	repository.bot.Scopes = input.Scopes
	repository.bot.ExpiresIn = input.ExpiresIn
	repository.bot.ObtainmentTimestamp = input.ObtainmentTimestamp
	return repository.bot, nil
}
