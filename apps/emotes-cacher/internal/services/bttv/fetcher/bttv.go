package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.betterttv.net/3/cached/users/twitch/"+channelID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reqData BttvResponse
	if err := json.Unmarshal(body, &reqData); err != nil {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.betterttv.net/3/cached/emotes/global", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emotes []BttvEmote
	if err := json.Unmarshal(body, &emotes); err != nil {
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
