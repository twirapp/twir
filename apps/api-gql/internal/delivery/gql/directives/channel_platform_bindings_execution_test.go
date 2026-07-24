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

func TestChannelPlatformOptionsGraphQLRequiresViewBotSettings(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	viewerID := uuid.New()
	deniedID := uuid.New()

	tests := []struct {
		name             string
		userID           uuid.UUID
		roles            []model.ChannelRole
		wantError        string
		wantOptionsCalls int
	}{
		{
			name:             "allows normalized owner",
			userID:           ownerID,
			wantOptionsCalls: 1,
		},
		{
			name:   "allows assigned view bot settings role",
			userID: viewerID,
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: viewerID.String()}},
				Permissions: pq.StringArray{"VIEW_BOT_SETTINGS"},
			}},
			wantOptionsCalls: 1,
		},
		{
			name:   "denies role without view bot settings permission",
			userID: deniedID,
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: deniedID.String()}},
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			wantError: "user has no permission to access this resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations := &channelPlatformBindingExecutionOperations{}
			server := newChannelPlatformBindingExecutionServer(t, dashboardID, ownerID, tt.userID, tt.roles, operations)

			response := executeChannelPlatformBindingGraphQL(t, server, `{ channelPlatformOptions { platform capabilities { name } } }`)
			if tt.wantError != "" {
				if len(response.Errors) != 1 || response.Errors[0].Message != tt.wantError {
					t.Fatalf("GraphQL errors = %#v, want %q", response.Errors, tt.wantError)
				}
			} else {
				if len(response.Errors) != 0 {
					t.Fatalf("GraphQL errors = %#v, want none", response.Errors)
				}
				if response.Data == nil {
					t.Fatal("GraphQL data is nil")
				}
			}
			if operations.optionsCalls != tt.wantOptionsCalls {
				t.Fatalf("Options() calls = %d, want %d", operations.optionsCalls, tt.wantOptionsCalls)
			}
		})
	}
}

func TestChannelPlatformConnectGraphQLRequiresManageBotSettings(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	managerID := uuid.New()
	viewerID := uuid.New()

	tests := []struct {
		name             string
		userID           uuid.UUID
		roles            []model.ChannelRole
		wantError        string
		wantConnectCalls int
	}{
		{
			name:   "allows assigned manage bot settings role",
			userID: managerID,
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: managerID.String()}},
				Permissions: pq.StringArray{"MANAGE_BOT_SETTINGS"},
			}},
			wantConnectCalls: 1,
		},
		{
			name:   "denies view-only bot settings role",
			userID: viewerID,
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: viewerID.String()}},
				Permissions: pq.StringArray{"VIEW_BOT_SETTINGS"},
			}},
			wantError: "user has no permission to access this resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations := &channelPlatformBindingExecutionOperations{}
			server := newChannelPlatformBindingExecutionServer(t, dashboardID, ownerID, tt.userID, tt.roles, operations)

			response := executeChannelPlatformBindingGraphQL(t, server, `mutation { channelPlatformConnect(platform: KICK) }`)
			if tt.wantError != "" {
				if len(response.Errors) != 1 || response.Errors[0].Message != tt.wantError {
					t.Fatalf("GraphQL errors = %#v, want %q", response.Errors, tt.wantError)
				}
			} else if len(response.Errors) != 0 {
				t.Fatalf("GraphQL errors = %#v, want none", response.Errors)
			}
			if operations.connectCalls != tt.wantConnectCalls {
				t.Fatalf("Connect() calls = %d, want %d", operations.connectCalls, tt.wantConnectCalls)
			}
		})
	}
}

func newChannelPlatformBindingExecutionServer(
	t *testing.T,
	dashboardID uuid.UUID,
	ownerID uuid.UUID,
	userID uuid.UUID,
	roles []model.ChannelRole,
	operations *channelPlatformBindingExecutionOperations,
) *handler.Server {
	t.Helper()

	resolver, err := resolvers.New(resolvers.Deps{
		ChannelPlatformBindingsService: operations,
		ChannelPlatformDashboard:       channelPlatformBindingExecutionDashboard{dashboardID: dashboardID},
	})
	if err != nil {
		t.Fatalf("create resolver: %v", err)
	}

	directive := &Directives{
		sessions: &selectedDashboardDirectiveSession{
			user:        &model.Users{ID: userID.String()},
			dashboardID: dashboardID.String(),
		},
		dashboardAccess: dashboardaccess.New(
			selectedDashboardDirectiveChannelReader{channel: channelentity.Channel{
				ID: dashboardID,
				Bindings: []channelplatformentity.ChannelPlatform{{
					ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: ownerID, PlatformChannelID: "owner-channel", Enabled: true,
				}},
			}},
			&selectedDashboardDirectiveStore{roles: roles},
		),
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

func executeChannelPlatformBindingGraphQL(t *testing.T, server *handler.Server, query string) channelPlatformBindingExecutionResponse {
	t.Helper()

	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatalf("encode GraphQL query: %v", err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	server.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GraphQL status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response channelPlatformBindingExecutionResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode GraphQL response %q: %v", recorder.Body.String(), err)
	}
	return response
}

type channelPlatformBindingExecutionResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type channelPlatformBindingExecutionDashboard struct {
	dashboardID uuid.UUID
}

func (d channelPlatformBindingExecutionDashboard) GetSelectedDashboard(context.Context) (string, error) {
	return d.dashboardID.String(), nil
}

type channelPlatformBindingExecutionOperations struct {
	optionsCalls int
	connectCalls int
}

func (*channelPlatformBindingExecutionOperations) List(context.Context, uuid.UUID) ([]channelplatformservice.Binding, error) {
	return nil, nil
}

func (o *channelPlatformBindingExecutionOperations) Options() []channelplatformservice.Option {
	o.optionsCalls++
	return []channelplatformservice.Option{{
		Platform:     platformentity.PlatformTwitch,
		Capabilities: platformentity.PlatformTwitch.Capabilities(),
	}}
}

func (o *channelPlatformBindingExecutionOperations) Connect(context.Context, uuid.UUID, platformentity.Platform) (string, error) {
	o.connectCalls++
	return "https://provider.example/authorize", nil
}

func (*channelPlatformBindingExecutionOperations) Disconnect(context.Context, uuid.UUID, platformentity.Platform) error {
	return nil
}

func (*channelPlatformBindingExecutionOperations) SetEnabled(context.Context, uuid.UUID, platformentity.Platform, bool) (channelplatformservice.Binding, error) {
	return channelplatformservice.Binding{}, nil
}
