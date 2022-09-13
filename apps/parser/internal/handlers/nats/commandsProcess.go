package natshandler

import (
	"strings"
	"tsuwari/parser/internal/permissions"

	parserproto "github.com/satont/tsuwari/nats/parser"
)

func (c natsService) HandleProcessCommand(data parserproto.Request) *[]string {
	if !strings.HasPrefix(data.Message.Text, "!") {
		return nil
	}
	data.Message.Text = data.Message.Text[1:]

	cmds, err := c.commands.GetChannelCommands(data.Channel.Id)
	if err != nil {
		return nil
	}

	cmd := c.commands.FindByMessage(data.Message.Text, cmds)

	if cmd.Cmd == nil {
		return nil
	}

	hasPerm := permissions.UserHasPermissionToCommand(data.Sender.Badges, cmd.Cmd.Permission)

	if !hasPerm {
		return nil
	}

	result := c.commands.ParseCommandResponses(cmd, data)

	return &result
}
