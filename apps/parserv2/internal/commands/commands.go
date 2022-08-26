package defaultcommands

import (
	testcommand "tsuwari/parser/internal/commands/test"
	"tsuwari/parser/internal/types"
)

var (
	Commands = make(map[string]types.DefaultCommand)
)

func NewCommands() {
	Commands[testcommand.Command.Name] = testcommand.Command
}