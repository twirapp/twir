package natshandler

import (
	"encoding/json"
	"strings"
	"tsuwari/parser/internal/permissions"

	"github.com/nats-io/nats.go"
)

func (c natsService) HandleProcessCommand(m *nats.Msg) *[]string {
	data := HandleProcessCommandData{}
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

	result := c.commands.ParseCommandResponses(cmd)

	return &result
}

type Sender struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Badges      []string `json:"badges"`
}

type Channel struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
}

type Message struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type HandleProcessCommandData struct {
	Channel Channel `json:"channel"`
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}
