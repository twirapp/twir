package emotes

import (
	"encoding/json"
	"errors"
	"github.com/samber/lo"
	"io"
	"net/http"
)

type SevenTvEmote struct {
	Name string `json:"name"`
}

type SevenTvResponse []SevenTvEmote

func GetSevenTvEmotes(channelID string) ([]string, error) {
	resp, err := http.Get("https://api.7tv.app/v2/users/" + channelID + "/emotes")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := SevenTvResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, errors.New("cannot fetch 7tv emotes")
	}

	mappedEmotes := lo.Map(reqData, func(item SevenTvEmote, _ int) string {
		return item.Name
	})

	return mappedEmotes, nil
}
