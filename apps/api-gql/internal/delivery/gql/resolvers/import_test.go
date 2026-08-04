package resolvers

import (
	"context"
	"errors"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlerrors"
	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	streamlabsintegration "github.com/twirapp/twir/apps/api-gql/internal/services/streamlabs_integration"
	cfg "github.com/twirapp/twir/libs/config"
	apperrors "github.com/twirapp/twir/libs/errors"
	model "github.com/twirapp/twir/libs/gomodels"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestImportSchemaUsesDashboardPermissionBoundariesAndRemovesRawExchange(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile("../schema/import.graphql")
	if err != nil {
		t.Fatalf("read import schema: %v", err)
	}
	schema := string(contents)
	for _, expected := range []string{
		"streamelementsGetData: StreamElementsIntegration @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: VIEW_INTEGRATIONS)",
		"streamelementsGetAuthorizationUrl: String! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_INTEGRATIONS)",
		"streamelementsPostCode(input: IntegrationOAuthCodeInput!): Boolean! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_INTEGRATIONS)",
		"streamelementsImportCommands: ImportReport! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_COMMANDS)",
		"streamelementsImportTimers: ImportReport! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_TIMERS)",
		"nightbotPostCode(input: IntegrationOAuthCodeInput!): Boolean! @isAuthenticated @hasAccessToSelectedDashboard @hasChannelRolesDashboardPermission(permission: MANAGE_INTEGRATIONS)",
	} {
		if !strings.Contains(schema, expected) {
			t.Fatalf("import schema does not contain permission-bound field %q", expected)
		}
	}
	if strings.Contains(schema, "streamelementsExchangeDataByCode") {
		t.Fatal("legacy raw StreamElements exchange remains exposed")
	}
}

func TestStreamlabsAuthLinkCreatesProviderBoundState(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	userID := uuid.New()
	sessions := &integrationOAuthSessionsFake{
		dashboardID: dashboardID.String(),
		user:        &model.Users{ID: userID.String()},
	}
	service := streamlabsintegration.New(streamlabsintegration.Opts{Config: cfg.Config{
		SiteBaseUrl:            "https://twir.test",
		StreamlabsClientId:     "client-id",
		StreamlabsClientSecret: "client-secret",
	}})
	resolver := &queryResolver{Resolver: &Resolver{deps: Deps{
		Sessions:                     sessions,
		StreamlabsIntegrationService: service,
	}}}

	link, err := resolver.StreamlabsAuthLink(context.Background())
	if err != nil {
		t.Fatalf("StreamlabsAuthLink() error = %v", err)
	}
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse Streamlabs authorization URL: %v", err)
	}
	if got, want := parsed.Query().Get("state"), "created-state"; got != want {
		t.Fatalf("Streamlabs state = %q, want %q", got, want)
	}
	if got, want := sessions.createdService, integrationsmodel.ServiceStreamLabs; got != want {
		t.Fatalf("created OAuth service = %q, want %q", got, want)
	}
	if got, want := sessions.createdChannelID, dashboardID; got != want {
		t.Fatalf("created OAuth channel ID = %s, want %s", got, want)
	}
	if got, want := sessions.createdUserID, userID; got != want {
		t.Fatalf("created OAuth user ID = %s, want %s", got, want)
	}
}

func TestCompleteIntegrationOAuthConsumesStateBeforeCodeExchange(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	userID := uuid.New()
	sessions := &integrationOAuthSessionsFake{
		dashboardID: dashboardID.String(),
		user:        &model.Users{ID: userID.String()},
	}
	postCalled := false

	err := completeIntegrationOAuth(
		context.Background(),
		sessions,
		integrationsmodel.ServiceStreamElements,
		"authorization-code",
		"single-use-state",
		func(_ context.Context, channelID, code string) error {
			postCalled = true
			if !sessions.consumed {
				t.Fatal("provider exchange ran before OAuth state consumption")
			}
			if channelID != dashboardID.String() || code != "authorization-code" {
				t.Fatalf("provider exchange args = (%q, %q)", channelID, code)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("completeIntegrationOAuth() error = %v", err)
	}
	if !postCalled {
		t.Fatal("provider exchange was not called")
	}
	if got, want := sessions.consumedService, integrationsmodel.ServiceStreamElements; got != want {
		t.Fatalf("consumed service = %q, want %q", got, want)
	}
	if got, want := sessions.consumedState, "single-use-state"; got != want {
		t.Fatalf("consumed state = %q, want %q", got, want)
	}
	if got, want := sessions.consumedChannelID, dashboardID; got != want {
		t.Fatalf("consumed channel ID = %s, want %s", got, want)
	}
	if got, want := sessions.consumedUserID, userID; got != want {
		t.Fatalf("consumed user ID = %s, want %s", got, want)
	}
}

func TestCompleteIntegrationOAuthRejectsConsumedOrMismatchedStateBeforeCodeExchange(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name string
		err  error
	}{
		{name: "replay", err: authsessions.ErrOAuthAttemptNotFound},
		{name: "wrong provider", err: authsessions.ErrOAuthAttemptMismatch},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			sessions := &integrationOAuthSessionsFake{
				dashboardID: uuid.NewString(),
				user:        &model.Users{ID: uuid.NewString()},
				consumeErr:  testCase.err,
			}
			postCalled := false

			err := completeIntegrationOAuth(
				context.Background(), sessions, integrationsmodel.ServiceNightbot, "code", "state",
				func(context.Context, string, string) error {
					postCalled = true
					return nil
				},
			)
			if err == nil || !strings.Contains(err.Error(), "invalid or expired") {
				t.Fatalf("completeIntegrationOAuth() error = %v, want sanitized OAuth rejection", err)
			}
			handled := gqlerrors.HandleError(err)
			var gqlErr *gqlerror.Error
			if !errors.As(handled, &gqlErr) {
				t.Fatalf("handled error type = %T, want *gqlerror.Error", handled)
			}
			if got, want := gqlErr.Extensions["code"], string(apperrors.ErrorCodeBadRequest); got != want {
				t.Fatalf("GraphQL error code = %#v, want %q", got, want)
			}
			if postCalled {
				t.Fatal("provider exchange ran after rejected OAuth state")
			}
		})
	}
}

func TestMapImportReportPreservesStableFailureReasons(t *testing.T) {
	t.Parallel()

	report := importer.Report{
		ImportedCount: 2,
		FailedCount:   2,
		Failures: []importer.Failure{
			{Name: "duplicate", Reason: importer.FailureDuplicate},
			{Name: "role", Reason: importer.FailureUnsupportedRole},
		},
	}

	got := mapImportReport(report)
	if got.ImportedCount != 2 || got.FailedCount != 2 {
		t.Fatalf("mapped counts = (%d, %d), want (2, 2)", got.ImportedCount, got.FailedCount)
	}
	wantNames := []string{"duplicate", "role"}
	gotNames := []string{got.Failures[0].Name, got.Failures[1].Name}
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("mapped failure names = %#v, want %#v", gotNames, wantNames)
	}
	if got.Failures[0].Reason.String() != string(importer.FailureDuplicate) ||
		got.Failures[1].Reason.String() != string(importer.FailureUnsupportedRole) {
		t.Fatalf("mapped failure reasons = (%s, %s)", got.Failures[0].Reason, got.Failures[1].Reason)
	}
}

type integrationOAuthSessionsFake struct {
	dashboardID string
	user        *model.Users
	consumeErr  error

	createdService   integrationsmodel.Service
	createdChannelID uuid.UUID
	createdUserID    uuid.UUID

	consumed          bool
	consumedState     string
	consumedService   integrationsmodel.Service
	consumedChannelID uuid.UUID
	consumedUserID    uuid.UUID
}

func (f *integrationOAuthSessionsFake) GetSelectedDashboard(context.Context) (string, error) {
	return f.dashboardID, nil
}

func (f *integrationOAuthSessionsFake) GetAuthenticatedUserModel(context.Context) (*model.Users, error) {
	return f.user, nil
}

func (f *integrationOAuthSessionsFake) CreateIntegrationOAuthAttempt(
	_ context.Context,
	service integrationsmodel.Service,
	channelID uuid.UUID,
	userID uuid.UUID,
) (string, error) {
	f.createdService = service
	f.createdChannelID = channelID
	f.createdUserID = userID
	return "created-state", nil
}

func (f *integrationOAuthSessionsFake) GetCurrentPlatform(context.Context) (string, error) {
	return "", nil
}

func (f *integrationOAuthSessionsFake) SetSessionSelectedDashboard(context.Context, string) error {
	return nil
}

func (f *integrationOAuthSessionsFake) SessionLogout(context.Context) error {
	return nil
}

func (f *integrationOAuthSessionsFake) ConsumeIntegrationOAuthAttempt(
	_ context.Context,
	state string,
	service integrationsmodel.Service,
	channelID, userID uuid.UUID,
	_ time.Time,
) error {
	f.consumed = true
	f.consumedState = state
	f.consumedService = service
	f.consumedChannelID = channelID
	f.consumedUserID = userID
	return f.consumeErr
}
