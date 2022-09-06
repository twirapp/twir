package emotesffz

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/samber/lo"
)

type Emote struct {
	Name string `json:"name"`
}

type Set struct {
	Emoticons []Emote
}

type FfzResponse struct {
	Sets map[string]Set `json:"sets"`
}

var Variable = types.Variable{
	Name:        "emotes.ffz",
	Description: lo.ToPtr("Emotes of channel from https://frankerfacez.com"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		resp, err := http.Get("https://api.frankerfacez.com/v1/room/id/" + ctx.Context.ChannelId)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		reqData := FfzResponse{}
		err = json.Unmarshal(body, &reqData)
		if err != nil {
			return nil, errors.New("cannot fetch ffz emotes")
		}

		emotions := []string{}
		for _, set := range reqData.Sets {
			mapped := helpers.Map(set.Emoticons, func(e Emote) string {
				return e.Name
			})

			emotions = append(emotions, mapped...)
		}

		result := types.VariableHandlerResult{
			Result: strings.Join(emotions, " "),
		}

		return &result, nil
	},
}
