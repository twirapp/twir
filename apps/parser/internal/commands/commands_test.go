package commands

import (
	"testing"

	commandmodel "github.com/twirapp/twir/libs/repositories/commands/model"
	commandsmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

func TestNew_whenQuoteCommandsAreRegistered(t *testing.T) {
	commands := New(&Opts{})
	for _, name := range []string{"quote", "quote add", "quote remove"} {
		if _, ok := commands.DefaultCommands[name]; !ok {
			t.Fatalf("expected %q to be registered", name)
		}
	}
}

func TestFindChannelCommandInInput_whenQuoteSubcommandExists(t *testing.T) {
	commands := &Commands{}
	channelCommands := []commandsmodel.CommandWithGroupAndResponses{
		{Command: commandmodel.Command{Name: "quote"}},
		{Command: commandmodel.Command{Name: "quote add", Aliases: []string{"quote +"}}},
		{Command: commandmodel.Command{Name: "quote remove", Aliases: []string{"quote -"}}},
	}
	cases := []struct {
		input    string
		expected string
	}{
		{input: "quote add hello", expected: "quote add"},
		{input: "quote + hello", expected: "quote add"},
		{input: "quote remove #42", expected: "quote remove"},
		{input: "quote - #42", expected: "quote remove"},
		{input: "quote #42", expected: "quote"},
	}

	for _, testCase := range cases {
		t.Run(testCase.input, func(t *testing.T) {
			result := commands.FindChannelCommandInInput(testCase.input, channelCommands)
			if result.Cmd == nil || result.Cmd.Name != testCase.expected {
				t.Fatalf("expected %q, got %#v", testCase.expected, result.Cmd)
			}
		})
	}
}
