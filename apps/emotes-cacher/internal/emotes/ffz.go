package emotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/samber/lo"
)

type FfzEmote struct {
	Name string `json:"name"`
}

type Set struct {
	Emoticons []FfzEmote
}

type FfzResponse struct {
	Sets map[string]Set `json:"sets"`
}

func GetChannelFfzEmotes(channelID string) ([]string, error) {
	resp, err := http.Get("https://api.frankerfacez.com/v1/room/id/" + channelID)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	reqData := FfzResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch ffz emotes: %w", err)
	}

	emotes := []string{}
	for _, set := range reqData.Sets {
		mapped := lo.Map(
			set.Emoticons, func(e FfzEmote, _ int) string {
				return e.Name
			},
		)

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}

func GetGlobalFfzEmotes() ([]string, error) {
	resp, err := http.Get("https://api.frankerfacez.com/v1/set/global")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	reqData := FfzResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, errors.New("cannot fetch ffz emotes")
	}

	emotes := []string{}
	for _, set := range reqData.Sets {
		mapped := lo.Map(
			set.Emoticons, func(e FfzEmote, _ int) string {
				return e.Name
			},
		)

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}
