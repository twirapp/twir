package natshandler

import (
	"encoding/json"
	"strings"
	"tsuwari/parser/internal/permissions"
	types "tsuwari/parser/internal/types"

	"github.com/nats-io/nats.go"
)

func (c natsService) HandleProcessCommand(m *nats.Msg) *[]string {
	data := types.HandleProcessCommandData{}
	json.Unmarshal(m.Data, &data)

	if !strings.HasPrefix(data.Message.Text, "!") {
		return nil
	}
	data.Message.Text = strings.ToLower(data.Message.Text[1:])

	cmds, err := c.commands.GetChannelCommands(data.Channel.Id)
	if err != nil {
		return nil
	}

	cmd := c.commands.FindByMessage(data.Message.Text, cmds)

	if cmd == nil {
		return nil
	}

	hasPerm := permissions.UserHasPermissionToCommand(data.Sender.Badges, "MODERATOR")

	if !hasPerm {
		return nil
	}

	result := c.commands.ParseCommandResponses(cmd, data)

	return &result
}
