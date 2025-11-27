package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.frankerfacez.com/v1/room/id/"+channelID, nil)
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

	var reqData FfzResponse
	if err := json.Unmarshal(body, &reqData); err != nil {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.frankerfacez.com/v1/set/global", nil)
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

	var reqData FfzResponse
	if err := json.Unmarshal(body, &reqData); err != nil {
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
