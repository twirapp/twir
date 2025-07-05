package emotes

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type FfzEmote struct {
	Name string `json:"name"`
}

type Set struct {
	Emoticons []FfzEmote
}

type FfzResponse struct {
	Sets map[string]Set `json:"sets"`
}

func GetChannelFfzEmotes(ctx context.Context, channelID string) ([]string, error) {
	reqData := FfzResponse{}

	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://api.frankerfacez.com/v1/room/id/" + channelID)
	if err != nil {
		return nil, err
	}

	var emotes []string
	for _, set := range reqData.Sets {
		mapped := lo.Map(
			set.Emoticons, func(e FfzEmote, _ int) string {
				return e.Name
			},
		)

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}

func GetGlobalFfzEmotes(ctx context.Context) ([]string, error) {
	reqData := FfzResponse{}

	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://api.frankerfacez.com/v1/set/global")
	if err != nil {
		return nil, err
	}

	var emotes []string
	for _, set := range reqData.Sets {
		mapped := lo.Map(
			set.Emoticons, func(e FfzEmote, _ int) string {
				return e.Name
			},
		)

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}
