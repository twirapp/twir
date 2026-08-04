package resolvers

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/uuid"
	streamelementsservice "github.com/twirapp/twir/apps/api-gql/internal/services/streamelements"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestStreamElementsAuthLinkCreatesStateBoundToSelectedDashboardAndUser(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	userID := uuid.New()
	sessions := &streamElementsOAuthSessionsFake{
		dashboardID: dashboardID.String(),
		user:        &model.Users{ID: userID.String()},
		state:       "streamelements-state",
	}
	service, err := streamelementsservice.New(streamelementsservice.Opts{Config: config.Config{
		SiteBaseUrl: "https://twir.test", StreamElementsClientId: "client-id", StreamElementsClientSecret: "client-secret",
	}})
	if err != nil {
		t.Fatalf("create StreamElements service: %v", err)
	}
	resolver := &queryResolver{Resolver: &Resolver{deps: Deps{
		Sessions:              sessions,
		StreamElementsService: service,
	}}}

	link, err := resolver.StreamelementsGetAuthorizationURL(context.Background())
	if err != nil {
		t.Fatalf("StreamelementsGetAuthorizationURL() error = %v", err)
	}
	parsedLink, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	if got, want := parsedLink.Query().Get("state"), "streamelements-state"; got != want {
		t.Fatalf("authorization URL state = %q, want %q", got, want)
	}
	if got, want := sessions.service, integrationsmodel.ServiceStreamElements; got != want {
		t.Fatalf("OAuth attempt service = %q, want %q", got, want)
	}
	if got, want := sessions.channelID, dashboardID; got != want {
		t.Fatalf("OAuth attempt channel ID = %s, want %s", got, want)
	}
	if got, want := sessions.initiatorUserID, userID; got != want {
		t.Fatalf("OAuth attempt initiator user ID = %s, want %s", got, want)
	}
}

type streamElementsOAuthSessionsFake struct {
	dashboardID     string
	user            *model.Users
	state           string
	service         integrationsmodel.Service
	channelID       uuid.UUID
	initiatorUserID uuid.UUID
}

func (f *streamElementsOAuthSessionsFake) GetSelectedDashboard(context.Context) (string, error) {
	return f.dashboardID, nil
}

func (f *streamElementsOAuthSessionsFake) GetAuthenticatedUserModel(context.Context) (*model.Users, error) {
	return f.user, nil
}

func (f *streamElementsOAuthSessionsFake) GetCurrentPlatform(context.Context) (string, error) {
	return "", nil
}

func (f *streamElementsOAuthSessionsFake) SetSessionSelectedDashboard(context.Context, string) error {
	return nil
}

func (f *streamElementsOAuthSessionsFake) SessionLogout(context.Context) error {
	return nil
}

func (f *streamElementsOAuthSessionsFake) CreateIntegrationOAuthAttempt(
	_ context.Context,
	service integrationsmodel.Service,
	channelID, initiatorUserID uuid.UUID,
) (string, error) {
	f.service = service
	f.channelID = channelID
	f.initiatorUserID = initiatorUserID
	return f.state, nil
}
