package natshandler

import (
	"strings"
	"tsuwari/parser/internal/permissions"

	parserproto "github.com/satont/tsuwari/nats/parser"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func (c natsService) HandleProcessCommand(m *nats.Msg) *[]string {
	data := parserproto.Request{}
	err := proto.Unmarshal(m.Data, &data)
	if err != nil {
		panic(err)
	}

	if !strings.HasPrefix(data.Message.Text, "!") {
		return nil
	}
	data.Message.Text = strings.ToLower(data.Message.Text[1:])

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
