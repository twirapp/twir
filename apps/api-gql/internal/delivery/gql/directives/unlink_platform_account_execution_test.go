package directives

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
)

func TestUnlinkPlatformAccountGraphQLRequiresOwnerOrBotAdmin(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	managerID := uuid.New()
	viewerID := uuid.New()
	unauthorizedUserID := uuid.New()
	botAdminID := uuid.New()

	tests := []struct {
		name                string
		user                model.Users
		roles               []model.ChannelRole
		wantError           string
		wantDisconnectCalls int
		wantBindingPresent  bool
	}{
		{
			name:               "denies stale selected dashboard before disconnecting",
			user:               model.Users{ID: unauthorizedUserID.String()},
			wantError:          "user does not have access to selected dashboard",
			wantBindingPresent: true,
		},
		{
			name: "denies selected-dashboard view collaborator before disconnecting",
			user: model.Users{ID: viewerID.String()},
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: viewerID.String()}},
				Permissions: pq.StringArray{"VIEW_BOT_SETTINGS"},
			}},
			wantError:          "user has no permission to access this resource",
			wantBindingPresent: true,
		},
		{
			name: "denies selected-dashboard manager from unlinking",
			user: model.Users{ID: managerID.String()},
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: managerID.String()}},
				Permissions: pq.StringArray{"MANAGE_BOT_SETTINGS"},
			}},
			wantError:          "only the channel owner or a bot admin can manage platform identities",
			wantBindingPresent: true,
		},
		{
			name:                "allows normalized binding owner to unlink",
			user:                model.Users{ID: ownerID.String()},
			wantDisconnectCalls: 1,
		},
		{
			name:                "allows bot admin to unlink",
			user:                model.Users{ID: botAdminID.String(), IsBotAdmin: true},
			wantDisconnectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations := &unlinkPlatformAccountOperations{bindingPresent: true}
			server := newUnlinkPlatformAccountGraphQLServer(t, dashboardID, ownerID, tt.user, tt.roles, operations)

			response := executeUnlinkPlatformAccount(t, server)
			if tt.wantError != "" {
				if len(response.Errors) != 1 || response.Errors[0].Message != tt.wantError {
					t.Fatalf("GraphQL errors = %#v, want %q", response.Errors, tt.wantError)
				}
			} else {
				if len(response.Errors) != 0 {
					t.Fatalf("GraphQL errors = %#v, want none", response.Errors)
				}
				if response.Data == nil || !response.Data.UnlinkPlatformAccount {
					t.Fatalf("GraphQL data = %#v, want unlinkPlatformAccount true", response.Data)
				}
			}

			if operations.disconnectCalls != tt.wantDisconnectCalls {
				t.Fatalf("Disconnect() calls = %d, want %d", operations.disconnectCalls, tt.wantDisconnectCalls)
			}
			if operations.bindingPresent != tt.wantBindingPresent {
				t.Fatalf("binding present = %t, want %t", operations.bindingPresent, tt.wantBindingPresent)
			}
			if tt.wantDisconnectCalls != 0 && (operations.dashboardID != dashboardID || operations.platform != platformentity.PlatformKick) {
				t.Fatalf("Disconnect() = (%s, %s), want (%s, %s)", operations.dashboardID, operations.platform, dashboardID, platformentity.PlatformKick)
			}
		})
	}
}

func newUnlinkPlatformAccountGraphQLServer(
	t *testing.T,
	dashboardID uuid.UUID,
	ownerID uuid.UUID,
	user model.Users,
	roles []model.ChannelRole,
	operations *unlinkPlatformAccountOperations,
) *handler.Server {
	t.Helper()
	sessionUser := &user
	dashboardAccess := dashboardaccess.New(
		selectedDashboardDirectiveChannelReader{channel: channelentity.Channel{
			ID: dashboardID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: ownerID, PlatformChannelID: "owner-channel", Enabled: true,
			}},
		}},
		&selectedDashboardDirectiveStore{roles: roles},
	)

	resolver, err := resolvers.New(resolvers.Deps{
		ChannelPlatformBindingsService: operations,
		ChannelPlatformDashboard:       unlinkPlatformAccountDashboard{dashboardID: dashboardID},
		CurrentPlatform:                unlinkPlatformAccountCurrentPlatform{},
		Sessions:                       channelPlatformBindingExecutionSession{user: sessionUser},
		DashboardAccess:                dashboardAccess,
	})
	if err != nil {
		t.Fatalf("create resolver: %v", err)
	}

	directive := &Directives{
		sessions: &selectedDashboardDirectiveSession{
			user:        sessionUser,
			dashboardID: dashboardID.String(),
		},
		dashboardAccess: dashboardAccess,
	}

	server := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			IsAuthenticated:                    directive.IsAuthenticated,
			HasAccessToSelectedDashboard:       directive.HasAccessToSelectedDashboard,
			HasChannelRolesDashboardPermission: directive.HasChannelRolesDashboardPermission,
		},
	}))
	server.AddTransport(transport.POST{})

	return server
}

func executeUnlinkPlatformAccount(t *testing.T, server *handler.Server) unlinkPlatformAccountResponse {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/query",
		bytes.NewBufferString(`{"query":"mutation { unlinkPlatformAccount(platform: \"kick\") }"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GraphQL status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response unlinkPlatformAccountResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode GraphQL response %q: %v", recorder.Body.String(), err)
	}

	return response
}

type unlinkPlatformAccountResponse struct {
	Data *struct {
		UnlinkPlatformAccount bool `json:"unlinkPlatformAccount"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type unlinkPlatformAccountDashboard struct {
	dashboardID uuid.UUID
}

func (d unlinkPlatformAccountDashboard) GetSelectedDashboard(context.Context) (string, error) {
	return d.dashboardID.String(), nil
}

type unlinkPlatformAccountCurrentPlatform struct{}

func (unlinkPlatformAccountCurrentPlatform) GetCurrentPlatform(context.Context) (string, error) {
	return platformentity.PlatformTwitch.String(), nil
}

type unlinkPlatformAccountOperations struct {
	disconnectCalls int
	dashboardID     uuid.UUID
	platform        platformentity.Platform
	bindingPresent  bool
}

func (*unlinkPlatformAccountOperations) List(context.Context, uuid.UUID) ([]channelplatformservice.Binding, error) {
	return nil, nil
}

func (*unlinkPlatformAccountOperations) Options() []channelplatformservice.Option {
	return nil
}

func (*unlinkPlatformAccountOperations) Connect(context.Context, uuid.UUID, platformentity.Platform) (string, error) {
	return "", nil
}

func (o *unlinkPlatformAccountOperations) Disconnect(_ context.Context, dashboardID uuid.UUID, platform platformentity.Platform) error {
	o.disconnectCalls++
	o.dashboardID = dashboardID
	o.platform = platform
	o.bindingPresent = false
	return nil
}

func (*unlinkPlatformAccountOperations) SetEnabled(context.Context, uuid.UUID, platformentity.Platform, bool) (channelplatformservice.Binding, error) {
	return channelplatformservice.Binding{}, nil
}
