package fetcher

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/emotes-cacher/internal/emote"
)

type FfzEmote struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Set struct {
	Emoticons []FfzEmote
}

type FfzResponse struct {
	Sets map[string]Set `json:"sets"`
}

func GetChannelFfzEmotes(ctx context.Context, channelID string) ([]emote.Emote, error) {
	reqData := FfzResponse{}

	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://api.frankerfacez.com/v1/room/id/" + channelID)
	if err != nil {
		return nil, err
	}

	var emotes []emote.Emote
	for _, set := range reqData.Sets {
		mapped := lo.Map(
			set.Emoticons, func(e FfzEmote, _ int) emote.Emote {
				return emote.Emote{
					ID:   emote.ID(e.ID),
					Name: e.Name,
				}
			},
		)

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}

func GetGlobalFfzEmotes(ctx context.Context) ([]emote.Emote, error) {
	reqData := FfzResponse{}

	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://api.frankerfacez.com/v1/set/global")
	if err != nil {
		return nil, err
	}

	var emotes []emote.Emote
	for _, set := range reqData.Sets {
		mapped := lo.Map(
			set.Emoticons, func(e FfzEmote, _ int) emote.Emote {
				return emote.Emote{
					ID:   emote.ID(e.ID),
					Name: e.Name,
				}
			},
		)

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}
