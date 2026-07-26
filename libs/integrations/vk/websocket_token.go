package vk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type OAuthAccessToken string

type WebSocketChannel string

type WebSocketConnectionToken string

type WebSocketSubscriptionToken string

type WebSocketTokenClient struct {
	videoChatClient *VideoChatClient
}

func NewWebSocketTokenClient(opts VideoChatClientOpts) (*WebSocketTokenClient, error) {
	videoChatClient, err := NewVideoChatClient(opts)
	if err != nil {
		return nil, fmt.Errorf("create VK Video WebSocket token client: %w", err)
	}

	return &WebSocketTokenClient{videoChatClient: videoChatClient}, nil
}

func (c *WebSocketTokenClient) DiscoverChatChannel(
	ctx context.Context,
	accessToken OAuthAccessToken,
) (WebSocketChannel, error) {
	currentUserEndpoint, err := url.JoinPath(c.videoChatClient.apiBaseURL, "v1", "current_user")
	if err != nil {
		return "", fmt.Errorf("build VK Video current user endpoint: %w", err)
	}

	currentUserRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, currentUserEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("create VK Video current user request: %w", err)
	}
	currentUserRequest.Header.Set("Authorization", "Bearer "+string(accessToken))

	var currentUserResponse discoverChatChannelCurrentUserResponse
	if err := c.videoChatClient.do(currentUserRequest, &currentUserResponse); err != nil {
		return "", err
	}

	currentChannelURL := currentUserResponse.Data.Channel.URL
	if strings.TrimSpace(currentChannelURL) == "" {
		return "", errors.New("VK Video current user response is missing channel URL")
	}

	channelsEndpoint, err := url.JoinPath(c.videoChatClient.apiBaseURL, "v1", "channels")
	if err != nil {
		return "", fmt.Errorf("build VK Video channels endpoint: %w", err)
	}

	body, err := json.Marshal(discoverChatChannelRequest{
		Channels: []discoverChatChannelRequestChannel{{URL: currentChannelURL}},
	})
	if err != nil {
		return "", fmt.Errorf("encode VK Video channels request: %w", err)
	}
	channelsRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, channelsEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create VK Video channels request: %w", err)
	}
	channelsRequest.Header.Set("Authorization", "Bearer "+string(accessToken))
	channelsRequest.Header.Set("Content-Type", "application/json")

	var channelsResponse discoverChatChannelChannelsResponse
	if err := c.videoChatClient.do(channelsRequest, &channelsResponse); err != nil {
		return "", err
	}

	for _, channel := range channelsResponse.Data.Channels {
		if channel.Channel.URL != currentChannelURL {
			continue
		}
		if strings.TrimSpace(string(channel.Channel.WebSocketChannels.Chat)) == "" {
			return "", errors.New("VK Video channels response is missing chat channel")
		}

		return channel.Channel.WebSocketChannels.Chat, nil
	}

	return "", errors.New("VK Video channels response is missing current channel")
}

func (c *WebSocketTokenClient) ConnectionToken(
	ctx context.Context,
	accessToken OAuthAccessToken,
) (WebSocketConnectionToken, error) {
	endpoint, err := url.JoinPath(c.videoChatClient.apiBaseURL, "v1", "websocket", "token")
	if err != nil {
		return "", fmt.Errorf("build VK Video WebSocket connection token endpoint: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("create VK Video WebSocket connection token request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+string(accessToken))

	var response webSocketConnectionTokenResponse
	if err := c.videoChatClient.do(req, &response); err != nil {
		return "", err
	}
	if strings.TrimSpace(string(response.Data.Token)) == "" {
		return "", errors.New("VK Video WebSocket connection token response is missing token")
	}

	return response.Data.Token, nil
}

func (c *WebSocketTokenClient) SubscriptionToken(
	ctx context.Context,
	accessToken OAuthAccessToken,
	channel WebSocketChannel,
) (WebSocketSubscriptionToken, error) {
	endpoint, err := url.JoinPath(c.videoChatClient.apiBaseURL, "v1", "websocket", "subscription_token")
	if err != nil {
		return "", fmt.Errorf("build VK Video WebSocket subscription token endpoint: %w", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse VK Video WebSocket subscription token URL: %w", err)
	}
	query := u.Query()
	query.Set("channels", string(channel))
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create VK Video WebSocket subscription token request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+string(accessToken))

	var response webSocketSubscriptionTokenResponse
	if err := c.videoChatClient.do(req, &response); err != nil {
		return "", err
	}

	for _, channelToken := range response.Data.ChannelTokens {
		if channelToken.Channel != channel {
			continue
		}
		if strings.TrimSpace(string(channelToken.Token)) == "" {
			break
		}

		return channelToken.Token, nil
	}

	return "", errors.New("VK Video WebSocket subscription token response is missing requested channel token")
}

type webSocketConnectionTokenResponse struct {
	Data struct {
		Token WebSocketConnectionToken `json:"token"`
	} `json:"data"`
}

type webSocketSubscriptionTokenResponse struct {
	Data struct {
		ChannelTokens []webSocketChannelToken `json:"channel_tokens"`
	} `json:"data"`
}

type webSocketChannelToken struct {
	Channel WebSocketChannel           `json:"channel"`
	Token   WebSocketSubscriptionToken `json:"token"`
}

type discoverChatChannelCurrentUserResponse struct {
	Data struct {
		Channel struct {
			URL string `json:"url"`
		} `json:"channel"`
	} `json:"data"`
}

type discoverChatChannelRequest struct {
	Channels []discoverChatChannelRequestChannel `json:"channels"`
}

type discoverChatChannelRequestChannel struct {
	URL string `json:"url"`
}

type discoverChatChannelChannelsResponse struct {
	Data struct {
		Channels []discoverChatChannelChannel `json:"channels"`
	} `json:"data"`
}

type discoverChatChannelChannel struct {
	Channel struct {
		URL               string `json:"url"`
		WebSocketChannels struct {
			Chat WebSocketChannel `json:"chat"`
		} `json:"web_socket_channels"`
	} `json:"channel"`
}
