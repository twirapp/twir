package fetcher

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
)

type BttvEmote struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type BttvResponse struct {
	ChannelEmotes []BttvEmote `json:"channelEmotes"`
	SharedEmotes  []BttvEmote `json:"sharedEmotes"`
}

func GetChannelBttvEmotes(ctx context.Context, channelID string) ([]emote.Emote, error) {
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

	result := make([]emote.Emote, 0, len(emotes))
	for _, e := range emotes {
		result = append(
			result,
			emote.Emote{
				ID:   emote.ID(e),
				Name: e,
			},
		)
	}

	return result, nil
}

func GetGlobalBttvEmotes(ctx context.Context) ([]emote.Emote, error) {
	var emotes []BttvEmote

	_, err := req.
		SetContext(ctx).
		SetSuccessResult(&emotes).
		Get("https://api.betterttv.net/3/cached/emotes/global")
	if err != nil {
		return nil, err
	}

	result := make([]emote.Emote, 0, len(emotes))
	for _, e := range emotes {
		result = append(
			result,
			emote.Emote{
				ID:   emote.ID(e.ID),
				Name: e.Code,
			},
		)
	}

	return result, nil
}
