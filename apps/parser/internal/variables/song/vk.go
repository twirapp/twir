package song

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
)

type vkService struct {
	integration *model.ChannelsIntegrations
}

func newVk(integration *model.ChannelsIntegrations) *vkService {
	if integration == nil || !integration.AccessToken.Valid {
		return nil
	}

	service := vkService{
		integration: integration,
	}

	return &service
}

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

func (c *vkService) GetTrack(ctx context.Context) *string {
	data := vkResponse{}
	var response string

	resp, err := req.R().
		SetContext(ctx).
		SetQueryParam("access_token", c.integration.AccessToken.String).
		SetQueryParam("v", "5.131").
		SetSuccessResult(&data).
		SetContentType("application/json").
		Get("https://api.vk.com/method/status.get")

	if err != nil || !resp.IsSuccess() {
		return nil
	}

	if data.Error != nil || data.Status == nil || data.Status.Audio == nil {
		return nil
	}

	status := *data.Status.Audio
	response = fmt.Sprintf(
		"%s — %s",
		status.Artist,
		status.Title,
	)

	return &response
}
