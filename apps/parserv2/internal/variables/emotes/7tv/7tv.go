package emotes7tv

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"
)

type Emote struct {
	Name string `json:"name"`
}

type SevenTvResponse []Emote

const Name = "emotes.7tv"
const Description = "Emotes of channel from https://7tv.app"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	resp, err := http.Get("https://api.7tv.app/v2/users/" + ctx.Context.ChannelId + "/emotes")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	reqData := SevenTvResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, errors.New("cannot fetch ffz emotes")
	}

	mapped := helpers.Map(reqData, func(e Emote) string {
		return e.Name
	})

	result := types.VariableHandlerResult{
		Result: strings.Join(mapped, " "),
	}

	return &result, nil
}
