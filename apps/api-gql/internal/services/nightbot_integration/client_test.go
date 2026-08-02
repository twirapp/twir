package nightbot_integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsmodel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestGetAuthLinkEncodesProvidedState(t *testing.T) {
	t.Parallel()

	clientID := "nightbot-client"
	clientSecret := "nightbot-secret"
	redirectURL := "https://twir.test/dashboard/integrations/callbacks/nightbot"
	service := &Service{
		integrationsRepo: fakeIntegrationsRepository{integration: integrationsmodel.Integration{
			Service:      integrationsmodel.ServiceNightbot,
			ClientID:     &clientID,
			ClientSecret: &clientSecret,
			RedirectURL:  &redirectURL,
		}},
	}

	link, err := service.GetAuthLink(context.Background(), "state value+/=")
	if err != nil {
		t.Fatalf("GetAuthLink() error = %v", err)
	}

	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if got, want := parsed.Query().Get("state"), "state value+/="; got != want {
		t.Fatalf("authorization URL state = %q, want %q", got, want)
	}
}

func TestGetAuthLinkRejectsBlankState(t *testing.T) {
	t.Parallel()

	clientID := "nightbot-client"
	clientSecret := "nightbot-secret"
	redirectURL := "https://twir.test/dashboard/integrations/callbacks/nightbot"
	service := &Service{integrationsRepo: fakeIntegrationsRepository{integration: integrationsmodel.Integration{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		RedirectURL:  &redirectURL,
	}}}
	if _, err := service.GetAuthLink(context.Background(), " \t\n "); err == nil {
		t.Fatal("GetAuthLink() error = nil, want empty state rejection")
	}
}

func TestImportProviderErrorsDoNotLeakResponseBodies(t *testing.T) {
	t.Parallel()

	const sentinel = "SENTINEL_PROVIDER_SECRET_BODY"
	accessToken := "nightbot-access"
	for _, testCase := range []struct {
		name      string
		path      string
		operation string
		call      func(*Service) error
	}{
		{
			name:      "commands",
			path:      "/1/commands",
			operation: "commands",
			call: func(service *Service) error {
				_, err := service.getCommands(context.Background(), "channel")
				return err
			},
		},
		{
			name:      "timers",
			path:      "/1/timers",
			operation: "timers",
			call: func(service *Service) error {
				_, err := service.getTimers(context.Background(), "channel")
				return err
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path != testCase.path {
					t.Fatalf("request path = %q, want %q", request.URL.Path, testCase.path)
				}
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Body:       io.NopCloser(strings.NewReader(sentinel)),
					Header:     make(http.Header),
				}, nil
			})}
			service := &Service{
				httpClient: client,
				channelIntegrationsRepo: fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
					ID:          "nightbot-integration",
					AccessToken: &accessToken,
				}},
			}

			err := testCase.call(service)
			if err == nil {
				t.Fatal("provider call error = nil")
			}
			if strings.Contains(err.Error(), sentinel) {
				t.Fatalf("provider error leaked response body: %v", err)
			}
			if !strings.Contains(err.Error(), "Nightbot "+testCase.operation) ||
				!strings.Contains(err.Error(), "status 418") {
				t.Fatalf("provider error = %q, want operation and status", err)
			}
		})
	}
}

func TestOAuthProviderErrorsDoNotLeakResponseBodies(t *testing.T) {
	t.Parallel()

	const sentinel = "SENTINEL_OAUTH_SECRET_BODY"
	clientID := "nightbot-client"
	clientSecret := "nightbot-secret"
	redirectURL := "https://twir.test/dashboard/integrations/callbacks/nightbot"
	for _, testCase := range []struct {
		name       string
		operation  string
		statusCode int
		profile    bool
	}{
		{
			name:       "token",
			operation:  "token",
			statusCode: http.StatusTooManyRequests,
		},
		{
			name:       "profile",
			operation:  "profile",
			statusCode: http.StatusBadGateway,
			profile:    true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path == "/oauth2/token" && testCase.profile {
					return providerResponse(http.StatusOK, `{"access_token":"access","refresh_token":"refresh"}`), nil
				}
				expectedPath := "/oauth2/token"
				if testCase.profile {
					expectedPath = "/1/me"
				}
				if request.URL.Path != expectedPath {
					t.Fatalf("request path = %q, want %q", request.URL.Path, expectedPath)
				}
				return providerResponse(testCase.statusCode, sentinel), nil
			})

			service := &Service{
				httpClient: &http.Client{Transport: transport},
				integrationsRepo: fakeIntegrationsRepository{integration: integrationsmodel.Integration{
					ID:           "nightbot-provider",
					ClientID:     &clientID,
					ClientSecret: &clientSecret,
					RedirectURL:  &redirectURL,
				}},
				channelIntegrationsRepo: fakeChannelIntegrationsRepository{},
			}

			err := service.PostCode(context.Background(), "channel", "authorization-code")
			if err == nil {
				t.Fatal("PostCode() error = nil")
			}
			if strings.Contains(err.Error(), sentinel) {
				t.Fatalf("OAuth provider error leaked response body: %v", err)
			}
			wantStatus := fmt.Sprintf("status %d", testCase.statusCode)
			if !strings.Contains(err.Error(), "Nightbot "+testCase.operation) ||
				!strings.Contains(err.Error(), wantStatus) {
				t.Fatalf("OAuth provider error = %q, want operation and %s", err, wantStatus)
			}
		})
	}
}

func providerResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestImportCommandsPrependsNormalizationFailuresToSharedImporterReport(t *testing.T) {
	accessToken := "nightbot-access"
	requests := 0
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests++
		if got, want := request.URL.String(), "https://api.nightbot.tv/1/commands"; got != want {
			t.Fatalf("request URL = %q, want %q", got, want)
		}
		if got, want := request.Header.Get("Authorization"), "Bearer nightbot-access"; got != want {
			t.Fatalf("authorization = %q, want %q", got, want)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
				"commands": [
					{"name":"!Admin","message":"admin response","userLevel":"admin"},
					{"name":"!Hello","message":"hello response","userLevel":"everyone","coolDown":5}
				]
			}`)),
			Header: make(http.Header),
		}, nil
	})}
	sharedImporter := &fakeSharedImporter{commandReport: importer.Report{
		ImportedCount: 1,
		FailedCount:   1,
		Failures:      []importer.Failure{{Name: "hello", Reason: importer.FailureDuplicate}},
	}}
	service := &Service{
		httpClient: client,
		importer:   sharedImporter,
		channelIntegrationsRepo: fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
			ID:          "nightbot-integration",
			AccessToken: &accessToken,
		}},
	}

	result, err := service.ImportCommands(context.Background(), "channel", "actor")
	if err != nil {
		t.Fatalf("ImportCommands() error = %v", err)
	}
	if want := (importer.Report{
		ImportedCount: 1,
		FailedCount:   2,
		Failures: []importer.Failure{
			{Name: "!Admin", Reason: importer.FailureUnsupportedRole},
			{Name: "hello", Reason: importer.FailureDuplicate},
		},
	}); !reflect.DeepEqual(result, want) {
		t.Fatalf("ImportCommands() result = %#v, want %#v", result, want)
	}
	if want := ([]importer.Command{{
		Name: "hello", Response: "hello response", Enabled: true, Visible: true, IsReply: true,
		Aliases: []string{}, Cooldown: 5, Role: importer.RoleEveryone,
	}}); !reflect.DeepEqual(sharedImporter.commands, want) {
		t.Fatalf("shared importer commands = %#v, want %#v", sharedImporter.commands, want)
	}
	if got, want := requests, 1; got != want {
		t.Fatalf("Nightbot command requests = %d, want %d", got, want)
	}
}

func TestImportTimersPrependsNormalizationFailuresToSharedImporterReport(t *testing.T) {
	accessToken := "nightbot-access"
	requests := 0
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests++
		if got, want := request.URL.String(), "https://api.nightbot.tv/1/timers"; got != want {
			t.Fatalf("request URL = %q, want %q", got, want)
		}
		if got, want := request.Header.Get("Authorization"), "Bearer nightbot-access"; got != want {
			t.Fatalf("authorization = %q, want %q", got, want)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
				"timers": [
					{"name":"daily","message":"daily response","interval":"0 12 * * *","lines":1,"enabled":true},
					{"name":"five","message":"five response","interval":"*/5 * * * *","lines":2,"enabled":true}
				]
			}`)),
			Header: make(http.Header),
		}, nil
	})}
	sharedImporter := &fakeSharedImporter{timerReport: importer.Report{
		ImportedCount: 1,
		FailedCount:   1,
		Failures:      []importer.Failure{{Name: "five", Reason: importer.FailureDuplicate}},
	}}
	service := &Service{
		httpClient: client,
		importer:   sharedImporter,
		channelIntegrationsRepo: fakeChannelIntegrationsRepository{integration: channelsintegrationsmodel.ChannelIntegration{
			ID:          "nightbot-integration",
			AccessToken: &accessToken,
		}},
	}

	result, err := service.ImportTimers(context.Background(), "channel", "actor")
	if err != nil {
		t.Fatalf("ImportTimers() error = %v", err)
	}
	if want := (importer.Report{
		ImportedCount: 1,
		FailedCount:   2,
		Failures: []importer.Failure{
			{Name: "daily", Reason: importer.FailureIncompatibleInterval},
			{Name: "five", Reason: importer.FailureDuplicate},
		},
	}); !reflect.DeepEqual(result, want) {
		t.Fatalf("ImportTimers() result = %#v, want %#v", result, want)
	}
	if want := ([]importer.Timer{{
		Name: "five", Message: "five response", Enabled: true, OnlineEnabled: true,
		TimeInterval: 5, MessageInterval: 2,
	}}); !reflect.DeepEqual(sharedImporter.timers, want) {
		t.Fatalf("shared importer timers = %#v, want %#v", sharedImporter.timers, want)
	}
	if got, want := requests, 1; got != want {
		t.Fatalf("Nightbot timer requests = %d, want %d", got, want)
	}
}

type fakeIntegrationsRepository struct {
	integration integrationsmodel.Integration
}

func (f fakeIntegrationsRepository) GetByService(
	_ context.Context,
	_ integrationsmodel.Service,
) (integrationsmodel.Integration, error) {
	return f.integration, nil
}

type fakeSharedImporter struct {
	commands      []importer.Command
	commandReport importer.Report
	timers        []importer.Timer
	timerReport   importer.Report
}

func (f *fakeSharedImporter) ImportCommands(
	_ context.Context,
	_, _ string,
	commands []importer.Command,
) (importer.Report, error) {
	f.commands = commands
	return f.commandReport, nil
}

func (f *fakeSharedImporter) ImportTimers(
	_ context.Context,
	_, _ string,
	timers []importer.Timer,
) (importer.Report, error) {
	f.timers = timers
	return f.timerReport, nil
}

type fakeChannelIntegrationsRepository struct {
	integration channelsintegrationsmodel.ChannelIntegration
}

func (f fakeChannelIntegrationsRepository) GetByChannelAndService(
	_ context.Context,
	_ string,
	_ integrationsmodel.Service,
) (channelsintegrationsmodel.ChannelIntegration, error) {
	return f.integration, nil
}

func (f fakeChannelIntegrationsRepository) Create(
	_ context.Context,
	_ channelsintegrations.CreateInput,
) (channelsintegrationsmodel.ChannelIntegration, error) {
	return channelsintegrationsmodel.ChannelIntegration{}, nil
}

func (f fakeChannelIntegrationsRepository) Update(
	_ context.Context,
	_ string,
	_ channelsintegrations.UpdateInput,
) error {
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}
