package emotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/samber/lo"
)

type BttvEmote struct {
	Code string `json:"code"`
}

type BttvResponse struct {
	ChannelEmotes []BttvEmote `json:"channelEmotes"`
	SharedEmotes  []BttvEmote `json:"sharedEmotes"`
}

func GetChannelBttvEmotes(channelID string) ([]string, error) {
	resp, err := http.Get("https://api.betterttv.net/3/cached/users/twitch/" + channelID)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := BttvResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch bttv emotes: %w", err)
	}

	emotes := []string{}

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

func GetGlobalBttvEmotes() ([]string, error) {
	resp, err := http.Get("https://api.betterttv.net/3/cached/emotes/global")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	emotes := []BttvEmote{}
	err = json.Unmarshal(body, &emotes)
	if err != nil {
		return nil, errors.New("cannot fetch bttv emotes")
	}

	return lo.Map(
		emotes, func(item BttvEmote, _ int) string {
			return item.Code
		},
	), nil
}
