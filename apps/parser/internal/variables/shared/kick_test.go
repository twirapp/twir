package shared_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	sharedvars "github.com/twirapp/twir/apps/parser/internal/variables/shared"
	streamvars "github.com/twirapp/twir/apps/parser/internal/variables/stream"
	subscribervars "github.com/twirapp/twir/apps/parser/internal/variables/subscribers"
	"github.com/twirapp/twir/apps/parser/internal/types"
	parserservices "github.com/twirapp/twir/apps/parser/internal/types/services"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	cfg "github.com/twirapp/twir/libs/config"
	"go.uber.org/zap"
)

func TestHandlerByPlatform(t *testing.T) {
	handler := sharedvars.HandlerByPlatform(map[platformentity.Platform]types.VariableHandler{
		sharedvars.PlatformKick: func(
			ctx context.Context,
			parseCtx *types.VariableParseContext,
			variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			return &types.VariableHandlerResult{Result: "kick"}, nil
		},
	})

	res, err := handler(context.Background(), &types.VariableParseContext{ParseContext: &types.ParseContext{Platform: sharedvars.PlatformKick}}, &types.VariableData{})
	require.NoError(t, err)
	require.Equal(t, "kick", res.Result)

	res, err = handler(context.Background(), &types.VariableParseContext{ParseContext: &types.ParseContext{Platform: sharedvars.PlatformTwitch}}, &types.VariableData{})
	require.NoError(t, err)
	require.Equal(t, "not supported on this platform", res.Result)
}

func TestKickBackedVariables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
		case "/public/v1/channels":
			require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"active_subscribers_count":42,"broadcaster_user_id":123,"category":{"id":99,"name":"Slots","thumbnail":"thumb"},"channel_description":"desc","slug":"kick-channel","stream":{"is_live":false,"viewer_count":0},"stream_title":"Kick Title"}],"message":"OK"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	restore := sharedvars.SetKickClientOptionsForTests(server.Client(), server.URL, server.URL)
	t.Cleanup(restore)

	parseCtx := &types.VariableParseContext{ParseContext: &types.ParseContext{
		Platform: sharedvars.PlatformKick,
		Channel: &types.ParseContextChannel{
			ID:   "123",
			Name: "kick-channel",
		},
		Services: &parserservices.Services{
			Config: &cfg.Config{KickClientId: "client-id", KickClientSecret: "client-secret"},
			Logger: zap.NewNop(),
		},
	}}

	titleResult, err := streamvars.Title.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.title"})
	require.NoError(t, err)
	require.Equal(t, "Kick Title", titleResult.Result)

	categoryResult, err := streamvars.Category.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.category"})
	require.NoError(t, err)
	require.Equal(t, "Slots", categoryResult.Result)

	subsCountResult, err := subscribervars.Count.Handler(context.Background(), parseCtx, &types.VariableData{Key: "subscribers.count"})
	require.NoError(t, err)
	require.Equal(t, "42", subsCountResult.Result)

	t.Log(fmt.Sprintf("kick stream.title=%s", titleResult.Result))
	t.Log(fmt.Sprintf("kick stream.category=%s", categoryResult.Result))
	t.Log(fmt.Sprintf("kick subscribers.count=%s", subsCountResult.Result))
}
