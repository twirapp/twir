package emotes

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type BttvEmote struct {
	Code string `json:"code"`
}

type BttvResponse struct {
	ChannelEmotes []BttvEmote `json:"channelEmotes"`
	SharedEmotes  []BttvEmote `json:"sharedEmotes"`
}

func GetChannelBttvEmotes(ctx context.Context, channelID string) ([]string, error) {
	reqData := BttvResponse{}

	_, err := req.
		SetContext(ctx).
		SetSuccessResult(&reqData).
		Get("https://api.betterttv.net/3/cached/users/twitch/" + channelID)
	if err != nil {
		return nil, err
	}

	var emotes []string

	mappedChannelEmotes := lo.Map(
		reqData.ChannelEmotes, func(e BttvEmote, _ int) string {
			return e.Code
		},
	)
	mappedSharedEmotes := lo.Map(
		reqData.SharedEmotes, func(e BttvEmote, _ int) string {
			return e.Code
		},
	)

	emotes = append(emotes, mappedChannelEmotes...)
	emotes = append(emotes, mappedSharedEmotes...)

	return emotes, nil
}

func GetGlobalBttvEmotes(ctx context.Context) ([]string, error) {
	var emotes []BttvEmote

	_, err := req.
		SetContext(ctx).
		SetSuccessResult(&emotes).
		Get("https://api.betterttv.net/3/cached/emotes/global")
	if err != nil {
		return nil, err
	}

	return lo.Map(
		emotes, func(item BttvEmote, _ int) string {
			return item.Code
		},
	), nil
}
