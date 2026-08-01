package resolvers

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/uuid"
	nightbotintegration "github.com/twirapp/twir/apps/api-gql/internal/services/nightbot_integration"
	model "github.com/twirapp/twir/libs/gomodels"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestNightbotAuthLinkCreatesStateBoundToSelectedDashboardAndUser(t *testing.T) {
	dashboardID := uuid.New()
	userID := uuid.New()
	sessions := &nightbotOAuthSessionsFake{dashboardID: dashboardID.String(), user: &model.Users{ID: userID.String()}, state: "nightbot-state"}
	clientID := "nightbot-client"
	clientSecret := "nightbot-secret"
	redirectURL := "https://twir.test/dashboard/integrations/callbacks/nightbot"
	service := nightbotintegration.New(nightbotintegration.Opts{
		IntegrationsRepository: nightbotIntegrationsRepositoryFake{integration: integrationsmodel.Integration{
			ClientID:     &clientID,
			ClientSecret: &clientSecret,
			RedirectURL:  &redirectURL,
		}},
	})
	resolver := &queryResolver{Resolver: &Resolver{deps: Deps{
		Sessions:                   sessions,
		NightbotIntegrationService: service,
	}}}

	link, err := resolver.NightbotGetAuthLink(context.Background())
	if err != nil {
		t.Fatalf("NightbotGetAuthLink() error = %v", err)
	}
	parsedLink, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if got, want := parsedLink.Query().Get("state"), "nightbot-state"; got != want {
		t.Fatalf("Nightbot authorization URL state = %q, want %q", got, want)
	}
	if got, want := sessions.service, integrationsmodel.ServiceNightbot; got != want {
		t.Fatalf("OAuth attempt service = %q, want %q", got, want)
	}
	if got, want := sessions.channelID, dashboardID; got != want {
		t.Fatalf("OAuth attempt channel ID = %s, want %s", got, want)
	}
	if got, want := sessions.initiatorUserID, userID; got != want {
		t.Fatalf("OAuth attempt initiator user ID = %s, want %s", got, want)
	}
}

type nightbotOAuthSessionsFake struct {
	dashboardID     string
	user            *model.Users
	state           string
	service         integrationsmodel.Service
	channelID       uuid.UUID
	initiatorUserID uuid.UUID
}

func (f *nightbotOAuthSessionsFake) GetSelectedDashboard(context.Context) (string, error) {
	return f.dashboardID, nil
}

func (f *nightbotOAuthSessionsFake) GetAuthenticatedUserModel(context.Context) (*model.Users, error) {
	return f.user, nil
}

func (f *nightbotOAuthSessionsFake) GetCurrentPlatform(context.Context) (string, error) {
	return "", nil
}

func (f *nightbotOAuthSessionsFake) SetSessionSelectedDashboard(context.Context, string) error {
	return nil
}

func (f *nightbotOAuthSessionsFake) SessionLogout(context.Context) error {
	return nil
}

func (f *nightbotOAuthSessionsFake) CreateIntegrationOAuthAttempt(
	_ context.Context,
	service integrationsmodel.Service,
	channelID, initiatorUserID uuid.UUID,
) (string, error) {
	f.service = service
	f.channelID = channelID
	f.initiatorUserID = initiatorUserID
	return f.state, nil
}

type nightbotIntegrationsRepositoryFake struct {
	integration integrationsmodel.Integration
}

func (f nightbotIntegrationsRepositoryFake) GetByService(
	_ context.Context,
	_ integrationsmodel.Service,
) (integrationsmodel.Integration, error) {
	return f.integration, nil
}
