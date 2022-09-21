package manage

import (
	"encoding/json"
	model "tsuwari/parser/internal/models"

	"github.com/samber/lo"
)

func CreateRedisBytes(cmd model.ChannelsCommands, response string, setAliases *bool, ) (*[]byte, error) {
	var commandInterface map[string]interface{}
	inrec, _ := json.Marshal(cmd)
	json.Unmarshal(inrec, &commandInterface)
	commandInterface["responses"] = []string{response}

	if setAliases != nil && *setAliases == true {
		commandInterface["aliases"] = []string{}
	}

	bytes, err := json.Marshal(commandInterface)

	if err != nil {
		return nil, err
	}

	return lo.ToPtr(bytes), nil
}