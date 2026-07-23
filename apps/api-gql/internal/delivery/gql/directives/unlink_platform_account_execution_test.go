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
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	channelplatformservice "github.com/twirapp/twir/apps/api-gql/internal/services/channel_platforms"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestUnlinkPlatformAccountGraphQLRequiresSelectedDashboardAccess(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	unauthorizedUserID := uuid.New()

	tests := []struct {
		name                string
		userID              uuid.UUID
		wantError           string
		wantDisconnectCalls int
	}{
		{
			name:      "denies stale selected dashboard before disconnecting",
			userID:    unauthorizedUserID,
			wantError: "user does not have access to selected dashboard",
		},
		{
			name:                "allows normalized binding owner to unlink",
			userID:              ownerID,
			wantDisconnectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations := &unlinkPlatformAccountOperations{}
			server := newUnlinkPlatformAccountGraphQLServer(t, dashboardID, ownerID, tt.userID, operations)

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
	userID uuid.UUID,
	operations *unlinkPlatformAccountOperations,
) *handler.Server {
	t.Helper()

	resolver, err := resolvers.New(resolvers.Deps{
		ChannelPlatformBindingsService: operations,
		ChannelPlatformDashboard:       unlinkPlatformAccountDashboard{dashboardID: dashboardID},
		CurrentPlatform:                unlinkPlatformAccountCurrentPlatform{},
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
			selectedDashboardDirectiveChannelReader{channel: channelsmodel.Channel{
				ID: dashboardID,
				Bindings: []channelplatformsmodel.ChannelPlatform{{
					ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformTwitch, UserID: ownerID, PlatformChannelID: "owner-channel", Enabled: true,
				}},
			}},
			&selectedDashboardDirectiveStore{},
		),
	}

	server := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			IsAuthenticated:              directive.IsAuthenticated,
			HasAccessToSelectedDashboard: directive.HasAccessToSelectedDashboard,
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
}

func (*unlinkPlatformAccountOperations) List(context.Context, uuid.UUID) ([]channelplatformservice.Binding, error) {
	return nil, nil
}

func (*unlinkPlatformAccountOperations) Connect(context.Context, uuid.UUID, platformentity.Platform, string) (string, error) {
	return "", nil
}

func (o *unlinkPlatformAccountOperations) Disconnect(_ context.Context, dashboardID uuid.UUID, platform platformentity.Platform) error {
	o.disconnectCalls++
	o.dashboardID = dashboardID
	o.platform = platform
	return nil
}

func (*unlinkPlatformAccountOperations) SetEnabled(context.Context, uuid.UUID, platformentity.Platform, bool) (channelplatformservice.Binding, error) {
	return channelplatformservice.Binding{}, nil
}
