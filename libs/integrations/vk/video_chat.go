package vk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var ErrActiveStreamNotFound = errors.New("VK Video active stream not found")

type VideoChatClientOpts struct {
	APIBaseURL string
	HTTPClient *http.Client
}

type VideoChatClient struct {
	apiBaseURL string
	httpClient *http.Client
}

type ActiveStream struct {
	ChannelURL string
	ID         string
}

type SendTextMessageInput struct {
	OwnerAccessToken string
	BotAccessToken   string
	Content          string
}

func NewVideoChatClient(opts VideoChatClientOpts) (*VideoChatClient, error) {
	if strings.TrimSpace(opts.APIBaseURL) == "" {
		return nil, errors.New("VK Video chat API base URL is required")
	}
	if err := validateBaseURL(opts.APIBaseURL); err != nil {
		return nil, fmt.Errorf("invalid VK Video chat API base URL: %w", err)
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultOAuthHTTPTimeout}
	}

	return &VideoChatClient{
		apiBaseURL: strings.TrimRight(opts.APIBaseURL, "/"),
		httpClient: httpClient,
	}, nil
}

func (c *VideoChatClient) ResolveActiveStream(
	ctx context.Context,
	ownerAccessToken string,
	botAccessToken string,
) (*ActiveStream, error) {
	currentUserEndpoint, err := url.JoinPath(c.apiBaseURL, "v1", "current_user")
	if err != nil {
		return nil, fmt.Errorf("build VK Video current user endpoint: %w", err)
	}

	currentUserRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, currentUserEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create VK Video current user request: %w", err)
	}
	currentUserRequest.Header.Set("Authorization", "Bearer "+ownerAccessToken)

	var currentUser currentUserResponse
	if err := c.do(currentUserRequest, &currentUser); err != nil {
		return nil, err
	}
	channelURL := strings.TrimSpace(currentUser.Data.Channel.URL)
	if channelURL == "" {
		return nil, errors.New("VK Video current user response is missing channel URL")
	}

	channelEndpoint, err := url.JoinPath(c.apiBaseURL, "v1", "channel")
	if err != nil {
		return nil, fmt.Errorf("build VK Video channel endpoint: %w", err)
	}
	channelRequestURL, err := url.Parse(channelEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parse VK Video channel URL: %w", err)
	}
	query := channelRequestURL.Query()
	query.Set("channel_url", channelURL)
	channelRequestURL.RawQuery = query.Encode()

	channelRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, channelRequestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create VK Video channel request: %w", err)
	}
	channelRequest.Header.Set("Authorization", "Bearer "+botAccessToken)

	var channel channelResponse
	if err := c.do(channelRequest, &channel); err != nil {
		return nil, err
	}
	if channel.Data.Stream == nil || strings.TrimSpace(channel.Data.Stream.ID) == "" {
		return nil, ErrActiveStreamNotFound
	}

	return &ActiveStream{ChannelURL: channelURL, ID: channel.Data.Stream.ID}, nil
}

func (c *VideoChatClient) SendTextMessage(ctx context.Context, input SendTextMessageInput) error {
	stream, err := c.ResolveActiveStream(ctx, input.OwnerAccessToken, input.BotAccessToken)
	if errors.Is(err, ErrActiveStreamNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("resolve VK Video active stream: %w", err)
	}

	endpoint, err := url.JoinPath(c.apiBaseURL, "v1", "chat", "message", "send")
	if err != nil {
		return fmt.Errorf("build VK Video chat send endpoint: %w", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse VK Video chat send URL: %w", err)
	}
	query := u.Query()
	query.Set("channel_url", stream.ChannelURL)
	query.Set("stream_id", stream.ID)
	u.RawQuery = query.Encode()

	body, err := json.Marshal(chatMessageRequest{
		Parts: []chatMessagePart{{Text: chatMessageText{Content: input.Content}}},
	})
	if err != nil {
		return fmt.Errorf("encode VK Video chat message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create VK Video chat send request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+input.BotAccessToken)
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, nil)
}

func (c *VideoChatClient) do(req *http.Request, target any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform VK Video chat request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseSize+1))
	if err != nil {
		return fmt.Errorf("read VK Video chat response: %w", err)
	}
	if int64(len(body)) > maxOAuthResponseSize {
		return fmt.Errorf("VK Video chat response exceeds %d bytes", maxOAuthResponseSize)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &ProviderError{StatusCode: resp.StatusCode}
	}
	if target == nil {
		return nil
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode VK Video chat response: %w", err)
	}

	return nil
}

type currentUserResponse struct {
	Data struct {
		Channel struct {
			URL string `json:"url"`
		} `json:"channel"`
	} `json:"data"`
}

type channelResponse struct {
	Data struct {
		Stream *channelStream `json:"stream"`
	} `json:"data"`
}

type channelStream struct {
	ID string `json:"id"`
}

type chatMessageRequest struct {
	Parts []chatMessagePart `json:"parts"`
}

type chatMessagePart struct {
	Text chatMessageText `json:"text"`
}

type chatMessageText struct {
	Content string `json:"content"`
}
