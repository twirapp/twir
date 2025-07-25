package fetcher

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
)

type SevenTvEmote struct {
	ID   string `json:"id"`
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

func GetChannelSevenTvEmotes(ctx context.Context, channelID string) ([]emote.Emote, error) {
	data := SevenUserTvResponse{}

	resp, err := req.
		SetContext(ctx).
		SetSuccessResult(&data).
		Get("https://7tv.io/v3/users/twitch/" + channelID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	if data.EmoteSet == nil {
		return nil, nil
	}

	result := make([]emote.Emote, 0, len(data.EmoteSet.Emotes))
	for _, e := range data.EmoteSet.Emotes {
		result = append(
			result,
			emote.Emote{
				ID:   emote.ID(e.ID),
				Name: e.Name,
			},
		)
	}

	return result, nil
}

func GetGlobalSevenTvEmotes(ctx context.Context) ([]emote.Emote, error) {
	data := SevenTvGlobalResponse{}
	resp, err := req.
		SetContext(ctx).
		SetSuccessResult(&data).
		Get("https://7tv.io/v3/emote-sets/global")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	result := make([]emote.Emote, 0, len(data.Emotes))
	for _, e := range data.Emotes {
		result = append(
			result,
			emote.Emote{
				ID:   emote.ID(e.ID),
				Name: e.Name,
			},
		)
	}

	return result, nil
}
