package vk

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
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
	Gorm        *gorm.DB
	Integration *model.ChannelsIntegrations
}

type VK struct {
	gorm        *gorm.DB
	integration *model.ChannelsIntegrations
}

func New(opts Opts) (*VK, error) {
	if !opts.Integration.AccessToken.Valid {
		return nil, fmt.Errorf("integration access token is not valid")
	}

	vk := &VK{
		gorm:        opts.Gorm,
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
	response := vkResponse{}

	resp, err := req.R().
		SetContext(ctx).
		SetQueryParam("access_token", c.integration.AccessToken.String).
		SetQueryParam("v", "5.199").
		SetSuccessResult(&response).
		SetContentType("application/json").
		Get("https://api.vk.com/method/status.get")

	if err != nil || !resp.IsSuccessState() || response.Error != nil {
		return nil, fmt.Errorf("failed to get track from VK: %s", resp.String())
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
