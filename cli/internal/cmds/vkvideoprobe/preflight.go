package vkvideoprobe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type VKClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

type vkRequest struct {
	path        string
	query       url.Values
	accessToken string
}

type activeChannelsResponse struct {
	Data []activeChannel `json:"data"`
}

type activeChannel struct {
	Channel struct {
		URL               string `json:"url"`
		WebSocketChannels struct {
			Chat string `json:"chat"`
		} `json:"web_socket_channels"`
	} `json:"channel"`
	Stream struct {
		ID string `json:"id"`
	} `json:"stream"`
}

type websocketTokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type subscriptionTokenResponse struct {
	Data struct {
		ChannelTokens []channelToken `json:"channel_tokens"`
	} `json:"data"`
}

type channelToken struct {
	Channel string `json:"channel"`
	Token   string `json:"token"`
}

func NewVKClient(httpClient *http.Client, baseURL string) (*VKClient, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil || (parsedBaseURL.Scheme != "http" && parsedBaseURL.Scheme != "https") || parsedBaseURL.Host == "" {
		return nil, fmt.Errorf("invalid VK Video API base URL")
	}
	if httpClient == nil {
		return nil, fmt.Errorf("VK Video HTTP client is required")
	}

	return &VKClient{httpClient: httpClient, baseURL: parsedBaseURL}, nil
}

func (client *VKClient) Preflight(ctx context.Context, channelURL, accessToken string) (PreflightResult, error) {
	channelsBody, err := client.get(ctx, vkRequest{path: "/v1/channels/active", accessToken: accessToken})
	if err != nil {
		return PreflightResult{}, fmt.Errorf("get active VK Video channels: %w", err)
	}

	var activeChannels activeChannelsResponse
	if err := json.Unmarshal(channelsBody, &activeChannels); err != nil {
		return PreflightResult{}, fmt.Errorf("decode active VK Video channels: %w", err)
	}

	selected, found := selectChannel(activeChannels.Data, channelURL)
	if !found {
		return PreflightResult{}, fmt.Errorf("active VK Video channel was not found")
	}

	connectionBody, err := client.get(ctx, vkRequest{path: "/v1/websocket/token", accessToken: accessToken})
	if err != nil {
		return PreflightResult{}, fmt.Errorf("get VK Video websocket token: %w", err)
	}

	var connection websocketTokenResponse
	if err := json.Unmarshal(connectionBody, &connection); err != nil {
		return PreflightResult{}, fmt.Errorf("decode VK Video websocket token: %w", err)
	}
	if connection.Data.Token == "" {
		return PreflightResult{}, fmt.Errorf("VK Video websocket token was empty")
	}

	subscriptionBody, err := client.get(ctx, vkRequest{
		path:        "/v1/websocket/subscription_token",
		query:       url.Values{"channels": []string{selected.Channel.WebSocketChannels.Chat}},
		accessToken: accessToken,
	})
	if err != nil {
		return PreflightResult{}, fmt.Errorf("get VK Video websocket subscription token: %w", err)
	}

	var subscription subscriptionTokenResponse
	if err := json.Unmarshal(subscriptionBody, &subscription); err != nil {
		return PreflightResult{}, fmt.Errorf("decode VK Video websocket subscription token: %w", err)
	}

	subscriptionToken, found := selectSubscriptionToken(subscription.Data.ChannelTokens, selected.Channel.WebSocketChannels.Chat)
	if !found {
		return PreflightResult{}, fmt.Errorf("VK Video chat subscription token was not found")
	}

	return PreflightResult{
		ChannelURL:        selected.Channel.URL,
		StreamID:          selected.Stream.ID,
		ChatChannel:       selected.Channel.WebSocketChannels.Chat,
		ConnectionToken:   connection.Data.Token,
		SubscriptionToken: subscriptionToken,
		accessToken:       accessToken,
	}, nil
}

func (client *VKClient) get(ctx context.Context, request vkRequest) ([]byte, error) {
	endpoint := client.baseURL.ResolveReference(&url.URL{Path: request.path})
	endpoint.RawQuery = request.query.Encode()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create VK Video request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+request.accessToken)
	httpRequest.Header.Set("Accept", "application/json")

	response, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("send VK Video request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, maxHTTPResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read VK Video response: %w", err)
	}
	if len(body) > maxHTTPResponseBytes {
		return nil, fmt.Errorf("VK Video response exceeded %d bytes", maxHTTPResponseBytes)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("VK Video returned HTTP %s", response.Status)
	}

	return body, nil
}

func selectChannel(channels []activeChannel, channelURL string) (activeChannel, bool) {
	for _, channel := range channels {
		if channel.Channel.URL == channelURL && channel.Channel.WebSocketChannels.Chat != "" && channel.Stream.ID != "" {
			return channel, true
		}
	}

	return activeChannel{}, false
}

func selectSubscriptionToken(tokens []channelToken, chatChannel string) (string, bool) {
	for _, token := range tokens {
		if token.Channel == chatChannel && strings.TrimSpace(token.Token) != "" {
			return token.Token, true
		}
	}

	return "", false
}
