package emotes

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type SevenTvEmote struct {
	Name string `json:"name"`
}

type SevenUserTvResponse struct {
	EmoteSet *struct {
		Emotes []SevenTvEmote `json:"emotes"`
	} `json:"emote_set"`
}

type SevenTvGlobalResponse struct {
	Emotes []SevenTvEmote `json:"emotes"`
}

func GetChannelSevenTvEmotes(ctx context.Context, channelID string) ([]string, error) {
	reqData := SevenUserTvResponse{}

	resp, err := req.
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://7tv.io/v3/users/twitch/" + channelID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	if reqData.EmoteSet == nil {
		return []string{}, nil
	}

	mappedEmotes := lo.Map(
		reqData.EmoteSet.Emotes, func(item SevenTvEmote, _ int) string {
			return item.Name
		},
	)

	return mappedEmotes, nil
}

func GetGlobalSevenTvEmotes(ctx context.Context) ([]string, error) {
	reqData := SevenTvGlobalResponse{}
	resp, err := req.
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://7tv.io/v3/emote-sets/global")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	mappedEmotes := lo.Map(
		reqData.Emotes, func(item SevenTvEmote, _ int) string {
			return item.Name
		},
	)

	return mappedEmotes, nil
}
