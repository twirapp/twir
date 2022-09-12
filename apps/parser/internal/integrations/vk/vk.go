package vk

import (
	"encoding/json"
	"fmt"
	model "tsuwari/parser/internal/models"

	req "github.com/imroc/req/v3"
)

type DbData struct {
	UserId *string `json:"userId"`
}

type Vk struct {
	DbData
	integration *model.ChannelInegrationWithRelation
}

func New(integration *model.ChannelInegrationWithRelation) *Vk {
	if integration == nil || !integration.Data.Valid || !integration.Integration.APIKey.Valid {
		return nil
	}

	dbData := DbData{}
	err := json.Unmarshal([]byte(integration.Data.String), &dbData)
	if err != nil || dbData.UserId == nil {
		return nil
	}

	service := Vk{
		DbData:      dbData,
		integration: integration,
	}

	return &service
}

type VkError struct {
	Code int    `json:"error_code,omitempty"`
	Msg  string `error_msg:"msg,omitempty"`
}

type VkAudio struct {
	Artist string `json:"artist,omitempty"`
	Title  string `json:"title,omitempty"`
}

type VkStatus struct {
	Text  *string  `json:"text,omitempty"`
	Audio *VkAudio `json:"audio,omitempty"`
}

type VkResponse struct {
	Error  *VkError  `json:"error,omitempty"`
	Status *VkStatus `json:"response"`
}

func (c *Vk) GetTrack() *string {
	data := VkResponse{}
	var response string

	resp, err := req.R().
		SetQueryParam("access_token", c.integration.Integration.APIKey.String).
		SetQueryParam("uid", *c.UserId).
		SetQueryParam("v", "5.131").
		SetResult(&data).
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
