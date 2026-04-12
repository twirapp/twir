package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/twirapp/twir/apps/twitch-mock/internal/admin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/config"
	"github.com/twirapp/twir/apps/twitch-mock/internal/handlers"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
	twitchws "github.com/twirapp/twir/apps/twitch-mock/internal/websocket"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestConfig() *config.Config {
	return &config.Config{
		HTTPAddr:    ":0",
		WSAddr:      ":0",
		AdminAddr:   ":0",
		SiteBaseURL: "http://localhost:3005",
	}
}

func newHandlerServer(t *testing.T) (*httptest.Server, *state.State) {
	t.Helper()
	st := state.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := handlers.New(newTestConfig(), st, logger)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts, st
}

func newAdminServer(t *testing.T, st *state.State, ws *twitchws.Server) *httptest.Server {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := admin.New(st, logger, ws)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func newWSServer(t *testing.T, st *state.State) (*httptest.Server, *twitchws.Server) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	wsSrv := twitchws.New(logger, st)
	ts := httptest.NewServer(wsSrv.Handler())
	t.Cleanup(ts.Close)
	return ts, wsSrv
}

func httpGet(t *testing.T, rawURL string, headers map[string]string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http.Do GET %s: %v", rawURL, err)
	}
	return resp
}

func httpPost(t *testing.T, rawURL string, body any) *http.Response {
	t.Helper()
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}
		buf = bytes.NewReader(b)
	} else {
		buf = strings.NewReader("")
	}
	resp, err := http.Post(rawURL, "application/json", buf)
	if err != nil {
		t.Fatalf("http.Post %s: %v", rawURL, err)
	}
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		t.Fatalf("json.Decode: %v", err)
	}
}

func postForm(t *testing.T, rawURL string, values url.Values) *http.Response {
	t.Helper()
	resp, err := http.PostForm(rawURL, values)
	if err != nil {
		t.Fatalf("http.PostForm %s: %v", rawURL, err)
	}
	return resp
}

func TestHealth(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/health", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %v", body["status"])
	}
}

func TestTokenClientCredentials(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := postForm(t, ts.URL+"/oauth2/token", url.Values{
		"grant_type": {"client_credentials"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	decodeJSON(t, resp, &body)

	if body["access_token"] != config.MockAppToken {
		t.Fatalf("expected access_token=%q, got %v", config.MockAppToken, body["access_token"])
	}
	expIn, ok := body["expires_in"].(float64)
	if !ok || expIn != 99999999 {
		t.Fatalf("expected expires_in=99999999, got %v", body["expires_in"])
	}
}

func TestTokenAuthorizationCode(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := postForm(t, ts.URL+"/oauth2/token", url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"mock_code"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	decodeJSON(t, resp, &body)

	if body["access_token"] != config.MockUserToken {
		t.Fatalf("expected access_token=%q, got %v", config.MockUserToken, body["access_token"])
	}
	if body["refresh_token"] != "mock-user-refresh" {
		t.Fatalf("expected refresh_token=mock-user-refresh, got %v", body["refresh_token"])
	}
}

func TestTokenRefreshToken(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := postForm(t, ts.URL+"/oauth2/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"mock-user-refresh"},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["access_token"] != config.MockUserToken {
		t.Fatalf("expected user token, got %v", body["access_token"])
	}
}

func TestTokenUnsupportedGrantType(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := postForm(t, ts.URL+"/oauth2/token", url.Values{
		"grant_type": {"unknown"},
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestValidateBroadcasterToken(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/oauth2/validate", map[string]string{
		"Authorization": "OAuth " + config.MockUserToken,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	decodeJSON(t, resp, &body)

	if body["login"] != config.MockBroadcasterLogin {
		t.Fatalf("expected login=%q, got %v", config.MockBroadcasterLogin, body["login"])
	}
	if body["user_id"] != config.MockBroadcasterID {
		t.Fatalf("expected user_id=%q, got %v", config.MockBroadcasterID, body["user_id"])
	}
}

func TestValidateBotToken(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/oauth2/validate", map[string]string{
		"Authorization": "OAuth " + config.MockBotToken,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	decodeJSON(t, resp, &body)

	if body["login"] != config.MockBotLogin {
		t.Fatalf("expected login=%q, got %v", config.MockBotLogin, body["login"])
	}
	if body["user_id"] != config.MockBotID {
		t.Fatalf("expected user_id=%q, got %v", config.MockBotID, body["user_id"])
	}
}

func TestHelixUsersBroadcaster(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/users", map[string]string{
		"Authorization": "Bearer " + config.MockUserToken,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &body)

	if len(body.Data) == 0 {
		t.Fatal("expected at least one user")
	}
	if body.Data[0]["login"] != config.MockBroadcasterLogin {
		t.Fatalf("expected login=%q, got %v", config.MockBroadcasterLogin, body.Data[0]["login"])
	}
}

func TestHelixUsersBot(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/users", map[string]string{
		"Authorization": "Bearer " + config.MockBotToken,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &body)

	if len(body.Data) == 0 {
		t.Fatal("expected at least one user")
	}
	if body.Data[0]["login"] != config.MockBotLogin {
		t.Fatalf("expected login=%q, got %v", config.MockBotLogin, body.Data[0]["login"])
	}
}

func TestHelixUsersByLogin(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/users?login="+config.MockBotLogin, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &body)
	if len(body.Data) == 0 {
		t.Fatal("expected user data")
	}
}

func TestHelixUsersByID(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/users?id="+config.MockBotID, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHelixStreamsOffline(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/streams", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []any `json:"data"`
	}
	decodeJSON(t, resp, &body)
	if len(body.Data) != 0 {
		t.Fatalf("expected empty data for offline stream, got %d", len(body.Data))
	}
}

func TestHelixStreamsOnline(t *testing.T) {
	ts, st := newHandlerServer(t)
	st.SetStreamOnline(true)

	resp := httpGet(t, ts.URL+"/helix/streams", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &body)
	if len(body.Data) == 0 {
		t.Fatal("expected stream data when online")
	}
	if body.Data[0]["type"] != "live" {
		t.Fatalf("expected type=live, got %v", body.Data[0]["type"])
	}
}

func TestHelixChannels(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/channels", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &body)
	if len(body.Data) == 0 {
		t.Fatal("expected channel data")
	}
	if body.Data[0]["broadcaster_id"] != config.MockBroadcasterID {
		t.Fatalf("expected broadcaster_id=%q, got %v", config.MockBroadcasterID, body.Data[0]["broadcaster_id"])
	}
}

func TestConduitsCreateAndList(t *testing.T) {
	ts, _ := newHandlerServer(t)

	createResp := httpPost(t, ts.URL+"/helix/eventsub/conduits", map[string]any{
		"shard_count": 1,
	})
	if createResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on create, got %d", createResp.StatusCode)
	}

	var createBody struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, createResp, &createBody)
	if len(createBody.Data) == 0 {
		t.Fatal("expected conduit in create response")
	}
	conduitID, _ := createBody.Data[0]["id"].(string)
	if conduitID == "" {
		t.Fatal("expected non-empty conduit id")
	}

	listResp := httpGet(t, ts.URL+"/helix/eventsub/conduits", nil)
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on list, got %d", listResp.StatusCode)
	}

	var listBody struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, listResp, &listBody)
	if len(listBody.Data) == 0 {
		t.Fatal("expected at least one conduit in list")
	}
	found := false
	for _, c := range listBody.Data {
		if c["id"] == conduitID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created conduit %q not found in list", conduitID)
	}
}

func TestConduitsUpdateShards(t *testing.T) {
	ts, _ := newHandlerServer(t)

	createResp := httpPost(t, ts.URL+"/helix/eventsub/conduits", map[string]any{
		"shard_count": 1,
	})
	var createBody struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, createResp, &createBody)
	conduitID := createBody.Data[0]["id"].(string)

	updatePayload := map[string]any{
		"conduit_id": conduitID,
		"shards": []map[string]any{
			{
				"id": 0,
				"transport": map[string]any{
					"method":     "websocket",
					"session_id": "test-session-id",
				},
			},
		},
	}

	reqBody, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/helix/eventsub/conduits/shards", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on update shards, got %d", resp.StatusCode)
	}
}

func TestConduitsUpdateShardsNotFound(t *testing.T) {
	ts, _ := newHandlerServer(t)

	updatePayload := map[string]any{
		"conduit_id": "nonexistent-id",
		"shards":     []map[string]any{},
	}
	reqBody, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/helix/eventsub/conduits/shards", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for nonexistent conduit, got %d", resp.StatusCode)
	}
}

func TestSubscriptionCreate(t *testing.T) {
	ts, _ := newHandlerServer(t)

	body := map[string]any{
		"type":      "channel.follow",
		"version":   "2",
		"condition": map[string]any{"broadcaster_user_id": config.MockBroadcasterID},
		"transport": map[string]any{
			"method":     "websocket",
			"session_id": "test-session",
		},
	}

	resp := httpPost(t, ts.URL+"/helix/eventsub/subscriptions", body)
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}

	var result struct {
		Data  []map[string]any `json:"data"`
		Total float64          `json:"total"`
	}
	decodeJSON(t, resp, &result)
	if len(result.Data) == 0 {
		t.Fatal("expected subscription in response")
	}
	if result.Data[0]["type"] != "channel.follow" {
		t.Fatalf("expected type=channel.follow, got %v", result.Data[0]["type"])
	}
}

func TestSubscriptionDelete(t *testing.T) {
	ts, st := newHandlerServer(t)

	sub := st.CreateSubscription("channel.follow", "2", nil, state.SubscriptionTransport{Method: "websocket"})

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/helix/eventsub/subscriptions?id="+sub.ID, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}

func TestSubscriptionDeleteByPath(t *testing.T) {
	ts, st := newHandlerServer(t)

	sub := st.CreateSubscription("channel.follow", "2", nil, state.SubscriptionTransport{Method: "websocket"})

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/helix/eventsub/subscriptions/"+sub.ID, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}

func TestSubscriptionDeleteMissingID(t *testing.T) {
	ts, _ := newHandlerServer(t)

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/helix/eventsub/subscriptions", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHelixUnknownEndpoint(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/helix/something/unknown", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		Data []any `json:"data"`
	}
	decodeJSON(t, resp, &body)
	if body.Data == nil {
		t.Fatal("expected data field to be present")
	}
}

func TestHelixChatMessages(t *testing.T) {
	ts, _ := newHandlerServer(t)

	body := map[string]any{
		"broadcaster_id": config.MockBroadcasterID,
		"sender_id":      config.MockBotID,
		"message":        "hello world",
	}
	resp := httpPost(t, ts.URL+"/helix/chat/messages", body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &result)
	if len(result.Data) == 0 {
		t.Fatal("expected data in response")
	}
}

func TestHelixModerationBans(t *testing.T) {
	ts, _ := newHandlerServer(t)

	body := map[string]any{
		"data": map[string]any{
			"user_id": "11111",
		},
	}
	resp := httpPost(t, ts.URL+"/helix/moderation/bans", body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		Data []map[string]any `json:"data"`
	}
	decodeJSON(t, resp, &result)
	if len(result.Data) == 0 {
		t.Fatal("expected ban data")
	}
	if result.Data[0]["user_id"] != "11111" {
		t.Fatalf("expected user_id=11111, got %v", result.Data[0]["user_id"])
	}
}

func TestHelixModerationBansNoBody(t *testing.T) {
	ts, _ := newHandlerServer(t)

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/helix/moderation/bans", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestNonHelixNotFound(t *testing.T) {
	ts, _ := newHandlerServer(t)

	resp := httpGet(t, ts.URL+"/completely/unknown/path", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestAuthorizeRedirect(t *testing.T) {
	ts, _ := newHandlerServer(t)

	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/oauth2/authorize?state=test123&response_type=code", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET authorize: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302, got %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if !strings.Contains(location, "code=mock_code") {
		t.Fatalf("expected code=mock_code in Location, got %q", location)
	}
	if !strings.Contains(location, "state=test123") {
		t.Fatalf("expected state=test123 in Location, got %q", location)
	}
}

func TestWebSocketSessionWelcome(t *testing.T) {
	st := state.New()
	wsTS, _ := newWSServer(t, st)

	wsURL := "ws" + strings.TrimPrefix(wsTS.URL, "http") + "/ws"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("websocket Dial: %v", err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	var msg map[string]any
	if err := conn.ReadJSON(&msg); err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}

	meta, ok := msg["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("expected metadata field, got %v", msg)
	}
	if meta["message_type"] != "session_welcome" {
		t.Fatalf("expected message_type=session_welcome, got %v", meta["message_type"])
	}

	payload, ok := msg["payload"].(map[string]any)
	if !ok {
		t.Fatalf("expected payload field")
	}
	session, ok := payload["session"].(map[string]any)
	if !ok {
		t.Fatalf("expected session in payload")
	}
	if session["status"] != "connected" {
		t.Fatalf("expected status=connected, got %v", session["status"])
	}
	if session["id"] == "" {
		t.Fatal("expected non-empty session id")
	}
}

func TestAdminTrigger(t *testing.T) {
	st := state.New()
	_, wsSrv := newWSServer(t, st)
	adminTS := newAdminServer(t, st, wsSrv)

	body := map[string]any{
		"from_user_id": "11111",
		"to_user_id":   config.MockBroadcasterID,
	}
	resp := httpPost(t, adminTS.URL+"/admin/trigger/channel.follow", body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	decodeJSON(t, resp, &result)
	if result["status"] != "ok" {
		t.Fatalf("expected status=ok, got %v", result["status"])
	}
	if result["event_type"] != "channel.follow" {
		t.Fatalf("expected event_type=channel.follow, got %v", result["event_type"])
	}
}

func TestAdminTriggerStreamOnline(t *testing.T) {
	st := state.New()
	_, wsSrv := newWSServer(t, st)
	adminTS := newAdminServer(t, st, wsSrv)

	resp := httpPost(t, adminTS.URL+"/admin/trigger/stream.online", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	_, online := st.StreamSnapshot()
	if !online {
		t.Fatal("expected stream to be online after trigger")
	}
}

func TestAdminTriggerStreamOffline(t *testing.T) {
	st := state.New()
	st.SetStreamOnline(true)
	_, wsSrv := newWSServer(t, st)
	adminTS := newAdminServer(t, st, wsSrv)

	resp := httpPost(t, adminTS.URL+"/admin/trigger/stream.offline", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	_, online := st.StreamSnapshot()
	if online {
		t.Fatal("expected stream to be offline after trigger")
	}
}

func TestAdminTriggerWithConnectedClient(t *testing.T) {
	st := state.New()
	wsTS, wsSrv := newWSServer(t, st)
	adminTS := newAdminServer(t, st, wsSrv)

	wsURL := "ws" + strings.TrimPrefix(wsTS.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("websocket Dial: %v", err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var welcome map[string]any
	if err := conn.ReadJSON(&welcome); err != nil {
		t.Fatalf("ReadJSON welcome: %v", err)
	}

	triggerBody := map[string]any{"from_user_id": "11111", "to_user_id": config.MockBroadcasterID}
	triggerResp := httpPost(t, adminTS.URL+"/admin/trigger/channel.follow", triggerBody)
	if triggerResp.StatusCode != http.StatusOK {
		t.Fatalf("admin trigger expected 200, got %d", triggerResp.StatusCode)
	}

	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var notification map[string]any
	if err := conn.ReadJSON(&notification); err != nil {
		t.Fatalf("ReadJSON notification: %v", err)
	}

	meta, ok := notification["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("expected metadata in notification, got %v", notification)
	}
	if meta["message_type"] != "notification" {
		t.Fatalf("expected message_type=notification, got %v", meta["message_type"])
	}
	if meta["subscription_type"] != "channel.follow" {
		t.Fatalf("expected subscription_type=channel.follow, got %v", meta["subscription_type"])
	}
}

func TestStateConduitCRUD(t *testing.T) {
	st := state.New()

	conduit := st.CreateConduit(2)
	if conduit.ID == "" {
		t.Fatal("expected non-empty conduit ID")
	}
	if conduit.ShardCount != 2 {
		t.Fatalf("expected shard_count=2, got %d", conduit.ShardCount)
	}

	conduits := st.ListConduits()
	if len(conduits) != 1 {
		t.Fatalf("expected 1 conduit, got %d", len(conduits))
	}
}

func TestStateConduitUpdateShards(t *testing.T) {
	st := state.New()
	conduit := st.CreateConduit(1)

	err := st.UpdateConduitShards(conduit.ID, []state.ConduitShard{
		{ID: 0, Transport: state.ConduitShardTransport{Method: "websocket", SessionID: "sess-1"}},
	})
	if err != nil {
		t.Fatalf("UpdateConduitShards: %v", err)
	}
}

func TestStateConduitUpdateShardsNotFound(t *testing.T) {
	st := state.New()
	err := st.UpdateConduitShards("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent conduit")
	}
}

func TestStateSubscriptionCRUD(t *testing.T) {
	st := state.New()

	sub := st.CreateSubscription("channel.follow", "2", nil, state.SubscriptionTransport{Method: "websocket"})
	if sub.ID == "" {
		t.Fatal("expected non-empty subscription ID")
	}
	if sub.Status != "enabled" {
		t.Fatalf("expected status=enabled, got %q", sub.Status)
	}

	found, ok := st.FindSubscriptionByType("channel.follow")
	if !ok {
		t.Fatal("expected to find subscription by type")
	}
	if found.ID != sub.ID {
		t.Fatalf("expected sub ID %q, got %q", sub.ID, found.ID)
	}

	st.DeleteSubscription(sub.ID)

	_, ok = st.FindSubscriptionByType("channel.follow")
	if ok {
		t.Fatal("expected subscription to be deleted")
	}
}

func TestStateStreamSnapshot(t *testing.T) {
	st := state.New()

	_, online := st.StreamSnapshot()
	if online {
		t.Fatal("expected stream offline initially")
	}

	st.SetStreamOnline(true)
	data, online := st.StreamSnapshot()
	if !online {
		t.Fatal("expected stream online")
	}
	if data["type"] != "live" {
		t.Fatalf("expected type=live, got %v", data["type"])
	}

	st.SetStreamOnline(false)
	_, online = st.StreamSnapshot()
	if online {
		t.Fatal("expected stream offline after SetStreamOnline(false)")
	}
}

func TestAdminIndex(t *testing.T) {
	st := state.New()
	_, wsSrv := newWSServer(t, st)
	adminTS := newAdminServer(t, st, wsSrv)

	resp := httpGet(t, adminTS.URL+"/admin", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for admin index, got %d", resp.StatusCode)
	}
}
