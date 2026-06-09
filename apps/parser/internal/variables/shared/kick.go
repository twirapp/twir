package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/apps/parser/internal/types"
	bustokens "github.com/twirapp/twir/libs/bus-core/tokens"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

var (
	kickHTTPClient  *http.Client
	kickAPIBaseURL  string
	kickAuthBaseURL string
	kickAppTokenRequester = func(ctx context.Context, parseCtx *types.VariableParseContext) (string, error) {
		resp, err := parseCtx.Services.Bus.Tokens.RequestAppToken.Request(
			ctx,
			bustokens.GetAppTokenRequest{Platform: platformentity.PlatformKick},
		)
		if err != nil {
			return "", fmt.Errorf("request kick app token: %w", err)
		}

		return resp.Data.AccessToken, nil
	}
)

func SetKickClientOptionsForTests(httpClient *http.Client, apiBaseURL, authBaseURL string) func() {
	oldHTTPClient := kickHTTPClient
	oldAPIBaseURL := kickAPIBaseURL
	oldAuthBaseURL := kickAuthBaseURL

	kickHTTPClient = httpClient
	kickAPIBaseURL = apiBaseURL
	kickAuthBaseURL = authBaseURL

	return func() {
		kickHTTPClient = oldHTTPClient
		kickAPIBaseURL = oldAPIBaseURL
		kickAuthBaseURL = oldAuthBaseURL
	}
}

func SetKickAppTokenRequesterForTests(requester func(ctx context.Context, parseCtx *types.VariableParseContext) (string, error)) func() {
	oldRequester := kickAppTokenRequester
	kickAppTokenRequester = requester

	return func() {
		kickAppTokenRequester = oldRequester
	}
}

func newKickAppClient(ctx context.Context, parseCtx *types.VariableParseContext) (*gokick.Client, error) {
	accessToken, err := kickAppTokenRequester(ctx, parseCtx)
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

func GetKickKicksLeaderboard(ctx context.Context, parseCtx *types.VariableParseContext, top int) (*gokick.KicksLeaderboardResponse, error) {
	client, err := newKickAppClient(ctx, parseCtx)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetKicksLeaderboard(ctx, gokick.NewKicksLeaderboardFilter().SetTop(top))
	if err != nil {
		return nil, fmt.Errorf("get kick kicks leaderboard: %w", err)
	}

	return &resp.Result, nil
}

type KickChannelUser struct {
	CreatedAt      string `json:"created_at"`
	FollowingSince string `json:"following_since"`
}

func GetKickChannelUser(
	ctx context.Context,
	parseCtx *types.VariableParseContext,
	channelSlug string,
	userSlug string,
) (*KickChannelUser, error) {
	accessToken, err := kickAppTokenRequester(ctx, parseCtx)
	if err != nil {
		return nil, err
	}

	httpClient := kickHTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	apiBaseURL := kickAPIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = "https://api.kick.com"
	}

	url := fmt.Sprintf("%s/v2/channels/%s/users/%s", apiBaseURL, channelSlug, userSlug)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create kick channel user request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get kick channel user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kick channel user returned status %d", resp.StatusCode)
	}

	var user KickChannelUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode kick channel user: %w", err)
	}

	return &user, nil
}
