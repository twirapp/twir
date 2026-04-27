package shared

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var (
	kickAppTokenMu      sync.Mutex
	kickAppToken        string
	kickAppTokenExpires time.Time

	kickHTTPClient  *http.Client
	kickAPIBaseURL  string
	kickAuthBaseURL string
)

func SetKickClientOptionsForTests(httpClient *http.Client, apiBaseURL, authBaseURL string) func() {
	kickAppTokenMu.Lock()
	oldHTTPClient := kickHTTPClient
	oldAPIBaseURL := kickAPIBaseURL
	oldAuthBaseURL := kickAuthBaseURL
	oldToken := kickAppToken
	oldExpires := kickAppTokenExpires

	kickHTTPClient = httpClient
	kickAPIBaseURL = apiBaseURL
	kickAuthBaseURL = authBaseURL
	kickAppToken = ""
	kickAppTokenExpires = time.Time{}
	kickAppTokenMu.Unlock()

	return func() {
		kickAppTokenMu.Lock()
		defer kickAppTokenMu.Unlock()

		kickHTTPClient = oldHTTPClient
		kickAPIBaseURL = oldAPIBaseURL
		kickAuthBaseURL = oldAuthBaseURL
		kickAppToken = oldToken
		kickAppTokenExpires = oldExpires
	}
}

func getKickAppAccessToken(ctx context.Context, parseCtx *types.VariableParseContext) (string, error) {
	kickAppTokenMu.Lock()
	defer kickAppTokenMu.Unlock()

	if kickAppToken != "" && time.Now().Before(kickAppTokenExpires) {
		return kickAppToken, nil
	}

	client, err := gokick.NewClient(&gokick.ClientOptions{
		ClientID:     parseCtx.Services.Config.KickClientId,
		ClientSecret: parseCtx.Services.Config.KickClientSecret,
		HTTPClient:   kickHTTPClient,
		APIBaseURL:   kickAPIBaseURL,
		AuthBaseURL:  kickAuthBaseURL,
	})
	if err != nil {
		return "", fmt.Errorf("create kick auth client: %w", err)
	}

	resp, err := client.GetAppAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get kick app access token: %w", err)
	}

	kickAppToken = resp.AccessToken
	kickAppTokenExpires = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second).Add(-5 * time.Minute)

	return kickAppToken, nil
}

func newKickAppClient(ctx context.Context, parseCtx *types.VariableParseContext) (*gokick.Client, error) {
	accessToken, err := getKickAppAccessToken(ctx, parseCtx)
	if err != nil {
		return nil, err
	}

	client, err := gokick.NewClient(&gokick.ClientOptions{
		AppAccessToken: accessToken,
		HTTPClient:     kickHTTPClient,
		APIBaseURL:     kickAPIBaseURL,
		AuthBaseURL:    kickAuthBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("create kick app client: %w", err)
	}

	return client, nil
}

func GetKickChannel(ctx context.Context, parseCtx *types.VariableParseContext) (*gokick.ChannelResponse, error) {
	client, err := newKickAppClient(ctx, parseCtx)
	if err != nil {
		return nil, err
	}

	filter := gokick.NewChannelListFilter()

	if broadcasterUserID, err := strconv.Atoi(parseCtx.Channel.ID); err == nil {
		filter = filter.SetBroadcasterUserIDs([]int{broadcasterUserID})
	} else if parseCtx.Channel.Name != "" {
		filter = filter.SetSlug([]string{parseCtx.Channel.Name})
	} else {
		return nil, fmt.Errorf("kick channel has neither numeric id nor slug")
	}

	resp, err := client.GetChannels(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get kick channels: %w", err)
	}

	if len(resp.Result) == 0 {
		return nil, nil
	}

	return &resp.Result[0], nil
}
