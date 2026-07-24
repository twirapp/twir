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
	AccessToken string
	OwnerID     string
	Content     string
}

func NewVideoChatClient(opts VideoChatClientOpts) (*VideoChatClient, error) {
	apiBaseURL := opts.APIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = defaultOAuthDevAPIBaseURL
	}
	if err := validateBaseURL(apiBaseURL); err != nil {
		return nil, fmt.Errorf("invalid VK Video chat API base URL: %w", err)
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultOAuthHTTPTimeout}
	}

	return &VideoChatClient{
		apiBaseURL: strings.TrimRight(apiBaseURL, "/"),
		httpClient: httpClient,
	}, nil
}

func (c *VideoChatClient) ResolveActiveStream(
	ctx context.Context,
	accessToken string,
	ownerID string,
) (*ActiveStream, error) {
	endpoint, err := url.JoinPath(c.apiBaseURL, "v1", "channels", "active")
	if err != nil {
		return nil, fmt.Errorf("build VK Video active channels endpoint: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create VK Video active channels request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var response activeChannelsResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	for _, channel := range response.Data.Channels {
		if string(channel.Owner.ID) != ownerID {
			continue
		}
		if channel.Channel.URL == "" || channel.Stream.ID == "" {
			continue
		}

		return &ActiveStream{
			ChannelURL: channel.Channel.URL,
			ID:         channel.Stream.ID,
		}, nil
	}

	return nil, ErrActiveStreamNotFound
}

func (c *VideoChatClient) SendTextMessage(ctx context.Context, input SendTextMessageInput) error {
	stream, err := c.ResolveActiveStream(ctx, input.AccessToken, input.OwnerID)
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
	req.Header.Set("Authorization", "Bearer "+input.AccessToken)
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

type activeChannelsResponse struct {
	Data struct {
		Channels []activeChannel `json:"channels"`
	} `json:"data"`
}

type activeChannel struct {
	Channel struct {
		URL string `json:"url"`
	} `json:"channel"`
	Owner struct {
		ID flexString `json:"id"`
	} `json:"owner"`
	Stream struct {
		ID string `json:"id"`
	} `json:"stream"`
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
