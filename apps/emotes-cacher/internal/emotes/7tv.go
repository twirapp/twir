package emotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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

func GetChannelSevenTvEmotes(channelID string) ([]string, error) {
	resp, err := http.Get("https://7tv.io/v3/users/twitch/988337552" + channelID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := SevenUserTvResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch 7tv emotes: %w", err)
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

func GetGlobalSevenTvEmotes() ([]string, error) {
	resp, err := http.Get("https://7tv.io/v3/emote-sets/global")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := SevenTvGlobalResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, errors.New("cannot fetch 7tv emotes")
	}

	mappedEmotes := lo.Map(
		reqData.Emotes, func(item SevenTvEmote, _ int) string {
			return item.Name
		},
	)

	return mappedEmotes, nil
}
