package emotes7tv

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

	"github.com/samber/lo"
)

type Emote struct {
	Name string `json:"name"`
}

type SevenTvResponse []Emote

var Variable = types.Variable{
	Name:        "emotes.7tv",
	Description: lo.ToPtr("Emotes of channel from https://7tv.app"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		resp, err := http.Get("https://api.7tv.app/v2/users/" + ctx.ChannelId + "/emotes")
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
	},
}
