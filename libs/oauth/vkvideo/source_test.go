package vkvideo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/crypto"
	entity "github.com/twirapp/twir/libs/entities/vk_video_bot"
	"github.com/twirapp/twir/libs/oauth"
	tokens "github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	"github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

const testCipherKey = "0123456789abcdef0123456789abcdef"

func TestUserTokenSourceRefreshesAndEncryptsCredential(t *testing.T) {
	// Given
	server, calls := newOAuthServer(t, `{"access_token":"new-access","refresh_token":"new-refresh","expires_in":3600,"scope":"chat"}`)
	userID := uuid.New()
	repo := &userRepository{token: encryptedUserToken(t, userID, "old-access", "old-refresh")}
	source := newUserSource(t, server.URL, repo)

	// When
	credential, err := source.Token(context.Background(), userID)

	// Then
	if err != nil || credential.AccessToken != "new-access" || calls.Load() != 1 {
		t.Fatalf("refresh credential = %#v, %v, calls=%d", credential, err, calls.Load())
	}
	accessToken, decryptErr := crypto.Decrypt(repo.token.AccessToken, testCipherKey)
	if decryptErr != nil || accessToken != "new-access" {
		t.Fatalf("stored access token = %q, %v", accessToken, decryptErr)
	}
}

func TestSingletonBotTokenSourcePreservesRefreshAndTransaction(t *testing.T) {
	// Given
	server, _ := newOAuthServer(t, `{"access_token":"new-access","expires_in":3600}`)
	refreshToken := "old-refresh"
	repo := &botRepository{bot: encryptedBot(t, "old-access", refreshToken)}
	runner := &transactionRunner{}
	source := newBotSource(t, server.URL, repo, runner)

	// When
	credential, err := source.Token(context.Background())

	// Then
	if err != nil || credential.AccessToken != "new-access" || runner.calls != 1 || repo.updates != 1 {
		t.Fatalf("refresh credential = %#v, %v, transactions=%d updates=%d", credential, err, runner.calls, repo.updates)
	}
	storedRefresh, decryptErr := crypto.Decrypt(repo.bot.EncryptedRefreshToken, testCipherKey)
	if decryptErr != nil || storedRefresh != refreshToken {
		t.Fatalf("stored refresh token = %q, %v", storedRefresh, decryptErr)
	}
}

func TestUserTokenSourceRefreshesOnceWhenConcurrent(t *testing.T) {
	// Given
	server, calls := newOAuthServer(t, `{"access_token":"new-access","refresh_token":"new-refresh","expires_in":3600}`)
	userID := uuid.New()
	repo := &userRepository{token: encryptedUserToken(t, userID, "old-access", "old-refresh")}
	first := newUserSource(t, server.URL, repo)
	second := newUserSource(t, server.URL, repo)
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, source := range []UserTokenSource{first, second} {
		go func(source UserTokenSource) {
			<-start
			_, err := source.Token(context.Background(), userID)
			results <- err
		}(source)
	}

	// When
	close(start)
	firstErr, secondErr := <-results, <-results

	// Then
	if firstErr != nil || secondErr != nil || calls.Load() != 1 {
		t.Fatalf("refresh errors = %v, %v; calls=%d", firstErr, secondErr, calls.Load())
	}
}

func TestUserTokenSourceReturnsTypedErrorWhenRefreshIsMissing(t *testing.T) {
	// Given
	userID := uuid.New()
	repo := &userRepository{token: encryptedUserToken(t, userID, "old-access", "")}
	source := newUserSource(t, "http://127.0.0.1:1", repo)

	// When
	_, err := source.Token(context.Background(), userID)

	// Then
	if !errors.Is(err, ErrMissingRefreshToken) || !errors.Is(err, oauth.ErrInvalidCredential) {
		t.Fatalf("missing refresh error = %v", err)
	}
}

func TestUserTokenSourceRedactsMalformedCiphertext(t *testing.T) {
	// Given
	userID := uuid.New()
	repo := &userRepository{token: &tokenmodel.Token{ID: uuid.New(), AccessToken: "not-hex", RefreshToken: "TWIR_VK_TEST_REFRESH_SENTINEL", ExpiresIn: 1, ObtainmentTimestamp: time.Now().UTC()}}
	source := newUserSource(t, "http://127.0.0.1:1", repo)

	// When
	_, err := source.Token(context.Background(), userID)

	// Then
	if err == nil || strings.Contains(err.Error(), "TWIR_VK_TEST_REFRESH_SENTINEL") {
		t.Fatalf("malformed ciphertext error = %v", err)
	}
}

func newOAuthServer(t *testing.T, response string) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/oauth/server/token" || request.FormValue("grant_type") != "refresh_token" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		calls.Add(1)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return server, &calls
}

func newUserSource(t *testing.T, serverURL string, repo *userRepository) UserTokenSource {
	t.Helper()
	source, err := NewUserTokenSource(testOptions(t, serverURL), repo)
	if err != nil {
		t.Fatal(err)
	}
	return source
}

func newBotSource(t *testing.T, serverURL string, repo *botRepository, runner *transactionRunner) SingletonBotTokenSource {
	t.Helper()
	options := testOptions(t, serverURL)
	options.TransactionRunner = runner
	source, err := NewSingletonBotTokenSource(options, repo)
	if err != nil {
		t.Fatal(err)
	}
	return source
}

func testOptions(t *testing.T, serverURL string) SourceOptions {
	t.Helper()
	address := os.Getenv("TWIR_VK_TEST_REDIS_ADDR")
	if address == "" {
		t.Skip("TWIR_VK_TEST_REDIS_ADDR is required")
	}
	client := redis.NewClient(&redis.Options{Addr: address})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return SourceOptions{ClientID: "client", ClientSecret: "secret", RedirectURL: "https://callback.invalid", APIBaseURL: serverURL, AuthBaseURL: serverURL, DevAPIBaseURL: serverURL, Redis: client, CipherKey: testCipherKey}
}

func encryptedUserToken(t *testing.T, _ uuid.UUID, accessToken string, refreshToken string) *tokenmodel.Token {
	t.Helper()
	access, err := crypto.Encrypt(accessToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := crypto.Encrypt(refreshToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	return &tokenmodel.Token{ID: uuid.New(), AccessToken: access, RefreshToken: refresh, ExpiresIn: 1, ObtainmentTimestamp: time.Now().Add(-time.Hour), Scopes: []string{"old"}}
}

func encryptedBot(t *testing.T, accessToken string, refreshToken string) entity.VKVideoBot {
	t.Helper()
	access, err := crypto.Encrypt(accessToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := crypto.Encrypt(refreshToken, testCipherKey)
	if err != nil {
		t.Fatal(err)
	}
	return entity.VKVideoBot{ID: uuid.New(), EncryptedAccessToken: access, EncryptedRefreshToken: refresh, ExpiresIn: 1, ObtainmentTimestamp: time.Now().Add(-time.Hour), VKUserID: uuid.New()}
}

type transactionRunner struct{ calls int }

func (r *transactionRunner) Do(ctx context.Context, fn func(context.Context) error) error {
	r.calls++
	return fn(ctx)
}

type userRepository struct {
	mu    sync.Mutex
	token *tokenmodel.Token
}

func (r *userRepository) GetByUserID(context.Context, uuid.UUID) (*tokenmodel.Token, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	token := *r.token
	return &token, nil
}
func (r *userRepository) UpdateTokenByID(_ context.Context, _ uuid.UUID, input tokens.UpdateTokenInput) (*tokenmodel.Token, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.token.AccessToken = *input.AccessToken
	r.token.RefreshToken = *input.RefreshToken
	r.token.ExpiresIn = *input.ExpiresIn
	r.token.ObtainmentTimestamp = *input.ObtainmentTimestamp
	if len(input.Scopes) > 0 {
		r.token.Scopes = input.Scopes
	}
	token := *r.token
	return &token, nil
}
func (r *userRepository) GetByID(context.Context, uuid.UUID) (*tokenmodel.Token, error) {
	return r.GetByUserID(context.Background(), uuid.Nil)
}
func (r *userRepository) GetByBotID(context.Context, string) (*tokenmodel.Token, error) {
	return nil, errors.New("not implemented")
}
func (r *userRepository) CreateUserToken(context.Context, tokens.CreateInput) (*tokenmodel.Token, error) {
	return nil, errors.New("not implemented")
}

type botRepository struct {
	mu      sync.Mutex
	bot     entity.VKVideoBot
	updates int
}

func (r *botRepository) Get(context.Context) (entity.VKVideoBot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.bot, nil
}
func (r *botRepository) Lock(context.Context) error { return nil }
func (r *botRepository) Upsert(context.Context, vk_video_bots.UpsertInput) (entity.VKVideoBot, error) {
	return entity.VKVideoBot{}, errors.New("not implemented")
}
func (r *botRepository) Update(_ context.Context, input vk_video_bots.UpdateInput) (entity.VKVideoBot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bot.EncryptedAccessToken = input.EncryptedAccessToken
	r.bot.EncryptedRefreshToken = input.EncryptedRefreshToken
	r.bot.Scopes = input.Scopes
	r.bot.ExpiresIn = input.ExpiresIn
	r.bot.ObtainmentTimestamp = input.ObtainmentTimestamp
	r.bot.VKUserID = input.VKUserID
	r.updates++
	return r.bot, nil
}
