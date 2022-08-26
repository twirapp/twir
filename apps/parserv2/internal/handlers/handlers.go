package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"tsuwari/parser/internal/config/redis"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/pkg/helpers"
)


func GetChannelCommands(channelId string) (*[]types.Command, error) {
	keys, err := redis.Rdb.Keys(redis.RdbCtx, fmt.Sprintf("commands:%s:*", channelId)).Result()

	if err != nil {
		return nil, err
	}

	var cmds = make([]types.Command, len(keys))
	rCmds, err := redis.Rdb.MGet(redis.RdbCtx, keys...).Result()

	if err != nil {
		return nil, err
	}

	for i, cmd := range rCmds {
		parsedCmd := types.Command{}

		err := json.Unmarshal([]byte(cmd.(string)), &parsedCmd)

		if err == nil {
			cmds[i] = parsedCmd
		}
	}

	return &cmds, nil
}

func FindCommandByMessage(input string, cmds *[]types.Command) *types.Command {
	if !strings.HasPrefix(input, "!") {
		return nil
	}

	msg := strings.ToLower(input[1:])
	splittedName := regexp.MustCompile(`[^\s]+`).FindAllString(msg, -1)

	var cmd *types.Command

	length := len(splittedName)

	for i := 0; i < length; i++ {
		query := strings.Join(splittedName, " ")
		for _, c := range *cmds {
			if c.Name == query {
				cmd = &c
				break
			}

			if helpers.Contains(c.Aliases, query) {
				cmd = &c
				break
			}
		}

		if cmd != nil {
			break
		} else {

			splittedName = splittedName[:len(splittedName) - 1]
			continue
		}
	}

	return cmd
}

func UserHasPermissionToCommand(badges []string, commandPermission string) bool {
	commandPermIndex := helpers.IndexOf(types.CommandPerms, commandPermission)

	res := false
	for _, b := range badges {
		idx := helpers.IndexOf(types.CommandPerms, b)
		if idx <= commandPermIndex {
			res = true
			break
		}
	}

	return res
}

func ParseCommandResponses(message types.ChatMessage) {

}


/* func ParseChatMessage(userId *string, channelId string, userName *string, userDisplayName *string, text string) (string, error) {
	return "", nil
} */