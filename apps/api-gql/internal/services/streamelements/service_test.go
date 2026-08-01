package streamelements

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	busintegrations "github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/oauthlock"
	streamelementsintegration "github.com/twirapp/twir/libs/integrations/streamelements"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsmodel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestGetAuthLinkPreservesOAuthState(t *testing.T) {
	t.Parallel()

	service := testService(t, nil, nil, nil, nil)
	stateURL, err := service.GetAuthLink(context.Background(), "oauth-state")
	if err != nil {
		t.Fatalf("GetAuthLink() error = %v", err)
	}
	if got, want := mustQuery(t, stateURL, "state"), "oauth-state"; got != want {
		t.Fatalf("authorization URL state = %q, want %q", got, want)
	}
}

func TestLifecycleMethodsReturnConfigurationErrors(t *testing.T) {
	t.Parallel()

	service := testService(t, nil, nil, nil, nil)
	service.config.StreamElementsClientId = ""

	if _, err := service.GetAuthLink(context.Background(), "oauth-state"); err == nil {
		t.Fatal("GetAuthLink() error = nil, want configuration error")
	}
	if err := service.PostCode(context.Background(), "channel", "code"); err == nil {
		t.Fatal("PostCode() error = nil, want configuration error")
	}
}

func TestGetDataReturnsNilUntilProfileDataExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		integration channelsintegrationsmodel.ChannelIntegration
	}{
		{name: "missing integration", integration: channelsintegrationsmodel.Nil},
		{name: "missing data", integration: channelsintegrationsmodel.ChannelIntegration{ID: "channel-integration"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeChannelIntegrationsRepository{integration: test.integration}
			service := testService(t, repo, nil, nil, nil)

			data, err := service.GetData(context.Background(), "channel")
			if err != nil {
				t.Fatalf("GetData() error = %v", err)
			}
			if data != nil {
				t.Fatalf("GetData() = %#v, want nil", data)
			}
		})
	}
}

func TestGetDataMapsStoredProfile(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
		ID: "channel-integration",
		Data: &channelsintegrationsmodel.Data{
			UserName: lo.ToPtr("streamer"),
			Avatar:   lo.ToPtr("https://avatar.test/profile.png"),
		},
	}}
	service := testService(t, repo, nil, nil, nil)

	data, err := service.GetData(context.Background(), "channel")
	if err != nil {
		t.Fatalf("GetData() error = %v", err)
	}
	if want := (&IntegrationData{UserName: "streamer", Avatar: "https://avatar.test/profile.png"}); !reflect.DeepEqual(data, want) {
		t.Fatalf("GetData() = %#v, want %#v", data, want)
	}
}

func TestPostCodeCreatesPersistentIntegrationAndPublishesAdd(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.Nil}
	client := &fakeProviderClient{
		exchangedTokens: &streamelementsintegration.TokenResponse{
			AccessToken: "new-access", RefreshToken: "new-refresh",
		},
		profile: &streamelementsintegration.UserProfile{
			ID: "provider-channel", DisplayName: "Streamer", Avatar: "https://avatar.test/profile.png",
		},
	}
	factory := &fakeProviderClientFactory{static: client}
	events := &fakeIntegrationEvents{}
	service := testService(t, repo, factory, nil, events)

	if err := service.PostCode(context.Background(), "channel", "authorization-code"); err != nil {
		t.Fatalf("PostCode() error = %v", err)
	}

	want := channelsintegrationsmodel.ChannelIntegration{
		ID: "created-channel-integration", ChannelID: "channel", IntegrationID: "provider-integration",
		Enabled: true, AccessToken: lo.ToPtr("new-access"), RefreshToken: lo.ToPtr("new-refresh"),
		Data: &channelsintegrationsmodel.Data{
			UserName: lo.ToPtr("Streamer"), Avatar: lo.ToPtr("https://avatar.test/profile.png"),
		},
	}
	if !reflect.DeepEqual(repo.integration, want) {
		t.Fatalf("stored integration = %#v, want %#v", repo.integration, want)
	}
	if got, wantEvent := events.added, []busintegrations.Request{{
		ID: "created-channel-integration", Service: busintegrations.StreamElements,
	}}; !reflect.DeepEqual(got, wantEvent) {
		t.Fatalf("add events = %#v, want %#v", got, wantEvent)
	}
	if got, wantCode := client.exchangedCode, "authorization-code"; got != wantCode {
		t.Fatalf("exchanged code = %q, want %q", got, wantCode)
	}
}

func TestPostCodeUpdatesPersistentIntegrationAndPublishesAdd(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
		ID: "existing-channel-integration", ChannelID: "channel", IntegrationID: "provider-integration",
		AccessToken: lo.ToPtr("old-access"), RefreshToken: lo.ToPtr("old-refresh"),
	}}
	client := &fakeProviderClient{
		exchangedTokens: &streamelementsintegration.TokenResponse{
			AccessToken: "new-access", RefreshToken: "new-refresh",
		},
		profile: &streamelementsintegration.UserProfile{ID: "provider-channel", Username: "streamer", ProfileImage: "https://avatar.test/profile.png"},
	}
	events := &fakeIntegrationEvents{}
	service := testService(t, repo, &fakeProviderClientFactory{static: client}, nil, events)

	if err := service.PostCode(context.Background(), "channel", "authorization-code"); err != nil {
		t.Fatalf("PostCode() error = %v", err)
	}
	if !repo.integration.Enabled {
		t.Fatal("stored integration enabled = false, want true")
	}
	if got := dereference(repo.integration.AccessToken); got != "new-access" {
		t.Fatalf("stored access token = %q, want new-access", got)
	}
	if got := dereference(repo.integration.RefreshToken); got != "new-refresh" {
		t.Fatalf("stored refresh token = %q, want new-refresh", got)
	}
	if got := dereference(repo.integration.Data.UserName); got != "streamer" {
		t.Fatalf("stored username = %q, want streamer", got)
	}
	if got := dereference(repo.integration.Data.Avatar); got != "https://avatar.test/profile.png" {
		t.Fatalf("stored avatar = %q, want profile image fallback", got)
	}
	if got, want := events.added, []busintegrations.Request{{
		ID: "existing-channel-integration", Service: busintegrations.StreamElements,
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("add events = %#v, want %#v", got, want)
	}
}

func TestLogoutClearsPersistentCredentialsAndPublishesRemove(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
		ID: "channel-integration", ChannelID: "channel", IntegrationID: "provider-integration", Enabled: true,
		AccessToken: lo.ToPtr("access"), RefreshToken: lo.ToPtr("refresh"),
		Data: &channelsintegrationsmodel.Data{UserName: lo.ToPtr("streamer"), Avatar: lo.ToPtr("avatar")},
	}}
	events := &fakeIntegrationEvents{}
	service := testService(t, repo, nil, nil, events)

	if err := service.Logout(context.Background(), "channel"); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if repo.integration.Enabled {
		t.Fatal("stored integration enabled = true, want false")
	}
	if repo.integration.AccessToken != nil || repo.integration.RefreshToken != nil {
		t.Fatalf("stored tokens = (%v, %v), want nil", repo.integration.AccessToken, repo.integration.RefreshToken)
	}
	if repo.integration.Data == nil || repo.integration.Data.UserName != nil || repo.integration.Data.Avatar != nil {
		t.Fatalf("stored profile data = %#v, want empty", repo.integration.Data)
	}
	if got, want := events.removed, []busintegrations.Request{{
		ID: "channel", Service: busintegrations.StreamElements,
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("remove events = %#v, want %#v", got, want)
	}
}

func TestLogoutMissingIntegrationIsNoop(t *testing.T) {
	t.Parallel()

	events := &fakeIntegrationEvents{}
	service := testService(t, &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.Nil}, nil, nil, events)
	if err := service.Logout(context.Background(), "channel"); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if len(events.removed) != 0 {
		t.Fatalf("remove events = %#v, want none", events.removed)
	}
}

func TestImportsComposeNormalizationAndSharedImporterReports(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
		ID: "channel-integration", ChannelID: "channel", Enabled: true,
		AccessToken: lo.ToPtr("access"), RefreshToken: lo.ToPtr("refresh"),
	}}
	timer := streamelementsintegration.Timer{Name: "timer", Message: "timer response", Enabled: true, ChatLines: 3}
	timer.Online.Enabled = true
	timer.Online.Interval = 5
	badTimer := streamelementsintegration.Timer{Name: "bad-timer", Message: "bad response", Enabled: true, ChatLines: 2}
	badTimer.Online.Enabled = true
	badTimer.Online.Interval = 5
	badTimer.Offline.Enabled = true
	badTimer.Offline.Interval = 10
	client := &fakeProviderClient{
		profile: &streamelementsintegration.UserProfile{ID: "provider-channel"},
		commands: []streamelementsintegration.Command{
			{Name: "command", Response: "response", Enabled: true, AccessLevel: 100, Type: "say"},
			{Name: "regular", Response: "regular response", Enabled: true, AccessLevel: 300, Type: "say"},
		},
		timers: []streamelementsintegration.Timer{timer, badTimer},
	}
	factory := &fakeProviderClientFactory{authorized: client}
	sharedImporter := &fakeProviderImporter{
		commandReport: importer.Report{
			ImportedCount: 1, FailedCount: 1,
			Failures: []importer.Failure{{Name: "command", Reason: importer.FailureDuplicate}},
		},
		timerReport: importer.Report{ImportedCount: 1, Failures: []importer.Failure{}},
	}
	service := testService(t, repo, factory, sharedImporter, nil)

	commandReport, err := service.ImportCommands(context.Background(), "channel", "actor")
	if err != nil {
		t.Fatalf("ImportCommands() error = %v", err)
	}
	if want := (importer.Report{
		ImportedCount: 1, FailedCount: 2,
		Failures: []importer.Failure{
			{Name: "regular", Reason: importer.FailureUnsupportedRole},
			{Name: "command", Reason: importer.FailureDuplicate},
		},
	}); !reflect.DeepEqual(commandReport, want) {
		t.Fatalf("ImportCommands() report = %#v, want %#v", commandReport, want)
	}
	if want := []importer.Command{{
		Name: "command", Response: "response", Enabled: true, Visible: true,
		Aliases: []string{}, Role: importer.RoleEveryone,
	}}; !reflect.DeepEqual(sharedImporter.commands, want) {
		t.Fatalf("shared importer commands = %#v, want %#v", sharedImporter.commands, want)
	}

	timerReport, err := service.ImportTimers(context.Background(), "channel", "actor")
	if err != nil {
		t.Fatalf("ImportTimers() error = %v", err)
	}
	if want := (importer.Report{
		ImportedCount: 1, FailedCount: 1,
		Failures: []importer.Failure{{Name: "bad-timer", Reason: importer.FailureIncompatibleInterval}},
	}); !reflect.DeepEqual(timerReport, want) {
		t.Fatalf("ImportTimers() report = %#v, want %#v", timerReport, want)
	}
	if want := []importer.Timer{{
		Name: "timer", Message: "timer response", Enabled: true, OnlineEnabled: true,
		TimeInterval: 5, MessageInterval: 3,
	}}; !reflect.DeepEqual(sharedImporter.timers, want) {
		t.Fatalf("shared importer timers = %#v, want %#v", sharedImporter.timers, want)
	}

	if got, want := factory.authorizedCalls, 2; got != want {
		t.Fatalf("authorized client creations = %d, want %d", got, want)
	}
	if factory.store != service {
		t.Fatalf("authorized client token store = %T, want service", factory.store)
	}
	if factory.locker != service.locker {
		t.Fatal("authorized client locker does not match service locker")
	}
}

func TestTokenStoreReadsAndAtomicallyUpdatesRepositoryTokens(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
		ID: "channel-integration", ChannelID: "channel", Enabled: true,
		AccessToken: lo.ToPtr("access"), RefreshToken: lo.ToPtr("refresh"),
	}}
	service := testService(t, repo, nil, nil, nil)

	tokens, err := service.GetTokens(context.Background(), "channel")
	if err != nil {
		t.Fatalf("GetTokens() error = %v", err)
	}
	if want := (streamelementsintegration.Tokens{AccessToken: "access", RefreshToken: "refresh"}); tokens != want {
		t.Fatalf("GetTokens() = %#v, want %#v", tokens, want)
	}

	newTokens := streamelementsintegration.Tokens{AccessToken: "new-access", RefreshToken: "new-refresh"}
	if err := service.UpdateTokens(context.Background(), "channel", newTokens); err != nil {
		t.Fatalf("UpdateTokens() error = %v", err)
	}
	if got := dereference(repo.integration.AccessToken); got != "new-access" {
		t.Fatalf("stored access token = %q, want new-access", got)
	}
	if got := dereference(repo.integration.RefreshToken); got != "new-refresh" {
		t.Fatalf("stored refresh token = %q, want new-refresh", got)
	}
}

func TestTokenStoreRejectsDisabledIntegration(t *testing.T) {
	t.Parallel()

	repo := &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
		ID: "channel-integration", ChannelID: "channel", Enabled: false,
		AccessToken: lo.ToPtr("access"), RefreshToken: lo.ToPtr("refresh"),
	}}
	service := testService(t, repo, nil, nil, nil)

	if _, err := service.GetTokens(context.Background(), "channel"); err == nil {
		t.Fatal("GetTokens() error = nil, want disabled integration rejection")
	}
}

func testService(
	t *testing.T,
	channelRepo *fakeChannelIntegrationsRepository,
	clientFactory *fakeProviderClientFactory,
	sharedImporter *fakeProviderImporter,
	events *fakeIntegrationEvents,
) *Service {
	t.Helper()
	if channelRepo == nil {
		channelRepo = &fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.Nil}
	}
	if clientFactory == nil {
		clientFactory = &fakeProviderClientFactory{realStatic: true}
	}
	if sharedImporter == nil {
		sharedImporter = &fakeProviderImporter{}
	}
	if events == nil {
		events = &fakeIntegrationEvents{}
	}

	return &Service{
		config: config.Config{
			StreamElementsClientId: "client-id", StreamElementsClientSecret: "client-secret",
		},
		redirectURL:             "https://twir.test/dashboard/integrations/callbacks/streamelements",
		importer:                sharedImporter,
		integrationsRepo:        fakeIntegrationsRepository{integration: integrationsmodel.Integration{ID: "provider-integration", Service: integrationsmodel.ServiceStreamElements}},
		channelIntegrationsRepo: channelRepo,
		clientFactory:           clientFactory,
		locker:                  fakeLocker{},
		events:                  events,
	}
}

type fakeIntegrationsRepository struct {
	integration integrationsmodel.Integration
}

func (f fakeIntegrationsRepository) GetByService(context.Context, integrationsmodel.Service) (integrationsmodel.Integration, error) {
	return f.integration, nil
}

type fakeChannelIntegrationsRepository struct {
	integration channelsintegrationsmodel.ChannelIntegration
}

func (f *fakeChannelIntegrationsRepository) GetByChannelAndService(
	context.Context,
	string,
	integrationsmodel.Service,
) (channelsintegrationsmodel.ChannelIntegration, error) {
	return f.integration, nil
}

func (f *fakeChannelIntegrationsRepository) Create(
	_ context.Context,
	input channelsintegrations.CreateInput,
) (channelsintegrationsmodel.ChannelIntegration, error) {
	f.integration = channelsintegrationsmodel.ChannelIntegration{
		ID: "created-channel-integration", ChannelID: input.ChannelID, IntegrationID: input.IntegrationID,
		Enabled: input.Enabled, AccessToken: input.AccessToken, RefreshToken: input.RefreshToken, Data: input.Data,
	}
	return f.integration, nil
}

func (f *fakeChannelIntegrationsRepository) Update(
	_ context.Context,
	_ string,
	input channelsintegrations.UpdateInput,
) error {
	if input.Enabled != nil {
		f.integration.Enabled = *input.Enabled
	}
	if input.AccessToken != nil {
		f.integration.AccessToken = input.AccessToken
	} else if input.ClearAccessToken {
		f.integration.AccessToken = nil
	}
	if input.RefreshToken != nil {
		f.integration.RefreshToken = input.RefreshToken
	} else if input.ClearRefreshToken {
		f.integration.RefreshToken = nil
	}
	if input.Data != nil {
		f.integration.Data = input.Data
	}
	return nil
}

type fakeProviderClientFactory struct {
	static          *fakeProviderClient
	authorized      *fakeProviderClient
	realStatic      bool
	authorizedCalls int
	store           streamelementsintegration.TokenStore
	locker          oauthlock.Locker
}

func (f *fakeProviderClientFactory) NewStatic(clientID, clientSecret string) providerClient {
	if f.realStatic {
		return streamelementsintegration.NewStatic(clientID, clientSecret)
	}
	return f.static
}

func (f *fakeProviderClientFactory) NewAuthorized(
	_, _, _, _ string,
	_ streamelementsintegration.Tokens,
	store streamelementsintegration.TokenStore,
	locker oauthlock.Locker,
) providerClient {
	f.authorizedCalls++
	f.store = store
	f.locker = locker
	return f.authorized
}

type fakeProviderClient struct {
	exchangedCode   string
	exchangedTokens *streamelementsintegration.TokenResponse
	profile         *streamelementsintegration.UserProfile
	commands        []streamelementsintegration.Command
	timers          []streamelementsintegration.Timer
}

func (f *fakeProviderClient) GetAuthLinkWithState(redirectURL, state string) (string, error) {
	u, _ := url.Parse("https://api.streamelements.test/oauth2/authorize")
	query := u.Query()
	query.Set("redirect_uri", redirectURL)
	query.Set("state", state)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (f *fakeProviderClient) ExchangeCode(_ context.Context, code, _ string) (*streamelementsintegration.TokenResponse, error) {
	f.exchangedCode = code
	return f.exchangedTokens, nil
}

func (f *fakeProviderClient) GetProfile(context.Context) (*streamelementsintegration.UserProfile, error) {
	return f.profile, nil
}

func (f *fakeProviderClient) GetCommands(context.Context, string) ([]streamelementsintegration.Command, error) {
	return f.commands, nil
}

func (f *fakeProviderClient) GetTimers(context.Context, string) ([]streamelementsintegration.Timer, error) {
	return f.timers, nil
}

type fakeProviderImporter struct {
	commands      []importer.Command
	commandReport importer.Report
	timers        []importer.Timer
	timerReport   importer.Report
}

func (f *fakeProviderImporter) ImportCommands(
	_ context.Context,
	_, _ string,
	commands []importer.Command,
) (importer.Report, error) {
	f.commands = commands
	return f.commandReport, nil
}

func (f *fakeProviderImporter) ImportTimers(
	_ context.Context,
	_, _ string,
	timers []importer.Timer,
) (importer.Report, error) {
	f.timers = timers
	return f.timerReport, nil
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

func mustQuery(t *testing.T, value, key string) string {
	t.Helper()
	parsed, err := url.Parse(value)
	if err != nil {
		t.Fatalf("parse URL %q: %v", value, err)
	}
	return parsed.Query().Get(key)
}

func dereference(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
