package shared

import (
	"context"
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
