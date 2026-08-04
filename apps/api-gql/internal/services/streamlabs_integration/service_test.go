package streamlabs_integration

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	busintegrations "github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/oauthlock"
	provider "github.com/twirapp/twir/libs/integrations/streamlabs"
	repository "github.com/twirapp/twir/libs/repositories/streamlabs_integration"
	"github.com/twirapp/twir/libs/repositories/streamlabs_integration/model"
)

func TestGetAuthLinkEncodesRequiredStateExactlyOnce(t *testing.T) {
	t.Parallel()

	client := &fakeProviderClient{}
	service := testService(&fakeRepository{}, &fakeProviderFactory{static: client}, nil, nil)

	response, err := service.GetAuthLink(context.Background(), "provider-bound-state")
	if err != nil {
		t.Fatalf("GetAuthLink() error = %v", err)
	}
	parsed, err := url.Parse(response.Link)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if got := parsed.Query()["state"]; !reflect.DeepEqual(got, []string{"provider-bound-state"}) {
		t.Fatalf("authorization URL state values = %#v, want one bound state", got)
	}
	if got := client.authStates; !reflect.DeepEqual(got, []string{"provider-bound-state"}) {
		t.Fatalf("provider GetAuthLink states = %#v, want exactly one state", got)
	}
}

func TestLifecycleRejectsMissingConfigurationAndBlankState(t *testing.T) {
	t.Parallel()

	service := testService(&fakeRepository{}, nil, nil, nil)
	if _, err := service.GetAuthLink(context.Background(), "  "); err == nil {
		t.Fatal("GetAuthLink() error = nil, want blank-state rejection")
	}
	service.config.StreamlabsClientId = ""
	if _, err := service.GetAuthLink(context.Background(), "state"); err == nil {
		t.Fatal("GetAuthLink() error = nil, want configuration error")
	}
	if err := service.PostCode(context.Background(), "channel", "code"); err == nil {
		t.Fatal("PostCode() error = nil, want configuration error")
	}
}

func TestPostCodeCreatesIntegrationFromProviderAndPublishesAdd(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{integration: model.Nil}
	client := &fakeProviderClient{
		tokens:  &provider.TokenResponse{AccessToken: "new-access", RefreshToken: "new-refresh"},
		profile: profile("Streamer", "https://avatar.test/profile.png"),
	}
	events := &fakeIntegrationEvents{}
	service := testService(repo, &fakeProviderFactory{static: client}, nil, events)

	if err := service.PostCode(context.Background(), "channel", "authorization-code"); err != nil {
		t.Fatalf("PostCode() error = %v", err)
	}
	if client.exchangedCode != "authorization-code" {
		t.Fatalf("provider exchanged code = %q, want authorization-code", client.exchangedCode)
	}
	if got := repo.integration; got.ChannelID != "channel" || !got.Enabled ||
		got.AccessToken != "new-access" || got.RefreshToken != "new-refresh" ||
		got.UserName != "Streamer" || got.Avatar != "https://avatar.test/profile.png" {
		t.Fatalf("stored integration = %#v, want mapped provider credentials/profile", got)
	}
	if want := []busintegrations.Request{{ID: repo.integration.ID.String(), Service: busintegrations.Streamlabs}}; !reflect.DeepEqual(events.added, want) {
		t.Fatalf("add events = %#v, want %#v", events.added, want)
	}
}

func TestPostCodeUpdatesIntegrationAndPublishesAdd(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	repo := &fakeRepository{integration: model.StreamlabsIntegration{
		ID: id, ChannelID: "channel", AccessToken: "old-access", RefreshToken: "old-refresh",
	}}
	client := &fakeProviderClient{
		tokens:  &provider.TokenResponse{AccessToken: "new-access", RefreshToken: "new-refresh"},
		profile: profile("Streamer", "avatar"),
	}
	events := &fakeIntegrationEvents{}
	service := testService(repo, &fakeProviderFactory{static: client}, nil, events)

	if err := service.PostCode(context.Background(), "channel", "authorization-code"); err != nil {
		t.Fatalf("PostCode() error = %v", err)
	}
	if got := repo.integration; got.ID != id || !got.Enabled || got.AccessToken != "new-access" ||
		got.RefreshToken != "new-refresh" || got.UserName != "Streamer" || got.Avatar != "avatar" {
		t.Fatalf("updated integration = %#v", got)
	}
	if want := []busintegrations.Request{{ID: id.String(), Service: busintegrations.Streamlabs}}; !reflect.DeepEqual(events.added, want) {
		t.Fatalf("add events = %#v, want %#v", events.added, want)
	}
}

func TestGetIntegrationDataMapsStoredRecordAndMissingDefault(t *testing.T) {
	t.Parallel()

	service := testService(&fakeRepository{getErr: repository.ErrNotFound}, nil, nil, nil)
	missing, err := service.GetIntegrationData(context.Background(), "channel")
	if err != nil {
		t.Fatalf("GetIntegrationData(missing) error = %v", err)
	}
	if missing.ChannelID != "channel" || missing.Enabled {
		t.Fatalf("GetIntegrationData(missing) = %#v, want disabled channel default", missing)
	}

	id := uuid.New()
	repo := &fakeRepository{integration: model.StreamlabsIntegration{
		ID: id, ChannelID: "channel", Enabled: true, UserName: "Streamer", Avatar: "avatar",
	}}
	service = testService(repo, nil, nil, nil)
	data, err := service.GetIntegrationData(context.Background(), "channel")
	if err != nil {
		t.Fatalf("GetIntegrationData(stored) error = %v", err)
	}
	if data.ID != id || !data.Enabled || data.UserName != "Streamer" || data.Avatar != "avatar" {
		t.Fatalf("GetIntegrationData(stored) = %#v", data)
	}
}

func TestAuthorizedClientUsesRepositoryTokenStoreAndSharedLocker(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{integration: model.StreamlabsIntegration{
		ChannelID: "channel", Enabled: true, AccessToken: "access", RefreshToken: "refresh",
	}}
	factory := &fakeProviderFactory{authorized: &fakeProviderClient{}}
	locker := fakeLocker{}
	service := testService(repo, factory, locker, nil)

	if _, err := service.authorizedClient(context.Background(), "channel"); err != nil {
		t.Fatalf("authorizedClient() error = %v", err)
	}
	if factory.store != repo {
		t.Fatalf("authorized client store = %T, want repository", factory.store)
	}
	if factory.locker != locker {
		t.Fatal("authorized client locker does not match service locker")
	}
}

func TestLogoutSerializesWithRefreshClearsCredentialsAndPublishesRemove(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{integration: model.StreamlabsIntegration{
		ID: uuid.New(), ChannelID: "channel", Enabled: true,
		AccessToken: "access", RefreshToken: "refresh", UserName: "Streamer", Avatar: "avatar",
	}}
	locker := &serialTestLocker{
		expectedKey: "twir:integration-token-refresh:streamlabs:channel",
		called:      make(chan string, 1),
	}
	events := &fakeIntegrationEvents{}
	service := testService(repo, nil, locker, events)

	locker.mu.Lock()
	locked := true
	defer func() {
		if locked {
			locker.mu.Unlock()
		}
	}()
	logoutDone := make(chan error, 1)
	go func() { logoutDone <- service.Logout(context.Background(), "channel") }()

	select {
	case got := <-locker.called:
		if got != locker.expectedKey {
			locker.mu.Unlock()
			t.Fatalf("Logout() lock key = %q, want %q", got, locker.expectedKey)
		}
	case err := <-logoutDone:
		locker.mu.Unlock()
		t.Fatalf("Logout() completed without refresh lock: %v", err)
	case <-time.After(time.Second):
		locker.mu.Unlock()
		t.Fatal("Logout() did not attempt to acquire refresh lock")
	}
	select {
	case err := <-logoutDone:
		t.Fatalf("Logout() completed while refresh lock held: %v", err)
	default:
	}

	if err := repo.UpdateTokens(context.Background(), "channel", provider.Tokens{
		AccessToken: "rotated-access", RefreshToken: "rotated-refresh",
	}); err != nil {
		t.Fatalf("simulated refresh UpdateTokens() error = %v", err)
	}
	locker.mu.Unlock()
	locked = false

	if err := <-logoutDone; err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if repo.integration != model.Nil {
		t.Fatalf("stored integration = %#v, want deleted credentials", repo.integration)
	}
	if want := []busintegrations.Request{{ID: "channel", Service: busintegrations.Streamlabs}}; !reflect.DeepEqual(events.removed, want) {
		t.Fatalf("remove events = %#v, want %#v", events.removed, want)
	}
}

func testService(
	repo *fakeRepository,
	factory *fakeProviderFactory,
	locker oauthlock.Locker,
	events *fakeIntegrationEvents,
) *Service {
	if factory == nil {
		factory = &fakeProviderFactory{static: &fakeProviderClient{}}
	}
	if locker == nil {
		locker = fakeLocker{}
	}
	if events == nil {
		events = &fakeIntegrationEvents{}
	}
	return &Service{
		streamlabsRepository: repo,
		config: config.Config{
			SiteBaseUrl: "https://twir.test", StreamlabsClientId: "client-id",
			StreamlabsClientSecret: "client-secret",
		},
		clientFactory: factory,
		locker:        locker,
		events:        events,
	}
}

func profile(displayName, thumbnail string) *provider.UserProfile {
	result := &provider.UserProfile{}
	result.StreamLabs.DisplayName = displayName
	result.StreamLabs.ThumbNail = thumbnail
	return result
}

type fakeRepository struct {
	integration model.StreamlabsIntegration
	getErr      error
}

func (f *fakeRepository) GetByChannelID(context.Context, string) (model.StreamlabsIntegration, error) {
	if f.getErr != nil {
		return model.Nil, f.getErr
	}
	return f.integration, nil
}

func (f *fakeRepository) Create(_ context.Context, opts repository.CreateOpts) error {
	f.integration = model.StreamlabsIntegration{
		ID: uuid.New(), ChannelID: opts.ChannelID, Enabled: opts.Enabled,
		AccessToken: opts.AccessToken, RefreshToken: opts.RefreshToken,
		UserName: opts.UserName, Avatar: opts.Avatar,
	}
	return nil
}

func (f *fakeRepository) Update(_ context.Context, opts repository.UpdateOpts) error {
	if opts.Enabled != nil {
		f.integration.Enabled = *opts.Enabled
	}
	if opts.AccessToken != nil {
		f.integration.AccessToken = *opts.AccessToken
	}
	if opts.RefreshToken != nil {
		f.integration.RefreshToken = *opts.RefreshToken
	}
	if opts.UserName != nil {
		f.integration.UserName = *opts.UserName
	}
	if opts.Avatar != nil {
		f.integration.Avatar = *opts.Avatar
	}
	return nil
}

func (f *fakeRepository) Delete(context.Context, string) error {
	f.integration = model.Nil
	return nil
}

func (f *fakeRepository) GetTokens(context.Context, string) (provider.Tokens, error) {
	if !f.integration.Enabled || f.integration.AccessToken == "" || f.integration.RefreshToken == "" {
		return provider.Tokens{}, repository.ErrNotFound
	}
	return provider.Tokens{AccessToken: f.integration.AccessToken, RefreshToken: f.integration.RefreshToken}, nil
}

func (f *fakeRepository) UpdateTokens(_ context.Context, _ string, tokens provider.Tokens) error {
	if !f.integration.Enabled {
		return repository.ErrNotFound
	}
	f.integration.AccessToken = tokens.AccessToken
	f.integration.RefreshToken = tokens.RefreshToken
	return nil
}

type fakeProviderFactory struct {
	static     *fakeProviderClient
	authorized *fakeProviderClient
	store      provider.TokenStore
	locker     oauthlock.Locker
}

func (f *fakeProviderFactory) New(string, string, string) providerClient { return f.static }

func (f *fakeProviderFactory) NewAuthorized(
	_, _, _, _ string,
	_ provider.Tokens,
	store provider.TokenStore,
	locker oauthlock.Locker,
) providerClient {
	f.store = store
	f.locker = locker
	return f.authorized
}

type fakeProviderClient struct {
	authStates    []string
	exchangedCode string
	tokens        *provider.TokenResponse
	profile       *provider.UserProfile
}

func (f *fakeProviderClient) GetAuthLink(state string) string {
	f.authStates = append(f.authStates, state)
	u, _ := url.Parse("https://streamlabs.test/api/v2.0/authorize")
	query := u.Query()
	query.Set("state", state)
	u.RawQuery = query.Encode()
	return u.String()
}

func (f *fakeProviderClient) ExchangeCode(_ context.Context, code string) (*provider.TokenResponse, error) {
	f.exchangedCode = code
	return f.tokens, nil
}

func (f *fakeProviderClient) GetProfile(context.Context) (*provider.UserProfile, error) {
	return f.profile, nil
}

type fakeIntegrationEvents struct {
	added   []busintegrations.Request
	removed []busintegrations.Request
}

func (f *fakeIntegrationEvents) PublishAdd(_ context.Context, request busintegrations.Request) error {
	f.added = append(f.added, request)
	return nil
}

func (f *fakeIntegrationEvents) PublishRemove(_ context.Context, request busintegrations.Request) error {
	f.removed = append(f.removed, request)
	return nil
}

type fakeLocker struct{}

func (fakeLocker) WithLock(ctx context.Context, _ string, fn func(context.Context) error) error {
	return fn(ctx)
}

type serialTestLocker struct {
	mu          sync.Mutex
	expectedKey string
	called      chan string
}

func (l *serialTestLocker) WithLock(ctx context.Context, key string, fn func(context.Context) error) error {
	l.called <- key
	if key != l.expectedKey {
		return fmt.Errorf("unexpected lock key %q", key)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return fn(ctx)
}
