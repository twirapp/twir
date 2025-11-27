package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://7tv.io/v3/users/twitch/"+channelID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data SevenUserTvResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://7tv.io/v3/emote-sets/global", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data SevenTvGlobalResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
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
