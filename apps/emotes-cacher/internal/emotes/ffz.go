package emotes

import (
	"encoding/json"
	"errors"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
	"io"
	"log"
	"net/http"
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

func GetFfzEmotes(channelID string) ([]string, error) {
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
		return nil, errors.New("cannot fetch ffz emotes")
	}

	emotes := []string{}
	for _, set := range reqData.Sets {
		mapped := helpers.Map(set.Emoticons, func(e FfzEmote) string {
			return e.Name
		})

		emotes = append(emotes, mapped...)
	}

	return emotes, nil
}
