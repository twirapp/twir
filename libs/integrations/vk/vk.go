package vk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/twirapp/twir/libs/entities/vk_integration"
)

type vkError struct {
	Code int    `json:"error_code,omitempty"`
	Msg  string `error_msg:"msg,omitempty"`
}

type vkAudio struct {
	Artist string `json:"artist,omitempty"`
	Title  string `json:"title,omitempty"`
}

type vkStatus struct {
	Text  *string  `json:"text,omitempty"`
	Audio *vkAudio `json:"audio,omitempty"`
}

type vkResponse struct {
	Error  *vkError  `json:"error,omitempty"`
	Status *vkStatus `json:"response"`
}

type Opts struct {
	Integration vk_integration.Entity
}

type VK struct {
	integration vk_integration.Entity
}

func New(opts Opts) (*VK, error) {
	if opts.Integration.AccessToken == "" {
		return nil, fmt.Errorf("integration access token is empty")
	}

	vk := &VK{
		integration: opts.Integration,
	}

	return vk, nil
}

type Track struct {
	Title  string
	Artist string
	Image  string
}

func (c *VK) GetTrack(ctx context.Context) (*Track, error) {
	u, _ := url.Parse("https://api.vk.com/method/status.get")
	q := u.Query()
	q.Set("access_token", c.integration.AccessToken)
	q.Set("v", "5.199")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get track from VK: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get track from VK: %s", string(body))
	}

	var response vkResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("failed to get track from VK: %s", string(body))
	}

	if response.Status == nil || response.Status.Audio == nil {
		return nil, nil
	}

	return &Track{
		Title:  response.Status.Audio.Title,
		Artist: response.Status.Audio.Artist,
		Image:  "",
	}, nil
}
