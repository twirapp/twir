package shared_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
			_, _ = w.Write([]byte(`{"data":[{"active_subscribers_count":42,"broadcaster_user_id":123,"category":{"id":99,"name":"Slots","thumbnail":"thumb"},"channel_description":"desc","slug":"kick-channel","stream":{"is_live":true,"viewer_count":321,"language":"en","start_time":"` + time.Now().UTC().Add(-2*time.Hour).Format(time.RFC3339) + `","thumbnail":"kick-thumb","custom_tags":["tag1","tag2"]},"stream_title":"Kick Title"}],"message":"OK"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	restore := sharedvars.SetKickClientOptionsForTests(server.Client(), server.URL, server.URL)
	t.Cleanup(restore)
	restoreRequester := sharedvars.SetKickAppTokenRequesterForTests(func(ctx context.Context, parseCtx *types.VariableParseContext) (string, error) {
		return "test-token", nil
	})
	t.Cleanup(restoreRequester)

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

	viewersResult, err := streamvars.Viewers.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.viewers"})
	require.NoError(t, err)
	require.Equal(t, "321", viewersResult.Result)

	uptimeResult, err := streamvars.Uptime.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.uptime"})
	require.NoError(t, err)
	require.NotEqual(t, "", uptimeResult.Result)
	require.NotEqual(t, "offline", uptimeResult.Result)

	languageResult, err := streamvars.Language.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.language"})
	require.NoError(t, err)
	require.Equal(t, "en", languageResult.Result)

	slugResult, err := streamvars.Slug.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.slug"})
	require.NoError(t, err)
	require.Equal(t, "kick-channel", slugResult.Result)

	descriptionResult, err := streamvars.Description.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.description"})
	require.NoError(t, err)
	require.Equal(t, "desc", descriptionResult.Result)

	thumbnailResult, err := streamvars.Thumbnail.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.thumbnail"})
	require.NoError(t, err)
	require.Equal(t, "kick-thumb", thumbnailResult.Result)

	tagsResult, err := streamvars.Tags.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.tags"})
	require.NoError(t, err)
	require.Equal(t, "tag1, tag2", tagsResult.Result)

	t.Log(fmt.Sprintf("kick stream.title=%s", titleResult.Result))
	t.Log(fmt.Sprintf("kick stream.category=%s", categoryResult.Result))
	t.Log(fmt.Sprintf("kick subscribers.count=%s", subsCountResult.Result))
	t.Log(fmt.Sprintf("kick stream.viewers=%s", viewersResult.Result))
	t.Log(fmt.Sprintf("kick stream.uptime=%s", uptimeResult.Result))
	t.Log(fmt.Sprintf("kick stream.language=%s", languageResult.Result))
	t.Log(fmt.Sprintf("kick stream.slug=%s", slugResult.Result))
	t.Log(fmt.Sprintf("kick stream.description=%s", descriptionResult.Result))
	t.Log(fmt.Sprintf("kick stream.thumbnail=%s", thumbnailResult.Result))
	t.Log(fmt.Sprintf("kick stream.tags=%s", tagsResult.Result))
}
